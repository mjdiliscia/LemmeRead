package model

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"

	"github.com/gotk3/gotk3/glib"
)

type AppModelConfiguration struct {
	config   ConfigData
	filepath string
}

const configDirName = "lemmeread"

type ConfigData struct {
	LemmyServer string      `json:"lemmyServer"`
	LemmyToken  string      `json:"lemmyToken"`
	Order       PostsOrder  `json:"order"`
	Filter      PostsFilter `json:"filter"`
}

type PostsOrder int

const (
	PostOrderActive = iota
	PostOrderHot
	PostOrderScaled
	PostOrderControversial
	PostOrderNew
	PostOrderOld
	PostOderMostComments
	PostOrderNewComments
)

type PostsFilter int

const (
	PostFilterSubscribed = iota
	PostFilterLocal
	PostFilterAll
)

func NewAppModelConfiguration(configFilename string) (amc AppModelConfiguration) {
	configDir := path.Join(glib.GetUserConfigDir(), configDirName)
	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	amc.filepath = path.Join(configDir, configFilename)

	_, err = os.Stat(amc.filepath)
	if os.IsNotExist(err) {
		err = amc.saveConfig()
		if err != nil {
			log.Println(err)
		}
		return
	} else if err != nil {
		log.Printf("Couldn't reach configuration file '%s': %s", amc.filepath, err)
		return
	}

	err = amc.loadConfig()
	if err != nil {
		log.Println(err)
	}

	return
}

func (amc *AppModelConfiguration) GetLemmyServer() string {
	return amc.config.LemmyServer
}

func (amc *AppModelConfiguration) SetLemmyServer(server string) {
	amc.config.LemmyServer = server
	amc.saveConfig()
}

func (amc *AppModelConfiguration) HaveLemmyData() bool {
	return amc.config.LemmyToken != "" && amc.config.LemmyServer != ""
}

func (amc *AppModelConfiguration) GetLemmyToken() string {
	return amc.config.LemmyToken
}

func (amc *AppModelConfiguration) SetLemmyToken(token string) {
	amc.config.LemmyToken = token
	amc.saveConfig()
}

func (amc *AppModelConfiguration) GetOrder() PostsOrder {
	return amc.config.Order
}

func (amc *AppModelConfiguration) SetOrder(order PostsOrder) {
	amc.config.Order = order
	amc.saveConfig()
}

func (amc *AppModelConfiguration) GetFilter() PostsFilter {
	return amc.config.Filter
}

func (amc *AppModelConfiguration) SetFilter(filter PostsFilter) {
	amc.config.Filter = filter
	amc.saveConfig()
}

func (amc *AppModelConfiguration) loadConfig() (err error) {
	var file *os.File
	file, err = os.Open(amc.filepath)
	if err != nil {
		return
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &amc.config)
	return
}

func (amc *AppModelConfiguration) saveConfig() (err error) {
	jsonData, err := json.MarshalIndent(&amc.config, "", "  ")
	if err != nil {
		return
	}

	err = os.WriteFile(amc.filepath, jsonData, 0644)
	return
}
