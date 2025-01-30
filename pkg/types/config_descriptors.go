package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

func (c *Config) ChainDescriptors() string {
	ret := []string{}
	for _, chain := range c.Chains {
		descr, err := Descriptor(c.General.DataFolder, &chain)
		if err != nil {
			ret = append(ret, err.Error())
		} else {
			ret = append(ret, descr.String())
		}
	}
	return strings.Join(ret, "\n")
}

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
	file.EstablishFolder(configPath)
	chainUrl := "https://chainid.network/chains.json"
	chainsFn := filepath.Join(configPath, "chains.json")
	if bytes, err := utils.DownloadAndStore(chainUrl, chainsFn, 24*time.Hour); err != nil {
		return &ChainList{}, err
	} else {
		var chainList ChainList
		err := json.Unmarshal(bytes, &chainList.Chains)
		if err != nil {
			return &ChainList{}, err
		}
		chainList.ChainsMap = make(map[int]*ChainListItem)
		for _, chain := range chainList.Chains {
			chainList.ChainsMap[chain.ChainID] = &chain
		}
		return &chainList, nil
	}
}

// ChainDescriptor represents the configuration of a single chain in the configurion file's template.
type ChainDescriptor struct {
	Chain          string `json:"chain"`
	ChainID        string `json:"chainId"`
	RemoteExplorer string `json:"remoteExplorer"`
	RpcProvider    string `json:"rpcProvider"`
	Symbol         string `json:"symbol"`
}

func (c *ChainDescriptor) String() string {
	bytes, _ := json.Marshal(c) // MarshalIndent(c, "", "  ")
	return string(bytes)
}

func Descriptor(configFolder string, ch *Chain) (ChainDescriptor, error) {
	return ChainDescriptor{
		Chain:          ch.Name,
		ChainID:        fmt.Sprintf("%d", ch.ChainID),
		Symbol:         "SYM",
		RpcProvider:    ch.RPCs[0],
		RemoteExplorer: "https://etherscan.io",
	}, nil
}

// func (c *Config) ChainDescriptors2() string {
// 	// 	dataFn := filepath.Join(c.General.DataFolder, "chains.json")
// 	// 	chainData := file.AsciiFileToString(dataFn)
// 	// 	if !file.FileExists(dataFn) || len(chainData) == 0 {
// 	// 		chainData = `{
// 	//   "mainnet": {
// 	//     "chain": "mainnet",
// 	//     "chainId": "1",
// 	//     "remoteExplorer": "https://etherscan.io",
// 	//     "symbol": "ETH"
// 	//   }
// 	// }
// 	// `
// 	// 	}

// 	// 	chainDescrs := make(map[string]ChainDescriptor)
// 	// 	if err := json.Unmarshal([]byte(chainData), &chainDescrs); err != nil {
// 	// 		return err.Error()
// 	// 	}

// 	tmpl, err := template.New("chainConfigTmpl").Parse(`  [chains.{{.Chain}}]
//     chain = "{{.Chain}}"
//     chainId = "{{.ChainID}}"
//     remoteExplorer = "{{.RemoteExplorer}}"
//     rpcProvider = "{{.RpcProvider}}"
//     symbol = "{{.Symbol}}"`)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	ret := []string{}
// 	for _, ch := range c.Chains {
// 		if chainConfig, ok := c.ChainList.ChainsMap[ch.ChainId]; ok {
// 			chainConfig.RpcProvider = ch.RPCs[0]
// 			var buf bytes.Buffer
// 			if err = tmpl.Execute(&buf, &chainConfig); err != nil {
// 				return err.Error()
// 			}

// 			ret = append(ret, buf.String())
// 		} else {
// 			ret = append(ret, "  # "+chain.Name+" is not supported")
// 		}
// 	}

// 	sort.Slice(ret, func(i, j int) bool {
// 		return strings.Compare(ret[i], ret[j]) < 0
// 	})

// 	return "\n" + strings.Join(ret, "\n")
// }
