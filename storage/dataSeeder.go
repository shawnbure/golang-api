package storage

import (
	"encoding/json"
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
		Priority:      50,

		CreatorID: 1,
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
		Priority:      25,

		CreatorID: 1,
	},
}

var _tokens = []entities.Token{
	{
		TokenID:          "APE-abcdef",
		Nonce:            1,
		PriceString:      "100",
		PriceNominal:     100,
		RoyaltiesPercent: 200,
		MetadataLink:     "https://galacticapes.mypinata.cloud/ipfs/QmcX6g2xXiFP5j1iAfXREuP9EucRRpuMCAnoYaVYjtrJeK",
		CreatedAt:        uint64(time.Now().Unix()),
		Listed:           true,
		Attributes: toJson(map[string]string{
			"background": "azure",
			"face":       "grey",
		}),
		TokenName: "Galactic apes",
		ImageLink: "https://galacticapes.mypinata.cloud/ipfs/QmPqKt7guhrCNS6DWy7gNeyR9ia7UgijVj8evWcUjFiQrc/1.png",
		Hash:      "",

		OwnerId:      2,
		CollectionID: 2,
	},
	{
		TokenID:          "APE-abcdef",
		Nonce:            2,
		PriceString:      "0",
		PriceNominal:     0,
		RoyaltiesPercent: 200,
		MetadataLink:     "https://galacticapes.mypinata.cloud/ipfs/QmcX6g2xXiFP5j1iAfXREuP9EucRRpuMCAnoYaVYjtrJeK",
		CreatedAt:        uint64(time.Now().Unix()),
		Listed:           false,
		Attributes: toJson(map[string]string{
			"background": "green",
			"face":       "blue",
			"eye":        "glasses",
		}),
		TokenName: "Galactic apes",
		ImageLink: "https://galacticapes.mypinata.cloud/ipfs/QmPqKt7guhrCNS6DWy7gNeyR9ia7UgijVj8evWcUjFiQrc/2.png",
		Hash:      "",

		OwnerId:      1,
		CollectionID: 2,
	},
}

var _txs = []entities.Transaction{
	{
		Hash:         "hash1",
		Type:         entities.ListToken,
		PriceNominal: 100,
		SellerID:     2,
		BuyerID:      0,
		TokenID:      1,
		CollectionID: 1,
	},
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
	for _, acc := range _accounts {
		if err := AddAccount(&acc); err != nil {
			return err
		}
	}

	return nil
}

func addCollections() error {
	for _, coll := range _collections {
		if err := AddCollection(&coll); err != nil {
			return err
		}
	}

	return nil
}

func addTokens() error {
	for _, t := range _tokens {
		if err := AddToken(&t); err != nil {
			return err
		}
	}

	return nil
}

func addTxs() error {
	for _, tx := range _txs {
		if err := AddTransaction(&tx); err != nil {
			return err
		}
	}

	return nil
}
