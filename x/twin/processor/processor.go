package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	toml "github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

const defaultTimeout time.Duration = 1 * time.Second

// TODO: Put these files together
const access_token_file = "access_token.toml"
const moniker_file = "moniker.toml"
const address_file = "address.toml"
const remote_file = "remote.toml"

const training_script = "train.py"
const validation_script = "validate.py"
const confirm_train_phase_ended_script = "confirm_train_phase_ended.sh"
const confirm_best_train_result_script = "confirm_best_train_result_is.sh"

type AccessTokenContent struct {
	AccessToken AccessToken
}

type AccessToken struct {
	Value string
}

type MonikerContent struct {
	Moniker string
}

type AddressContent struct {
	Address string
}

type RemoteURLContent struct {
	Remote string
	File   string
}

type TrainDataContent struct {
	Dataset_csv         string
	Model_csv           string
	Push_to             PushTo
	ValidatorsTrainData []ValidatorTrainData
}

type PushTo struct {
	Owner     string
	Repo      string
	Branch    string
	Path      string
	CommitMsg string
}

type TrainConfFileData struct {
	TwinRemoteDirectory string
	Dataset_csv         string
	Model_csv           string
	Push_to             PushTo
	Moniker             string
	Lr                  float64
}
type ValidatorTrainData struct {
	Moniker string
	Lr      json.Number
}

type ValidatorTrainingState struct {
	Moniker  string
	Complete bool
}

type ValidatorTrainingResults struct {
	Validator      string
	Min_val_loss   float32
	Err_perc       float32
	R2             float32
	RAAE           float32
	RMAE           float32
	SHA256         string
	SHAComputation []string
	NNParams       NNParams
}

type NNParams struct {
	NNWeights []NNWeight
	NNBiases  []NNBias
}

type NNWeight struct {
	Value [][]float64
}

type NNBias struct {
	Value []float64
}

type Processor struct {
	nodeHome    string
	Logger      log.Logger
	address     string
	moniker     string
	accessToken string
}

func CheckPathFormat(path string) string {
	lastok := strings.Compare(path[len(path)-1:], "/")
	if lastok != 0 {
		path = path + "/"
	}
	return path
}

func NewProcessor(nodeHome string, log log.Logger) (Processor, error) {

	nodeHome = CheckPathFormat(nodeHome)

	p := Processor{
		nodeHome: nodeHome,
		Logger:   log,
	}

	access_token, err := p.getAccessToken()
	if err != nil {
		p.Logger.Error(err.Error())
		return p, err
	}

	moniker, err := p.getMoniker()
	if err != nil {
		p.Logger.Error(err.Error())
		return p, err
	}

	address, err := p.getAddress()
	if err != nil {
		p.Logger.Error(err.Error())
		return p, err
	}

	p.accessToken = access_token
	p.address = address
	p.moniker = moniker

	return p, nil
}

func (p Processor) GetNodeHome() string {
	return p.nodeHome
}

func (p Processor) GetAccessToken() string {
	return p.accessToken
}

func (p Processor) GetAddress() string {
	return p.address
}

func (p Processor) GetMoniker() string {
	return p.moniker
}

func (p Processor) getAccessToken() (string, error) {

	accessTokenFile := p.GetNodeHome() + access_token_file

	bz, err := ioutil.ReadFile(accessTokenFile)
	if err != nil {
		return "", err
	}

	var c AccessTokenContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", err
	}

	return c.AccessToken.Value, nil
}

func (p Processor) getMoniker() (string, error) {

	file := p.GetNodeHome() + moniker_file

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var c MonikerContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", err
	}

	return c.Moniker, nil
}

func (p Processor) getAddress() (string, error) {

	file := p.GetNodeHome() + address_file

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var c AddressContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", err
	}

	return c.Address, nil
}

func (p Processor) GetRemoteURL(twinName string) (trainConfURL string, twinRemoteURL string, remoteURL string, err error) {

	file := p.GetNodeHome() + remote_file

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return "", "", "", err
	}

	var content RemoteURLContent
	err = toml.Unmarshal(bz, &content)
	if err != nil {
		return "", "", "", err
	}

	if content.Remote == "" || content.File == "" {
		return "", "", "", fmt.Errorf("remote configuration file misconfigured")
	}

	content.Remote = CheckPathFormat(content.Remote)
	trainConfURL = content.Remote + twinName + "/" + content.File
	remoteURL = content.Remote
	twinRemoteURL = content.Remote + twinName

	return trainConfURL, twinRemoteURL, remoteURL, nil
}

