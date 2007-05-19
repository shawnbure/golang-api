package dtos

import (
	"encoding/json"
	"fmt"
	"github.com/erdsea/erdsea-api/data/entities"
	"gorm.io/datatypes"
	"testing"
)

var a = entities.Asset{
	Attributes: datatypes.JSON(`{"x": "y", "y": "x"}`),
}

func TestTest(t *testing.T) {
	var rr map[string]string
	_ = json.Unmarshal(a.Attributes, &rr)

	fmt.Println(rr)

	s, _ := json.Marshal(rr)

	fmt.Println(string(s))
}
