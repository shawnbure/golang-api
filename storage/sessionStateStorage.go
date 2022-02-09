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

func GetSessionStateByAccountId(accountId uint64, stateType uint64) (*entities.SessionState, error) {

	var sessionState entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&sessionState, "account_id = ? AND state_type = ?", accountId, stateType)
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

func DeleteSessionStateForAccountIdStateType(accountId uint64, stateType uint64) error {
	var sessionStates []entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(sessionStates, "account_id = ? AND state_type = ?", accountId, stateType)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
