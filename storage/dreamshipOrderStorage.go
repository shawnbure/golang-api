package storage

import (
	"fmt"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddOrUpdateUserPayment(userPayment entities.UserPayments) error{
	db, err := GetDBOrError()
	if err != nil {
		return err
	}
	result := db.Model(&entities.UserPayments{}).
				Where("tx_hash = ?", userPayment.TxHash).
				Updates(entities.UserPayments {
					Status: userPayment.Status,
				})
	if result.RowsAffected == 0{
		db.Create(&userPayment)
	}
	return nil
}

func AddOrUpdateOrderItem(orderItem entities.UserOrders) error{
	db, err := GetDBOrError()
	if err != nil {
		return err
	}
	result := db.Model(&entities.UserOrders{}).
				Where("order_id = ?", orderItem.OrderId).
				Updates(entities.UserOrders {
					Amount: orderItem.Amount,
					OrderStatus: orderItem.OrderStatus,
					CheckoutStatus: orderItem.CheckoutStatus,
					PaymentMethod: orderItem.PaymentMethod,
					OrderId: orderItem.OrderId,
				})
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0{
			db.Create(&orderItem)
	}
	return nil
}

func RetrievesUserOrders(userAddress string) ([]entities.UserOrders, error) {
	db, err := GetDBOrError()
	if err != nil {
		return nil, err
	}
	var result []entities.UserOrders
	txRead := db.Model(&entities.UserOrders{}).Where("user_address = ?", userAddress).Find(&result)
	if txRead.Error != nil {
		fmt.Println(txRead.Error)
		return nil, txRead.Error
	}
	return result, nil
}

func RetrievesAnOrders(userAddress string, orderId string) ([]entities.UserOrders, error) {
	db, err := GetDBOrError()
	if err != nil {
		return nil, err
	}
	var result []entities.UserOrders
	txRead := db.Model(&entities.UserOrders{}).Where("user_address = ? AND order_id = ?", userAddress, orderId).Find(&result)
	if txRead.Error != nil {
		fmt.Println(txRead.Error)
		return nil, txRead.Error
	}
	return result, nil
}