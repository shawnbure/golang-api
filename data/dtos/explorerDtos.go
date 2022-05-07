package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type ExplorerTokenList struct {
	Tokens []entities.TokenExplorer `json:"tokens"`
}
