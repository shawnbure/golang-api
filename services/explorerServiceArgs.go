package services

import "github.com/ENFT-DAO/youbei-api/data/entities"

type GetAllExplorerTokensArgs struct {
	LastTimestamp    int64
	Limit            int
	CurrentPage      int
	NextPage         int
	Filter           *entities.QueryFilter
	CollectionFilter *entities.QueryFilter
	SortOptions      *entities.SortOptions
	Attributes       [][]string
}
