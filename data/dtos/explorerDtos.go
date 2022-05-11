package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type ExplorerTokenList struct {
	Tokens     []entities.TokenExplorer `json:"tokens"`
	TotalCount int64                    `json:"total"`
	MinPrice   float64                  `json:"min_price"`
	MaxPrice   float64                  `json:"max_price"`
}
