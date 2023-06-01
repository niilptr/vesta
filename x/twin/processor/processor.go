package processor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"vesta/x/twin/types"

	toml "github.com/BurntSushi/toml"
	"github.com/tendermint/tendermint/libs/log"
)

// ====================================================================================
// Settings
// ====================================================================================

// Paths for the local node directory to reach twin module configuration,
// training and confirmation files.
const nodeHomeToModuleHome = "twin-module/"
const moduleConfigurationDir = nodeHomeToModuleHome + "configuration/"
const moduleTrainingDir = nodeHomeToModuleHome + "training/"
const moduleTrainingCoreDir = moduleTrainingDir + "core/"
const moduleConfirmationDir = nodeHomeToModuleHome + "confirmation/"

// Local configuration, training, and confirmation files.
const twin_module_configuration_file = "twin.toml"
const training_script = "train.py"
const validation_script = "validate.py"
const confirm_train_phase_ended_script = "confirm_train_phase_ended.sh"
const confirm_best_train_result_script = "confirm_best_train_result_is.sh"

// Default time-out for http get request.
const defaultTimeout time.Duration = 1 * time.Second

// ====================================================================================
// Data structures
// ====================================================================================

// Local twin module configuration file (twin.toml) data structure.
type TwinModuleConfigurationContent struct {
	TrainConfigurationPath TrainConfigurationPath
	AccessToken            AccessToken
	Trainer                Trainer
}

type TrainConfigurationPath struct {
	RemoteURL string
	File      string
}

type AccessToken struct {
	Token string
}

type Trainer struct {
	Address string
	Moniker string
}

// Remote json training configuration file in which the information about the training are saved.
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

// Remote results json file that will be uploaded by the trainer after local training is complete.
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

// Processor is an abstracted element. It's role is to manage the training process,
// upload the training results, choose the best result and broadcast the information.
type Processor struct {
	nodeHome                     string
	Logger                       log.Logger
	address                      string
	moniker                      string
	accessToken                  string
	remoteURL                    string
	remoteTrainConfigurationFile string
}

// ====================================================================================
// Utilities
// ====================================================================================

// CheckPathFormat ensures a string terminates with a slash.
func CheckPathFormat(path string) string {
	lastok := strings.Compare(path[len(path)-1:], "/")
	if lastok != 0 {
		path = path + "/"
	}
	return path
}

