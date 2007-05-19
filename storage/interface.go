package storage

type Storer interface {
	GetConnection()
}
/*
users
    ID unique
    Address

nfts
    ID unique
    ID creator
    ID owner
    ID collection
    tokenId
    nonce
    price
    link

transactions
    ID unique
    ID nft
    ID seller
    ID buyer
    tx hash

collections
    ID unique
    name
    tokenID

*/
