package processor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	Url string
}

type ValidatorMonikerContent struct {
	Moniker string
}

type TrainDataContent struct {
	ValidatorsTrainData []ValidatorTrainData
}

type ValidatorTrainData struct {
	Name string
	Lr   json.Number
}

type Processor struct {
	NodeHome string
	Logger   log.Logger
}

func NewProcessor(nodeHome string, log log.Logger) Processor {

	lastok := strings.Compare(nodeHome[len(nodeHome)-1:], "/")
	if lastok != 0 {
		nodeHome = nodeHome + "/"
	}

	p := Processor{
		NodeHome: nodeHome,
		Logger:   log,
	}
	return p
}

func (p Processor) PrepareTraining(ctx sdk.Context) (ValidatorTrainData, error) {

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

	c, err := p.ReadTrainConfiguration(acctoken)
	if err != nil {
		p.Logger.Error(err.Error())
		return vtd, err
	}

	var found bool

	for _, t := range c.ValidatorsTrainData {
		lr, err := t.Lr.Float64()
		if err != nil {
			return vtd, err
		}
		p.Logger.Error(fmt.Sprintf("val: %s  --> lr: %f", t.Name, lr))

		if t.Name == mon {
			vtd = t
			found = true
		}
	}

	if !found {
		return vtd, fmt.Errorf("Train data not found")
	}

	return vtd, nil
}

func (p Processor) StartTraining(ctx sdk.Context, lr float64) {

	name := p.NodeHome + "train.py"
	cmd := exec.Command(name)

	err := cmd.Start()
	if err != nil {
		p.Logger.Error(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		p.Logger.Error(err.Error())
	}
}

func (p Processor) ReadTrainConfiguration(accessToken string) (TrainDataContent, error) {

	var r TrainDataContent

	fileURL, err := p.GetRemoteURL()
	if err != nil {
		return r, err
	}

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return r, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return r, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	json.Unmarshal(body, &r)

	return r, nil
}

func (p Processor) GetAccessToken() (string, error) {

	accessTokenFile := p.NodeHome + "accessToken.toml"

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

func (p Processor) GetRemoteURL() (string, error) {

	file := p.NodeHome + "remote.toml"

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var c RemoteURLContent
	err = toml.Unmarshal(bz, &c)
	if err != nil {
		return "", err
	}

	return c.Url, nil
}

func (p Processor) GetValidatorMoniker() (string, error) {

	file := p.NodeHome + "validatorMoniker.toml"

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