// Read a file stored in the remote repository, returning its body.
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

	if resp.StatusCode == 404 {
		return []byte{}, fmt.Errorf("Error 404 file not found")
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func computeHash(bz []byte) string {

	hash := sha256.Sum256(bz)

	return hex.EncodeToString(hash[:])
}

// ====================================================================================
// Processor
// ====================================================================================
func NewProcessor(nodeHome string, log log.Logger) (Processor, error) {

	nodeHome = CheckPathFormat(nodeHome)

	p := Processor{
		nodeHome: nodeHome,
		Logger:   log,
	}

	accessToken, trainerAddress, trainerMoniker, remoteURL, trainConfigurationFile, err := p.getTwinModuleConfiguration()

	if err != nil {
		log.Error(err.Error())
		return p, err
	}

	p.accessToken = accessToken
	p.address = trainerAddress
	p.moniker = trainerMoniker
	p.remoteURL = remoteURL
	p.remoteTrainConfigurationFile = trainConfigurationFile

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

func (p Processor) GetRemoteURL() string {
	return p.remoteURL
}

func (p Processor) GetRemoteTrainConfigurationFile() string {
	return p.remoteTrainConfigurationFile
}

func (p Processor) getTwinModuleConfiguration() (
	accessToken string,
	trainerAddress string,
	trainerMoniker string,
	remoteURL string,
	remoteTrainConfigurationFile string,
	err error,
) {

	twinModuleConfiguration := p.GetNodeHome() + moduleConfigurationDir + twin_module_configuration_file

	bz, err := ioutil.ReadFile(twinModuleConfiguration)
	if err != nil {
		return "", "", "", "", "", err
	}

	var c TwinModuleConfigurationContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", "", "", "", "", err
	}

	return c.AccessToken.Token, c.Trainer.Address, c.Trainer.Moniker, c.TrainConfigurationPath.RemoteURL, c.TrainConfigurationPath.File, nil
}

// Read the remote json training configuration file.
func (p Processor) ReadTrainConfigurationAndVerifyHash(twinName string, trainHash string) (tdc TrainDataContent, twinRemoteURL string, err error) {

	twinRemoteURL = CheckPathFormat(p.remoteURL) + twinName + "/"
	fileURL := twinRemoteURL + p.remoteTrainConfigurationFile

	body, err := DoHttpRequestAndReturnBody(fileURL, p.GetAccessToken())
	if err != nil {
		return tdc, "", err
	}

	hash := computeHash(body)
	if hash != trainHash {
		return tdc, "", types.ErrTrainConfigurationHashNotMatch
	}

	json.Unmarshal(body, &tdc)

	return tdc, twinRemoteURL, nil
}

// Read the train configuration settings from the remote repository and write the specific trainer
// node settings to a local file. This file will be later used by the training program to get
// the information it needs.
func (p Processor) VerifyHashAndPrepareTraining(twinName string, trainConfHash string) (ValidatorTrainData, error) {

	var vtd ValidatorTrainData

	trainDataContent, twinRemoteURL, err := p.ReadTrainConfigurationAndVerifyHash(twinName, trainConfHash)
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
	trainConfFile := p.GetNodeHome() + moduleTrainingDir + "train_conf.json"
	err = ioutil.WriteFile(trainConfFile, bz, 0644)
	if err != nil {
		return vtd, err
	}

	return vtd, nil
}

// Start the training process. It will call train.py, an external script that will
// perform the training and upload the results to the remote repository.
func (p Processor) Train() error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + moduleTrainingCoreDir + training_script

	cmd := exec.Command("python", scriptPath, "--module-home", p.GetNodeHome()+nodeHomeToModuleHome, "--access-token", p.GetAccessToken())
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

// Check the remote repository for commited training results and returns if trainers
// completed the training or not.
func (p Processor) CheckValidatorsTrainingState(twinName string, trainConfHash string) (vts []ValidatorTrainingState, err error) {

	trainDataContent, _, err := p.ReadTrainConfigurationAndVerifyHash(twinName, trainConfHash)
	if err != nil {
		return vts, err
	}

	resultsDirURL := filepath.Join(p.remoteURL + trainDataContent.Push_to.Path)

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

// Read the training results from the remote repository.
func (p Processor) ReadValidatorsTrainingResults(twinName string, trainConfHash string) (vtr []ValidatorTrainingResults, err error) {

	trainDataContent, _, err := p.ReadTrainConfigurationAndVerifyHash(twinName, trainConfHash)
	if err != nil {
		return vtr, err
	}

	resultsURL := CheckPathFormat(p.remoteURL) + trainDataContent.Push_to.Path

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

// Select best training result from all the results committed.
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

	scriptPath := p.GetNodeHome() + moduleTrainingCoreDir + validation_script

	cmd := exec.Command("python", scriptPath, "--module-home", p.GetNodeHome()+nodeHomeToModuleHome, "--twin-hash", twinHash, "--access-token", p.GetAccessToken())
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

// Broadcast the network that all training have been completed (after completed or
// timeout reached).It will call an external bash script that will perform an
// automatic confirm-train-phase-ended transaction, singing from the trainer key.
func (p Processor) BroadcastConfirmationTrainingPhaseEnded() error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + moduleConfirmationDir + confirm_train_phase_ended_script

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

// Broadcast the network that the best result among all is the one identified with the
// provided hash. It will call an external bash script that will perform an
// automatic confirm-best-train-result-is transaction, signing from the trainer key.
func (p Processor) BroadcastConfirmationBestResultIsValid(resultTwinHash string) error {

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	scriptPath := p.GetNodeHome() + moduleConfirmationDir + confirm_best_train_result_script

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

// ====================================================================================
// Go routines
// ====================================================================================
func StartProcessorForTrainingTwin(nodeHome string, log log.Logger, twinName string, trainConfigurationsHash string) {

	p, err := NewProcessor(nodeHome, log)

	if err != nil {
		p.Logger.Error(
			fmt.Sprintf("Local training not start: %s", err.Error()),
		)
		return
	}

	_, err = p.VerifyHashAndPrepareTraining(twinName, trainConfigurationsHash)
	if err == nil {
		p.Train()

	} else {
		// TODO: If due to hash mismatch a notification sould be broadcasted to
		// stop training upon majority agrees.
		p.Logger.Error(
			fmt.Sprintf("Local training not start: %s", err.Error()),
		)
		return
	}

	return
}
