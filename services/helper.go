package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

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

func GetResponse(url string) ([]byte, error) {
	var client http.Client
	client.Timeout = time.Second * 5
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
		return nil, fmt.Errorf(resp.Status, req.URL.RawPath)
	}
	return body, nil
}
func TurnIntoBigInt18Dec(num int64) *big.Int {
	bigNum := big.NewInt(num)
	bigNum = bigNum.Mul(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil), bigNum)
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

func RandomName() string {
	guid := xid.New()
	return guid.String()
}
