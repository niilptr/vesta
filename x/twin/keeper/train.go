package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"

	toml "github.com/BurntSushi/toml"
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

func (k Keeper) StartTraining(ctx sdk.Context, name string, creator string) error {

	acctoken, err := k.GetAccessToken()
	if err != nil {
		return err
	}

	mon, err := k.GetValidatorMoniker()
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}

	c, err := k.ReadTrainConfiguration(acctoken)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}

	for _, t := range c.ValidatorsTrainData {
		lr, err := t.Lr.Float64()
		if err != nil {
			return err
		}
		k.Logger(ctx).Error(fmt.Sprintf("val: %s  --> lr: %f", t.Name, lr))

		if t.Name == mon {
			go runTraining(lr)
		}
	}

	return nil
}

func runTraining(lr float64) {

}

func (k Keeper) ReadTrainConfiguration(accessToken string) (TrainDataContent, error) {

	var r TrainDataContent

	fileURL, err := k.GetRemoteURL()
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

func (k Keeper) GetAccessToken() (string, error) {

	accessTokenFile := k.nodeHome + "accessToken.toml"

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

func (k Keeper) GetRemoteURL() (string, error) {

	file := k.nodeHome + "remote.toml"

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

func (k Keeper) GetValidatorMoniker() (string, error) {

	file := k.nodeHome + "validatorMoniker.toml"

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
