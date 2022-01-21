package entities

type Whitelist struct {
	ID           uint64 `gorm:"primaryKey" json:"id"` //PK
	CollectionID uint64 `json:"collectionId"`         //FK to the
	Address      string `json:"address"`
	Amount       uint64 `json:"amount"`
	Type         uint64 `json:"type"`
	CreatedAt    uint64 `json:"createdAt"`
	ModifiedAt   uint64 `json:"modifiedAt"`
}

const (
	WhitelistType_none     = 0
	WhitelistType_buy      = 1
	WhitelistType_mint     = 2
	WhitelistType_buy_mint = 3
)

//test

/*

	1. Add a whitelist check in the collection table

	2. When user click the buy button per a collection,
		it check if the LoggedIn account 'Address' is in the whitelist table.
			- if not, show message 'Sorry, you are not part of the whitelist'
		if it' in there, check if the WhiteListType is allowed to 'buy'
			- if not, show message, 'Sorry, you are not allow to buy'
		if allowed to buy, then check the 'Amount' to see if it is not zero
			- if it is zero, show message, 'Sorry, you already bought your allocated whitelist amount'
		if it's not zero, then deduct it by 1 and proceed on with the 'buy' process


*/
