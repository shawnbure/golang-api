package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type ActivityLogsList struct {
	Activities []entities.Activity `json:"activities"`
	TotalCount int64               `json:"total"`
}
