package storage

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"gorm.io/datatypes"
)

var toJson = func(m map[string]string) datatypes.JSON {
	attrs, _ := json.Marshal(m)
	return attrs
}

var _accounts = []entities.Account{
	{
		Address:       "erd1tfpkdly5zmsnst78pnr7fpfy6qpmmuvhdlff870lpufdaz9vnsrsnqen8j",
		Name:          "bouncing banana",
		Description:   "fuck-off",
		Website:       "https://github.com/",
		TwitterLink:   "https://twitter.com/disclosetv",
		InstagramLink: "https://www.instagram.com/",
		CreatedAt:     uint64(time.Now().Unix()),
	},
	{
		Address:       "erd1k62knymy0ypa8wusut7254q3tvnfy73kstvcggljjd7z44vsdz0qe5r59x",
		Name:          "vip cookie 90",
		Description:   "fuck-off twice",
		Website:       "https://github.com/",
		TwitterLink:   "https://twitter.com/Snowden",
		InstagramLink: "https://www.instagram.com/",
		CreatedAt:     uint64(time.Now().Unix()),
	},
}

var _collections = []entities.Collection{
	{
		Name:          "Apes",
		TokenID:       "APE-abcdef",
		Description:   "The best apes",
		Website:       "https://github.com/",
		DiscordLink:   "https://discord.com",
		TwitterLink:   "https://twitter.com/BoredApeYC",
		InstagramLink: "https://www.instagram.com",
		TelegramLink:  "https://telegram.com",
		CreatedAt:     uint64(time.Now().Unix()),
		Priority:      25,
	},
	{
		Name:          "Women",
		TokenID:       "WMEN-abcdef",
		Description:   "The best women",
		Website:       "https://github.com/",
		DiscordLink:   "https://discord.com",
		TwitterLink:   "https://twitter.com/WOW_Accesorios",
		InstagramLink: "https://www.instagram.com",
		TelegramLink:  "https://telegram.com",
		CreatedAt:     uint64(time.Now().Unix()),
		Priority:      50,
	},
}

var _tokens = []entities.Token{
	{
		TokenID:          "APE-abcdef",
		Nonce:            1,
		PriceString:      "100000000000000000000",
		PriceNominal:     100,
		RoyaltiesPercent: 200,
		MetadataLink:     "https://galacticapes.mypinata.cloud/ipfs/QmcX6g2xXiFP5j1iAfXREuP9EucRRpuMCAnoYaVYjtrJeK/1",
		CreatedAt:        uint64(time.Now().Unix()),
		Status:           entities.List,
		Attributes: toJson(map[string]string{
			"background": "azure",
			"face":       "grey",
		}),
		TokenName: "Galactic apes",
		ImageLink: "https://galacticapes.mypinata.cloud/ipfs/QmPqKt7guhrCNS6DWy7gNeyR9ia7UgijVj8evWcUjFiQrc/1.png",
		Hash:      "",
	},
	{
		TokenID:          "APE-abcdef",
		Nonce:            2,
		PriceString:      "1000000000000000000",
		PriceNominal:     1,
		RoyaltiesPercent: 200,
		MetadataLink:     "https://galacticapes.mypinata.cloud/ipfs/QmcX6g2xXiFP5j1iAfXREuP9EucRRpuMCAnoYaVYjtrJeK/2",
		CreatedAt:        uint64(time.Now().Unix()),
		Status:           entities.List,
		Attributes: toJson(map[string]string{
			"background": "green",
			"face":       "blue",
			"eye":        "glasses",
		}),
		TokenName: "Galactic apes",
		ImageLink: "https://galacticapes.mypinata.cloud/ipfs/QmPqKt7guhrCNS6DWy7gNeyR9ia7UgijVj8evWcUjFiQrc/2.png",
		Hash:      "",
	},
}

var _txTemplate = entities.Transaction{
	Hash:         "hash1",
	Type:         entities.ListToken,
	PriceNominal: 100,
}

func SeedDatabase(cfg config.DatabaseConfig) {
	Connect(cfg)

	err := addAccounts()
	if err != nil {
		panic(err)
	}

	err = addCollections()
	if err != nil {
		panic(err)
	}

	err = addTokens()
	if err != nil {
		panic(err)
	}

	err = addTxs()
	if err != nil {
		panic(err)
	}
}

func addAccounts() error {
	for index := range _accounts {
		if err := AddAccount(&_accounts[index]); err != nil {
			return err
		}
	}

	return nil
}

func addCollections() error {
	for index := range _collections {
		_collections[index].CreatorID = _accounts[rand.Intn(len(_accounts))].ID
		if err := AddCollection(&_collections[index]); err != nil {
			return err
		}
	}

	return nil
}

func addTokens() error {
	for index := range _tokens {
		_tokens[index].OwnerId = _accounts[rand.Intn(len(_accounts))].ID
		_tokens[index].CollectionID = _collections[rand.Intn(len(_collections))].ID
		if err := AddToken(&_tokens[index]); err != nil {
			return err
		}
	}

	return nil
}

func addTxs() error {
	for i := 1; i < 20; i++ {
		var txType entities.TxType
		if i%3 == 0 {
			txType = entities.ListToken
		}
		if i%3 == 1 {
			txType = entities.BuyToken
		}
		if i%3 == 2 {
			txType = entities.WithdrawToken
		}

		tx := _txTemplate
		randToken := _tokens[rand.Intn(len(_tokens))]

		tx.Type = txType
		tx.TokenID = randToken.ID
		tx.PriceNominal = randToken.PriceNominal
		tx.CollectionID = randToken.CollectionID
		tx.BuyerID = _accounts[rand.Intn(len(_accounts))].ID
		tx.SellerID = _accounts[rand.Intn(len(_accounts))].ID

		if err := AddTransaction(&tx); err != nil {
			return err
		}
	}

	return nil
}
