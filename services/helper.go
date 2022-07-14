package services

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	urlp "net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/proxier"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ConvertFilterToQuery converts a querystring conversion filter to a sql where clause
func ConvertFilterToQuery(tableName string, filter string) (string, []interface{}, error) { //Field|Value|Operator;AND/OR;...
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in %v", r)
		}
	}()
	var stm string
	var values []interface{}
	if filter == "" {
		return stm, values, nil
	}
	clauses := strings.Split(filter, ";")
	for _, c := range clauses {
		if c == "OR" || c == "AND" {
			stm += " " + c + " "
			continue
		}
		params := strings.Split(c, "|")
		// We must have only 3 length of params. field|value|operator
		if len(params) != 3 {
			fmt.Printf("we have wrong params structure '%s' in '%s'", c, filter)
			continue
		}
		field := params[0]
		value := params[1]
		operator := params[2]
		prefix := tableName
		subObjects := strings.Split(field, ".")

		if len(subObjects) > 1 {
			prefix = "\"" + subObjects[0] + "\""
			field = strings.Join(subObjects[1:], ".")
		}

		var query string

		if operator == "BETWEEN" {
			ranges := strings.Split(value, "AND")
			if len(ranges) != 2 {
				err := errors.New("bad given between range")
				return "", nil, err
			}
			query = prefix + "." + field + " BETWEEN " + "?" + " AND " + "?"
			values = append(values, ranges[0])
			values = append(values, ranges[1])
		} else {
			query = prefix + "." + field + " " + operator + " " + "?"
			values = append(values, value)
		}
		stm += query
	}
	return stm, values, nil
}

func ConvertSortToQuery(tableName string, sort string) (string, []interface{}, error) { //Field|(asc/desc);...
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in %v", r)
		}
	}()
	var stm string
	var values []interface{}
	if sort == "" {
		return stm, values, nil
	}
	clauses := strings.Split(sort, ";")
	if len(clauses) > 1 {
		return stm, values, nil
	}

	for index, c := range clauses {
		params := strings.Split(c, "|")
		// We must have only 3 length of params. field|value|operator
		if len(params) != 2 {
			fmt.Printf("we have wrong params structure '%s' in '%s'", c, sort)
			continue
		}
		field := params[0]
		value := params[1]
		if !(value == "asc" || value == "desc") {
			fmt.Printf("we have wrong params structure '%s' in '%s'", c, sort)
			continue
		}

		prefix := tableName
		subObjects := strings.Split(field, ".")

		if len(subObjects) > 1 {
			prefix = "\"" + subObjects[0] + "\""
			field = strings.Join(subObjects[1:], ".")
		}

		var query string

		query = prefix + "." + field + " " + "%s"
		values = append(values, value)
		if index < len(clauses)-1 {
			query += ", "
		}
		stm += query
	}
	return stm, values, nil
}

func ConvertAttributeFilterToQuery(filter string) ([][]string, error) { //Field|(asc/desc);...

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in %v", r)
		}
	}()

	var values [][]string

	if filter == "" {
		return values, nil
	}

	clauses := strings.Split(filter, ";")
	for _, c := range clauses {
		params := strings.Split(c, "|")

		field := params[0]
		value := params[1]

		values = append(values, []string{field, value})
	}
	return values, nil
}

var timeElapsed = 0

