package types

import (
	"path/filepath"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

type ChainList struct {
	Chains    []ChainListItem `json:"chains"`
	ChainsMap map[int]*ChainListItem
}

type ChainListItem struct {
	Name           string         `json:"name"`
	Chain          string         `json:"chain"`
	Icon           string         `json:"icon"`
	Rpc            []string       `json:"rpc"`
	Faucets        []string       `json:"faucets"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	InfoURL        string         `json:"infoURL"`
	ShortName      string         `json:"shortName"`
	ChainID        int            `json:"chainId"`
	NetworkID      int            `json:"networkId"`
	Explorers      []Explorer     `json:"explorers"`
}

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

type Explorer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Standard string `json:"standard"`
}

func UpdateChainList() (*ChainList, error) {
	configPath := utils.ResolvePath("~/.khedra")
	_ = file.EstablishFolder(configPath)

	chainURL := "https://chainid.network/chains.json"
	chainsFile := filepath.Join(configPath, "chains.json")

	chainData, err := utils.DownloadAndStoreJSON[[]ChainListItem](chainURL, chainsFile, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	var chainList ChainList
	chainList.Chains = *chainData
	chainList.ChainsMap = make(map[int]*ChainListItem)

	for _, chain := range chainList.Chains {
		chainCopy := chain
		chainList.ChainsMap[chain.ChainID] = &chainCopy
	}

	return &chainList, nil
}
