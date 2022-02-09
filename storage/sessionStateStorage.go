package storage

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddSessionState(sesionState *entities.SessionState) error {

	fmt.Println("AddSessionState 1")

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	fmt.Println("AddSessionState 2")

	txCreate := database.Create(&sesionState)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	fmt.Println("AddSessionState3")

	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	fmt.Println("AddSessionState 4")

	return nil
}

func GetSessionStateByAddressByStateType(address string, stateType uint64) (*entities.SessionState, error) {

	var sessionState entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&sessionState, "address = ? AND state_type = ?", address, stateType)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &sessionState, nil
}

func UpdateSessionState(sessionState *entities.SessionState) error {
	database, err := GetDBOrError()

	if err != nil {
		return err
	}

	txCreate := database.Save(&sessionState)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteSessionStateForAddressStateType(address string, stateType uint64) error {
	var sessionStates []entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(sessionStates, "address = ? AND state_type = ?", address, stateType)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
