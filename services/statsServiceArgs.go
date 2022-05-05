package services

import "github.com/ENFT-DAO/youbei-api/data/entities"

type GetAllTransactionsWithPaginationArgs struct {
	LastTimestamp int64
	Limit         int
	CurrentPage   int
	NextPage      int
	Filter        *entities.QueryFilter
}
