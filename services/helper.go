package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	urlp "net/url"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/proxier"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/rs/xid"
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