func GetResponse(url string) ([]byte, error) {
	for int(time.Now().UnixMilli())-timeElapsed < 1000 && rand.Int63n(10000) < 7000 {
		time.Sleep(time.Millisecond * 100)
	}
	timeElapsed = int(time.Now().UnixMilli())
	ipStr := proxier.GetCurrentIP()
	var client http.Client
	client.Timeout = time.Second * 10
	if ipStr != "" {
		proxyUrl, err := urlp.Parse(ipStr)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client.Transport = transport
	}
	req, err := http.
		NewRequest("GET", url,
			nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.Status != "200 OK" {
		if strings.Contains(resp.Status, "429") {
			proxier.ChangeIP()
		}
		return nil, fmt.Errorf("status %s %s", resp.Status, req.URL.RawPath)
	}
	return body, nil
}

func GetTransactionBC(hash string, api string) (entities.TransactionBC, error) {

	reqUrl := fmt.Sprintf("%s/transactions/%s",
		hash,
		api)
	body, err := GetResponse(reqUrl)
	if err != nil {
		zlog.Error(err.Error())
	}
	var tx entities.TransactionBC
	err = json.Unmarshal(body, &tx)
	return tx, err
}

// ConvertAttributeFilterToJsonQuery converts a querystring conversion attribute filter to a sql jsonb where clause
func ConvertAttributeFilterToJsonQuery(tableName string, filter string) (string, []interface{}, error) { //Field|Value|Operator;AND;...
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in %v", r)
		}
	}()
	var stm string
	var values []interface{}
	if filter == "" {
		return stm, values, nil
	}
	clauses := strings.Split(filter, ";")
	for _, c := range clauses {
		if c == "OR" || c == "AND" {
			stm += " " + c + " "
			continue
		}
		params := strings.Split(c, "|")
		// We must have only 3 length of params. field|value|operator
		if len(params) != 3 {
			fmt.Printf("we have wrong params structure '%s' in '%s'", c, filter)
			continue
		}
		field := params[0]
		value := params[1]
		operator := params[2]
		prefix := tableName
		subObjects := strings.Split(field, ".")

		if len(subObjects) > 1 {
			prefix = "\"" + subObjects[0] + "\""
			field = strings.Join(subObjects[1:], ".")
		}

		var query string

		if operator == "BETWEEN" {
			ranges := strings.Split(value, "AND")
			if len(ranges) != 2 {
				err := errors.New("bad given between range")
				return "", nil, err
			}
			query = prefix + "." + field + " BETWEEN " + "?" + " AND " + "?"
			values = append(values, ranges[0])
			values = append(values, ranges[1])
		} else {
			query = prefix + "." + field + " " + operator + " " + "?"
			values = append(values, value)
		}
		stm += query
	}
	return stm, values, nil
}

func TurnIntoBigInt18Dec(num int64) *big.Int {
	bigNum := big.NewInt(num)
	bigNum = bigNum.Mul(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil), bigNum)
	return bigNum
}

func Zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(256)
	return r
}

func Mul(a, b *big.Float) *big.Float {
	return Zero().Mul(a, b)
}
func Pow(a *big.Float, e uint64) *big.Float {
	result := Zero().Copy(a)
	for i := uint64(0); i < e-1; i++ {
		result = Mul(result, a)
	}
	return result
}
func TurnIntoBigFloat18Dec(num float64) *big.Float {
	bigNum := big.NewFloat(num)
	bigNum = Mul(bigNum, Pow(big.NewFloat(10), 18))
	return bigNum
}

func TurnIntoBigInt8Dec(num int64) *big.Int {
	bigNum := big.NewInt(num)
	bigNum = bigNum.Mul(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(8), nil), bigNum)
	return bigNum
}

func TurnIntoBigIntNDec(num int64, decimal int64) *big.Int {
	bigNum := big.NewInt(num)
	bigNum = bigNum.Mul(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(decimal), nil), bigNum)
	return bigNum
}

func TurnBigIntoBigIntNDec(num *big.Int, decimal int64) *big.Int {
	num = num.Mul(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(decimal), nil), num)
	return num
}

func TurnBigFloatoBigFloatNDec(num *big.Float, decimal int64) (*big.Float, bool) {
	powerString := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimal), nil).String()
	powerFloat, ok := big.NewFloat(0).SetString(powerString)
	if !ok {
		return nil, ok
	}
	num = num.Mul(powerFloat, num)
	return num, true
}

func ConvertBigFloatToFloat(num string) (*big.Float, error) {
	bigFloat, ok := big.NewFloat(0).SetString(num)
	if !ok {
		return nil, fmt.Errorf("failed to convert big float to float")
	}
	return bigFloat, nil
}
func RandomName() string {
	guid := xid.New()
	return guid.String()
}
func ConvertHexToBehc32(addrHex string) (string, error) {
	hexByte, err := hex.DecodeString(addrHex)
	if err != nil {
		return "", err
	}
	byte32, err := bech32.ConvertBits(hexByte, 8, 5, true)
	if err != nil {
		return "", err
	}
	bech32Addr, err := bech32.Encode("erd", byte32)
	if err != nil {
		return "", err
	}
	return bech32Addr, nil
}

