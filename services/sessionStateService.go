package services

import (
	"fmt"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

type CreateUpdateSessionStateRequest struct {
	Address   string `json:"address"`
	StateType uint64 `json:"stateType"`
	JsonData  string `json:"jsonData"`
}

type RetreiveDeleteSessionStateRequest struct {
	Address   string `json:"address"`
	StateType uint64 `json:"stateType"`
}

func CreateSessionState(request *CreateUpdateSessionStateRequest) (*entities.SessionState, error) {

	fmt.Println("CreateSessionState 1")

	fmt.Println("request.AccountId: " + request.JsonData)

	sessionState := &entities.SessionState{
		ID:        0,
		Address:   request.Address,
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

func RetrieveSessionState(request *RetreiveDeleteSessionStateRequest) (*entities.SessionState, error) {

	sessionState, err := storage.GetSessionStateByAddressByStateType(request.Address, request.StateType)
	if err != nil {
		return nil, err
	}

	return sessionState, nil
}

func UpdateSessionState(sessionState *entities.SessionState, request *CreateUpdateSessionStateRequest) error {

	sessionState.Address = request.Address
	sessionState.StateType = request.StateType
	sessionState.JsonData = request.JsonData

	err := storage.UpdateSessionState(sessionState)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSessionState(request *RetreiveDeleteSessionStateRequest) (string, error) {

	err := storage.DeleteSessionStateForAddressStateType(request.Address, request.StateType)

	if err != nil {
		log.Debug("could not delete Sesion State", "err", err)
	}

	return "Successful Delete", err
}
