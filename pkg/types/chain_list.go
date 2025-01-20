package types

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

type ChainList struct {
	Chains    []ChainListItem `json:"chains"`
	ChainsMap map[string]*ChainListItem
}

type ChainListItem struct {
	Name           string         `json:"name"`
	ChainId        int            `json:"chainId"`
	ShortName      string         `json:"shortName"`
	NetworkId      int            `json:"networkId"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	Rpc            []string       `json:"rpc"`
	Faucets        []string       `json:"faucets"`
	InfoURL        string         `json:"infoURL"`
}

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

func UpdateChainList() (*ChainList, error) {
	configPath := utils.ResolvePath("~/.khedra")
	file.EstablishFolder(configPath)
	chainUrl := "https://chainid.network/chains_mini.json"
	chainsFn := filepath.Join(configPath, "chains.json")
	if bytes, err := utils.DownloadAndStore(chainUrl, chainsFn, 24*time.Hour); err != nil {
		return &ChainList{}, err
	} else {
		var chainList ChainList
		err := json.Unmarshal(bytes, &chainList.Chains)
		if err != nil {
			return &ChainList{}, err
		}
		chainList.ChainsMap = make(map[string]*ChainListItem)
		for _, chain := range chainList.Chains {
			chainList.ChainsMap[chain.ShortName] = &chain
		}
		return &chainList, nil
	}
}