func IndexTokenAttribute(tokenIdentifier string, nonceStr string, api string) (*entities.TokenBC, error) {
tokenLoop:
	token, err := getTokenBC(tokenIdentifier, nonceStr, api)
	if err != nil {
		return nil, err
	}
	colObj, err := storage.GetCollectionByTokenId(token.Collection)
	if err != nil {
		return nil, err
	}
	imageURI, attributeURI := GetTokenUris(token)
	// nonce10Str := strconv.FormatUint(token.Nonce, 10)
	// nonceStr := strconv.FormatUint(token.Nonce, 16)
	// if len(nonceStr)%2 != 0 {
	// 	nonceStr = "0" + nonceStr
	// }
	// Convert URI to elrond url for faster retreive
	if strings.Contains(api, "devnet") {
		imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://devnet-media.elrond.com/nfts/asset/", 1)
	} else {
		imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
		imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
		imageURI = strings.Replace(imageURI, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)
	}
	youbeiMeta := strings.Replace(attributeURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
	youbeiMeta = strings.Replace(youbeiMeta, "https://media.youbei.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
	youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
	youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
	youbeiMeta = strings.Replace(youbeiMeta, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)

	var url string = youbeiMeta
	var attrbs []byte
	metadataJSON := make(map[string]interface{})

	if token.Attributes == "" {
		attrbs, err = GetResponse(url)
		if err != nil {
			zlog.Error(err.Error(), zap.String("url", string(url)), zap.Strings("URIS", token.URIs), zap.String("collection", token.Collection), zap.String("attributes", token.Attributes), zap.String("identifier", token.Identifier), zap.Any("media", token.Media), zap.Any("Metadata", token.Metadata))
		}
		err = json.Unmarshal(attrbs, &metadataJSON)
		if err != nil {
			zlog.Error(err.Error(), zap.String("url", string(url)), zap.Strings("URIS", token.URIs), zap.String("collection", token.Collection), zap.String("attributes", token.Attributes), zap.String("identifier", token.Identifier), zap.Any("media", token.Media), zap.Any("Metadata", token.Metadata))
		}
	}
	if !strings.Contains(youbeiMeta, "http") {
		youbeiMeta = ""
	}
	var attributes datatypes.JSON
	if token.Attributes != "" {
		if _, ok := metadataJSON["attributes"]; !ok {
			attributesStr, err := base64.StdEncoding.DecodeString(token.Attributes)
			if strings.Contains(string(attributesStr), ".json") {
				if strings.Contains(string(attributesStr), "metadata:") {
					attributeParts := strings.Split(string(attributesStr), ";")
					for _, part := range attributeParts {
						if strings.Contains("metadata:", part) {
							part = part[9:]
							// attributesStr = []byte(strings.Replace(string(attributesStr), "metadata:", "", 1))
							url = (`https://media.elrond.com/nfts/asset/` + string(part))
							attrbs, err := GetResponse(url)
							if err != nil {
								zlog.Error(err.Error(), zap.String("url", string(url)), zap.Strings("URIS", token.URIs), zap.String("collection", token.Collection), zap.String("attributes", token.Attributes), zap.String("identifier", token.Identifier), zap.Any("media", token.Media), zap.Any("Metadata", token.Metadata))
							}

							metadataJSON = make(map[string]interface{})
							err = json.Unmarshal(attrbs, &metadataJSON)
							if err != nil {
								zlog.Error(err.Error(), zap.String("url", string(url)), zap.Strings("URIS", token.URIs), zap.String("collection", token.Collection), zap.String("attributes", token.Attributes), zap.String("identifier", token.Identifier), zap.Any("media", token.Media), zap.Any("Metadata", token.Metadata))
							}
							attributesBytes, err := json.Marshal(metadataJSON["attributes"])
							if err != nil {
								zlog.Error(err.Error())
								attributesBytes = []byte{}
							}
							err = json.Unmarshal(attributesBytes, &attributes)
							if err != nil {
								zlog.Error(err.Error())
							}
						}
					}

				}
			}
			if attributes.String() == "" {
				resultStr := `[`
				if err != nil {
					zlog.Error("attribute decoding failed", zap.Error(err), zap.String("attribute", token.Attributes))
				} else {
					attrbutesParts := strings.Split(string(attributesStr), ";")
					var prefix string = ""
					for i, ap := range attrbutesParts {
						if i != 0 {
							prefix = ","
						}
						traitKeyValue := strings.Split(ap, ":")
						if len(traitKeyValue) < 2 {
							continue
						}
						resultStr = resultStr + prefix + `{"` + traitKeyValue[0] + `":"` + traitKeyValue[1] + `"}`
					}
					resultStr = resultStr + "]"
				}
				attributes = datatypes.JSON(resultStr)
			}
		}

	} else {
		attributesBytes, err := json.Marshal(metadataJSON["attributes"])
		if err != nil {
			zlog.Error(err.Error())
			attributesBytes = []byte{}
		}
		err = json.Unmarshal(attributesBytes, &attributes)
		if err != nil {
			zlog.Error(err.Error())
		}
	}

	//get owner of token from database TODO
	if token.Owner == "" {
		tokenRes, err := GetResponse(fmt.Sprintf("%s/nfts/%s", api, token.Identifier))
		if err != nil {
			zlog.Error("CRITICAL can't get nft data", zap.Error(err))
			if strings.Contains(err.Error(), "deadline") {
				goto tokenLoop
			}
		}
		json.Unmarshal(tokenRes, &token)
	}
	if token.Owner == "" {
		token.Owner = token.Creator
	}
	acc, err := storage.GetAccountByAddress(token.Owner)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			name := RandomName()
			acc = &entities.Account{
				Address: token.Owner,
				Name:    name,
			}
			err := storage.AddAccount(acc)
			if err != nil {
				if !strings.Contains(err.Error(), "duplicate") {
					zlog.Error("CRITICAL can't create user", zap.Error(err))

				} else {
					acc, err = storage.GetAccountByAddress(token.Owner)
					if err != nil {
						zlog.Error("CRITICAL can't get user", zap.Error(err))

					}
				}
			}
		}
	}
	//try get token from database TODO
	dbToken, err := storage.GetTokenByTokenIdAndNonce(token.Collection, token.Nonce)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			zlog.Error("getTokenByTokenIdAndNonce_error", zap.Error(err))
		} else {

		}
	}
	var js interface{}
	if json.Unmarshal(attributes, &js) != nil {
		attributes = []byte("{}")
	}
	if !utf8.ValidString(youbeiMeta) {
		youbeiMeta = ""
	}
	if dbToken == nil {
		dbToken = &entities.Token{
			TokenID:      token.Collection,
			MintTxHash:   "",
			CollectionID: colObj.ID,
			Nonce:        token.Nonce,
			NonceStr:     nonceStr,
			MetadataLink: strings.ToValidUTF8(youbeiMeta, " "),
			ImageLink:    imageURI,
			TokenName:    token.Name,
			Attributes:   attributes,
			OwnerID:      acc.ID,
			PriceString:  "0",
			PriceNominal: 0,
		}
	}
	err = storage.AddOrUpdateToken(&entities.Token{
		TokenID:      token.Collection,
		MintTxHash:   dbToken.MintTxHash,
		CollectionID: colObj.ID,
		Nonce:        token.Nonce,
		NonceStr:     nonceStr,
		MetadataLink: strings.ToValidUTF8(youbeiMeta, " "),
		ImageLink:    imageURI,
		TokenName:    token.Name,
		Attributes:   attributes,
		OwnerID:      dbToken.OwnerID,
		PriceString:  dbToken.PriceString,
		PriceNominal: dbToken.PriceNominal,
	})
	if err != nil {
		zlog.Error("BADERR", zap.Error(err), zap.Any("token", token))
		if err == gorm.ErrRegistered {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &token, nil
}
