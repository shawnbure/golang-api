package storage

import (
	"fmt"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"gorm.io/gorm"
)

// Refresh the session state (if not existing, create new, else just do an update/save)
func RefreshCreateOrUpdateSessionState(address string, stateType uint64, jsonData string) error {
	var sessionState entities.SessionState

	//get db
	database, err := GetDBOrError()

	if err != nil {
		return err
	}

	//do a query to see address and state type exists
	txRead := database.Find(&sessionState, "address = ? AND state_type = ?", address, stateType)

	if txRead.Error != nil {
		return txRead.Error
	}

	if txRead.RowsAffected == 0 {

		// -------- CREATE PROCESS --------

		//create new session state
		sessionStateNew := &entities.SessionState{
			ID:        0,
			Address:   address,
			StateType: stateType,
			JsonData:  jsonData,
			CreatedAt: uint64(time.Now().Unix()),
		}

		//create
		txCreate := database.Create(&sessionStateNew)

		//check for errors
		if txCreate.Error != nil {
			return txCreate.Error
		}
	} else {

		// -------- UPDATE PROCESS --------

		//set the json data
		sessionState.JsonData = jsonData

		if stateType == entities.SessionState_type_create_collection {

			//check the JSON and see the steps - default
			fmt.Println("inside")
		}

		//save / update
		txCreate := database.Save(&sessionState)

		//check for error
		if txCreate.Error != nil {
			return txCreate.Error
		}
	}

	return nil //successful so return nil for ERROR
}

//create the session state
func CreateSessionState(sesionState *entities.SessionState) error {

	//get db
	database, err := GetDBOrError()

	if err != nil {
		return err
	}

	//create
	txCreate := database.Create(&sesionState)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	//error checks
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil //succesful
}

//get session state by addresss and by state type
func GetSessionStateByAddressByStateType(address string, stateType uint64) (*entities.SessionState, error) {

	var sessionState entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	//query
	txRead := database.Find(&sessionState, "address = ? AND state_type = ?", address, stateType)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	//if no records are found
	if txRead.RowsAffected == 0 {

		//create a empty one and return it
		sessionStateEmpty := &entities.SessionState{
			ID:        0,
			Address:   address,
			StateType: stateType,
			JsonData:  "{ \"step\": 1, \"tokenID\": \"Empty\", \"scAddress\": \"Empty\", \"price\": 0 }",
			CreatedAt: uint64(time.Now().Unix()),
		}

		return sessionStateEmpty, nil
	}

	//found one and return it
	return &sessionState, nil
}

//update the session state
func UpdateSessionState(sessionState *entities.SessionState) error {

	//get DB
	database, err := GetDBOrError()

	if err != nil {
		return err
	}

	//save / update the session state
	txCreate := database.Save(&sessionState)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	//check if it's successful
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

//delete session state by address and state type
func DeleteSessionStateForAddressStateType(address string, stateType uint64) error {
	var sessionStates []entities.SessionState

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	//query delete
	txCreate := database.Delete(sessionStates, "address = ? AND state_type = ?", address, stateType)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
