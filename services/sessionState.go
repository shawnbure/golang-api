package services

import (
	"fmt"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

type CreateSessionStateRequest struct {
	AccountId uint64 `json:"accountId"`
	StateType uint64 `json:"stateType"`
	JsonData  string `json:"jsonData"`
}

func CreateSessionState(request *CreateSessionStateRequest) (*entities.SessionState, error) {

	fmt.Println("CreateSessionState 1")

	//string(request.AccountId)

	fmt.Println("request.AccountId: " + request.JsonData)

	sessionState := &entities.SessionState{
		ID:        0,
		AccountID: request.AccountId,
		StateType: request.StateType,
		JsonData:  request.JsonData,
		CreatedAt: uint64(time.Now().Unix()),
	}

	fmt.Println("CreateSessionState 2")

	err := storage.AddSessionState(sessionState)
	if err != nil {
		return nil, err
	}

	return sessionState, nil
}

func DeleteSessionState(accountId uint64, stateType uint64) {
	err := storage.DeleteSessionStateForAccountIdStateType(accountId, stateType)

	if err != nil {
		log.Debug("could not delete proffers for token", "err", err)
	}
}
