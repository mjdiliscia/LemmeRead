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

type ConfigData struct {
	Order  PostsOrder  `json:"order"`
	Filter PostsFilter `json:"filter"`
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
	amc.filepath = path.Join(glib.GetUserConfigDir(), configFilename)

	_, err := os.Stat(amc.filepath)
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
