// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	ethcommon "github.com/hacpy/go-ethereum/common"
	"os"
	"path/filepath"

	//ethcommon "github.com/hacpy/go-ethereum/common"

	"github.com/hacpy/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"
const DefaultKeystorePath = "./keystore"
const DefaultBlockTimeout = int64(180) // 3 minutes
const InitialEndPointId = 0

var EndPointParseError = errors.New("json: cannot unmarshal string into Go struct field RawChainConfig.chains.endpoint of type []string")

type Config struct {
	Chains       []RawChainConfig `json:"chains"`
	KeystorePath string           `json:"keystorePath,omitempty"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
	Name     		string              `json:"name"`
	Type     		string              `json:"type"`
	Id       		string              `json:"id"`       // ChainID
	Endpoint 		[]string   	     `json:"endpoint"` // url for rpc endpoint
	From     		string              `json:"from"`     // address of key to use
	Opts     		map[string]string   `json:"opts"`
	OtherRelayer 	[]string 		 `json:"otherRelayer"`
}

func NewConfig() *Config {
	return &Config{
		Chains: []RawChainConfig{},
	}
}

func (c *Config) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(*c); err != nil {
		log.Warn("error marshalling json", "err", err)
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		log.Warn("error creating config file", "err", err)
	}
	_, err = newFile.Write(raw)
	if err != nil {
		log.Warn("error writing to config file", "err", err)
	}

	if err := newFile.Close(); err != nil {
		log.Warn("error closing file", "err", err)
	}
	return newFile
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		if chain.Type == "" {
			return fmt.Errorf("required field chain.Type empty for chain %s", chain.Id)
		}
		if len(chain.Endpoint) == 0 {
			return fmt.Errorf("required field chain.Endpoint empty for chain %s", chain.Id)
		} else if chain.Endpoint[0] == "" {
			return fmt.Errorf("required field chain.Endpoint empty for chain %s", chain.Id)
		}
		if chain.Name == "" {
			return fmt.Errorf("required field chain.Name empty for chain %s", chain.Id)
		}
		if chain.Id == "" {
			return fmt.Errorf("required field chain.Id empty for chain %s", chain.Id)
		}
		if chain.From == "" {
			return fmt.Errorf("required field chain.From empty for chain %s", chain.Id)
		}
		/// Convert to eth address
		if len(chain.From) > 3 && chain.From[:3] == chain.Opts["prefix"] {
			addr, _ := ethcommon.PlatonToEth(chain.From)
			chain.From = string(addr)
		}
		if len(chain.Opts["bridge"]) > 3 && chain.Opts["bridge"][:3] == chain.Opts["prefix"] {
			addr, _ := ethcommon.PlatonToEth(chain.Opts["bridge"])
			address := ethcommon.BytesToAddress(addr)
			chain.Opts["bridge"] = address.String()
		}
		if len(chain.Opts["erc20Handler"]) > 3 && chain.Opts["erc20Handler"][:3] == chain.Opts["prefix"] {
			addr, _ := ethcommon.PlatonToEth(chain.Opts["erc20Handler"])
			address := ethcommon.BytesToAddress(addr)
			chain.Opts["erc20Handler"] = address.String()
		}
	}

	return nil
}

func GetConfig(ctx *cli.Context) (*Config, error) {
	var fig Config
	path := DefaultConfigPath
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &fig)
	if err != nil {
		log.Warn("err loading json file", "err", err.Error())
		return &fig, err
	}
	if ksPath := ctx.String(KeystorePathFlag.Name); ksPath != "" {
		fig.KeystorePath = ksPath
	}
	log.Debug("Loaded config", "path", path)
	err = fig.validate()
	if err != nil {
		return nil, err
	}
	return &fig, nil
}

func loadConfig(file string, config *Config) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
