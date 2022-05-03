package services

import "github.com/ENFT-DAO/youbei-api/data/entities"

type GetAllActivityArgs struct {
	LastTimestamp int64
	Limit         int
	Filter        *entities.QueryFilter
}
