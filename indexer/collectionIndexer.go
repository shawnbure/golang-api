package indexer

// import (
// 	"encoding/base64"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/ENFT-DAO/youbei-api/storage"
// 	"github.com/btcsuite/btcutil/bech32"
// 	"gorm.io/gorm"
// )

// type CollectionIndexer struct {
// 	DeployerAddr string `json:"deployerAddr"`
// }

// func NewCollectionIndexer(deployerAddr string) (*CollectionIndexer, error) {
// 	return &CollectionIndexer{DeployerAddr: deployerAddr}, nil
// }

// func (ci *CollectionIndexer) StartWorker() {
// 	client := &http.Client{}
// 	for {
// 		time.Sleep(time.Second * 5)
// 		deployerStat, err := storage.GetDeployerStat(ci.DeployerAddr)
// 		if err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				deployerStat, err = storage.CreateDeployerStat(ci.DeployerAddr)
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					fmt.Println("something went wrong creating marketstat")
// 				}
// 			}
// 		}
// 		req, err := http.
// 			NewRequest("GET",
// 				fmt.Sprintf("https://devnet-api.elrond.com/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false",
// 					ci.DeployerAddr,
// 					deployerStat.LastIndex),
// 				nil)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			fmt.Println("error creating request for get nfts deployer")
// 			continue
// 		}
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			fmt.Println("error running request get nfts deployer")
// 			continue
// 		}

// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			resp.Body.Close()
// 			fmt.Println("error readall response get nfts deployer")
// 			continue
// 		}
// 		resp.Body.Close()
// 		if resp.Status != "200 OK" {
// 			fmt.Println("response not successful  get nfts deployer")
// 			continue
// 		}
// 		var ColResults []map[string]interface{}
// 		err = json.Unmarshal(body, &ColResults)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			fmt.Println("error unmarshal nfts deployer")
// 			continue
// 		}
// 		if len(ColResults) == 0 {
// 			time.Sleep(time.Minute * 2)
// 		}
// 		for _, colR := range ColResults {
// 			name := (colR["action"].(map[string]interface{}))["name"].(string)
// 			if name == "deployNFTTemplateContract" {
// 				results := (colR["results"].([]interface{}))
// 				result := results[0]
// 				data := result.(map[string]interface{})["data"].(string)
// 				decodedData64, _ := base64.StdEncoding.DecodeString(data)
// 				decodedData := strings.Split(string(decodedData64), "@")
// 				hexByte, err := hex.DecodeString(decodedData[2])
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					continue
// 				}
// 				byte32, err := bech32.ConvertBits(hexByte, 8, 5, true)
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					continue
// 				}
// 				bech32Addr, err := bech32.Encode("erd", byte32)
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					continue
// 				}

// 			}
// 		}
// 		newStat, err := storage.UpdateDeployerIndexer(deployerStat.LastIndex+1, ci.DeployerAddr)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			fmt.Println("error update deployer index nfts ")
// 			continue
// 		}
// 		if newStat.LastIndex <= deployerStat.LastIndex {
// 			fmt.Println("error something went wrong updating last index of deployer  ")
// 			continue
// 		}
// 	}
// }
