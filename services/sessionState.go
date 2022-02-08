package services

import (
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

type CreateSessionStateRequest struct {
	AccountId uint64 `json:"accountId"`
	StateType uint64 `json:"stateType"`
	JSONData  string `json:"jsonData"`
}

func CreateSessionState(request *CreateSessionStateRequest) (*entities.SessionState, error) {

	sessionState := &entities.SessionState{
		ID:        0,
		AccountID: request.AccountId,
		StateType: request.StateType,
		JSONData:  request.JSONData,
		CreatedAt: uint64(time.Now().Unix()),
	}

	return sessionState, nil
}