func (p Processor) ReadTrainConfiguration(accessToken string, twinName string) (tdc TrainDataContent, twinRemoteURL string, remoteURL string, err error) {

	fileURL, twinRemoteURL, remoteURL, err := p.GetRemoteURL(twinName)
	if err != nil {
		return tdc, "", "", err
	}

	body, err := DoHttpRequestAndReturnBody(fileURL, accessToken)
	if err != nil {
		return tdc, "", "", err
	}

	json.Unmarshal(body, &tdc)

	return tdc, twinRemoteURL, remoteURL, nil
}

// Read the train configuration settings from the remote repository and write the specific trainer
// node settings to a local file. This file will be later used by the training program to get
// the information it needs.
func (p Processor) PrepareTraining(ctx sdk.Context, twinName string) (ValidatorTrainData, error) {

	var vtd ValidatorTrainData

	trainDataContent, twinRemoteURL, _, err := p.ReadTrainConfiguration(p.GetAccessToken(), twinName)
	if err != nil {
		p.Logger.Error(err.Error())
		return vtd, err
	}

	var found bool
	for _, t := range trainDataContent.ValidatorsTrainData {
		lr, err := t.Lr.Float64()
		if err != nil {
			return vtd, err
		}
		p.Logger.Error(fmt.Sprintf("val: %s  --> lr: %f", t.Moniker, lr))

		if t.Moniker == p.GetMoniker() {
			vtd = t
			found = true
		}
	}

	if !found {
		return vtd, fmt.Errorf("Train data not found")
	}

	lr, err := vtd.Lr.Float64()
	if err != nil {
		return vtd, err
	}

	tcfd := TrainConfFileData{
		TwinRemoteDirectory: twinRemoteURL,
		Dataset_csv:         trainDataContent.Dataset_csv,
		Model_csv:           trainDataContent.Model_csv,
		Push_to:             trainDataContent.Push_to,
		Moniker:             vtd.Moniker,
		Lr:                  lr,
	}

	bz, err := json.MarshalIndent(tcfd, "", " ")
	if err != nil {
		return vtd, err
	}
	trainConfFile := p.GetNodeHome() + "train_conf.json"
	err = ioutil.WriteFile(trainConfFile, bz, 0644)
	if err != nil {
		return vtd, err
	}

	return vtd, nil
}

func (p Processor) Train() error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + training_script

	cmd := exec.Command("python", scriptPath, "--input-dir", p.GetNodeHome())
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		p.Logger.Error(err.Error())
	}

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	p.Logger.Error(fmt.Sprintf("Captured training stdout\n%s\n", stdout))
	p.Logger.Error(fmt.Sprintf("Captured training stderr\n%s\n", stderr))

	return err
}

func (p Processor) CheckValidatorsTrainingState(twinName string) (vts []ValidatorTrainingState, err error) {

	trainDataContent, _, remoteURL, err := p.ReadTrainConfiguration(p.GetAccessToken(), twinName)
	if err != nil {
		return vts, err
	}

	resultsDirURL := filepath.Join(remoteURL + trainDataContent.Push_to.Path)

	for _, vtd := range trainDataContent.ValidatorsTrainData {

		monikerResultsURL := filepath.Join(resultsDirURL, vtd.Moniker+".json")

		body, err := DoHttpRequestAndReturnBody(monikerResultsURL, p.GetAccessToken())

		if err != nil || body == nil {
			p.Logger.Error(err.Error())
			vts = append(vts, ValidatorTrainingState{Moniker: vtd.Moniker, Complete: false})
		}

		if body != nil {
			vts = append(vts, ValidatorTrainingState{Moniker: vtd.Moniker, Complete: true})
		}

	}

	return vts, nil
}

