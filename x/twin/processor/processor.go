package processor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	toml "github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type AccessTokenContent struct {
	AccessToken AccessToken
}

type AccessToken struct {
	Value string
}

type RemoteURLContent struct {
	Remote string
	File   string
}

type ValidatorMonikerContent struct {
	Moniker string
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
	Remote      string
	Dataset_csv string
	Model_csv   string
	Push_to     PushTo
	Validator   string
	Lr          float64
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
	Validator    string
	Min_val_loss float64
	Err_perc     float64
	R2           float64
	RAAE         float64
	RMAE         float64
}

type Processor struct {
	NodeHome string
	Logger   log.Logger
}

func CheckPathFormat(path string) string {
	lastok := strings.Compare(path[len(path)-1:], "/")
	if lastok != 0 {
		path = path + "/"
	}
	return path
}

func NewProcessor(nodeHome string, log log.Logger) Processor {

	nodeHome = CheckPathFormat(nodeHome)

	p := Processor{
		NodeHome: nodeHome,
		Logger:   log,
	}
	return p
}

func (p Processor) GetAccessToken() (string, error) {

	accessTokenFile := p.NodeHome + "access_token.toml"

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

func (p Processor) GetRemoteURL(twinName string) (trainConfURL string, twinRemoteURL string, remoteURL string, err error) {

	file := p.NodeHome + "remote.toml"

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

func (p Processor) GetValidatorMoniker() (string, error) {

	file := p.NodeHome + "validator_moniker.toml"

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var c ValidatorMonikerContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", err
	}

	return c.Moniker, nil
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

func (p Processor) PrepareTraining(ctx sdk.Context, twinName string) (ValidatorTrainData, error) {

	var vtd ValidatorTrainData

	acctoken, err := p.GetAccessToken()
	if err != nil {
		return vtd, err
	}

	mon, err := p.GetValidatorMoniker()
	if err != nil {
		p.Logger.Error(err.Error())
		return vtd, err
	}

	trainDataContent, twinRemoteURL, _, err := p.ReadTrainConfiguration(acctoken, twinName)
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

		if t.Moniker == mon {
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
		Remote:      twinRemoteURL,
		Dataset_csv: trainDataContent.Dataset_csv,
		Model_csv:   trainDataContent.Model_csv,
		Push_to:     trainDataContent.Push_to,
		Validator:   vtd.Moniker,
		Lr:          lr,
	}

	bz, err := json.MarshalIndent(tcfd, "", " ")
	if err != nil {
		return vtd, err
	}
	trainConfFile := p.NodeHome + "train_conf.json"
	err = ioutil.WriteFile(trainConfFile, bz, 0644)
	if err != nil {
		return vtd, err
	}

	return vtd, nil
}

func (p Processor) StartTraining(ctx sdk.Context, lr float64) {

	fileToRun := p.NodeHome + "train.py"
	cmd := exec.Command(fileToRun)

	err := cmd.Start()
	if err != nil {
		p.Logger.Error(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		p.Logger.Error(err.Error())
	}
}

func (p Processor) CheckValidatorsTrainingState(twinName string) (vts []ValidatorTrainingState, err error) {

	accessToken, err := p.GetAccessToken()
	if err != nil {
		return vts, err
	}

	trainDataContent, _, remoteURL, err := p.ReadTrainConfiguration(accessToken, twinName)
	if err != nil {
		return vts, err
	}

	resultsDirURL := filepath.Join(remoteURL + trainDataContent.Push_to.Path)

	for _, vtd := range trainDataContent.ValidatorsTrainData {
		monikerResultsURL := filepath.Join(resultsDirURL, vtd.Moniker+".json")
		body, err := DoHttpRequestAndReturnBody(monikerResultsURL, accessToken)
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

func (p Processor) ReadValidatorsResults(twinName string) (vtr []ValidatorTrainingResults, err error) {

	accessToken, err := p.GetAccessToken()
	if err != nil {
		return vtr, err
	}

	trainDataContent, _, remoteURL, err := p.ReadTrainConfiguration(accessToken, twinName)
	if err != nil {
		return vtr, err
	}

	resultsDirURL := filepath.Join(remoteURL + trainDataContent.Push_to.Path)

	for _, vtd := range trainDataContent.ValidatorsTrainData {
		monikerResultsURL := filepath.Join(resultsDirURL, vtd.Moniker+".json")
		body, err := DoHttpRequestAndReturnBody(monikerResultsURL, accessToken)
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

func DoHttpRequestAndReturnBody(fileURL string, accessToken string) ([]byte, error) {

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
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