func (p Processor) ReadValidatorsTrainingResults(twinName string) (vtr []ValidatorTrainingResults, err error) {

	trainDataContent, _, remoteURL, err := p.ReadTrainConfiguration(p.GetAccessToken(), twinName)
	if err != nil {
		return vtr, err
	}

	resultsURL := CheckPathFormat(remoteURL) + trainDataContent.Push_to.Path

	for _, vtd := range trainDataContent.ValidatorsTrainData {

		// resultsURL will be a string like "https://.../twin01/results/{moniker}.json" .
		// To get the actual results "{moniker}" must be replaced with the actual
		// trainer moniker.
		monikerResultsURL := strings.ReplaceAll(resultsURL, "{moniker}", vtd.Moniker)

		body, err := DoHttpRequestAndReturnBody(monikerResultsURL, p.GetAccessToken())
		if err != nil {
			p.Logger.Error(err.Error())
			continue
		}

		var trainingResults ValidatorTrainingResults

		json.Unmarshal(body, &trainingResults)

		vtr = append(vtr, trainingResults)
	}

	return vtr, nil
}

func (p Processor) GetBestTrainingResult(vtr []ValidatorTrainingResults) (idx int, trainerMoniker string, newTwinHash string) {

	var best_score float32 = 0

	for i, v := range vtr {

		score := (100-v.Err_perc)/100 + v.R2

		if score > best_score {

			best_score = score
			idx = i
		}
	}

	return idx, vtr[idx].Validator, vtr[idx].SHA256

}

func DoHttpRequestAndReturnBody(fileURL string, accessToken string) ([]byte, error) {

	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	client := &http.Client{}

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return []byte{}, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

// ValidateBestTrainingResult will call validate.py, a python script that will:
// 1. Compute the training results hash;
// 2. Compare computed hash with the one saved in the results data structure;
// 3. Compare computed hash with the one given in input to ValidateBestTrainingResult function;
// 4. Compute accuracy metrics of the twin model;
// 5. Compare computed accuracy metrics with the ones saved in the results data structure.
// If all these checks are positive then validate.py will exit with status code 0, meaning that
// the results are valid.
// If checks fails, an exit code different from zero is returned by validate.py, with reason
// in the stderr.
// If checks cannot be performed due to misconfiguration or problems reaching remote resurces,
// an exit code different from zero is returned by validate.py with an error message in the
// stderr starting with "Fail".
func (p Processor) ValidateBestTrainingResult(twinName string, trainerMoniker string, twinHash string) (isResultValid bool, reasonWhyFalse string, err error) {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + validation_script

	cmd := exec.Command("python", scriptPath, "--input-dir", p.GetNodeHome(), "--twin-hash", twinHash)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	p.Logger.Error(fmt.Sprintf("Captured training stdout\n%s\n", stdout))
	p.Logger.Error(fmt.Sprintf("Captured training stderr\n%s\n", stderr))

	if err != nil {
		p.Logger.Error(err.Error())

		// Stderr is empty if results are valid.
		if len(stderr) > 3 {

			if stderr[:4] == "Fail" {
				// Errors related to process failures returns an error code different from 0 and
				// a message in the stderr that starts with "Fail".
				return false, "", fmt.Errorf(strings.ReplaceAll(stderr, "\n", ""))

			} else {
				// When the results are not valid an error code different from 0 is returned and a
				// message in the stderr that does not start with "Fail", explaining the cause.
				return false, stderr, nil
			}
		}
	}

	return true, "", nil
}

func (p Processor) BroadcastConfirmationTrainingPhaseEnded() error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + confirm_train_phase_ended_script

	cmd := exec.Command("bash", scriptPath, "-f", p.address)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	p.Logger.Error(fmt.Sprintf("Captured training stdout\n%s\n", stdout))
	p.Logger.Error(fmt.Sprintf("Captured training stderr\n%s\n", stderr))

	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}

	return nil

}

func (p Processor) BroadcastConfirmationBestResultIsValid(resultTwinHash string) error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + confirm_best_train_result_script

	cmd := exec.Command("bash", scriptPath, "-f", p.address, "-r", resultTwinHash)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	p.Logger.Error(fmt.Sprintf("Captured training stdout\n%s\n", stdout))
	p.Logger.Error(fmt.Sprintf("Captured training stderr\n%s\n", stderr))

	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}

	return nil

}
