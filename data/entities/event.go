package entities

type BlockEvents struct {
	Hash   string  `json:"hash"`
	Events []Event `json:"events"`
}

type FinalizedBlock struct {
	Hash string `json:"hash"`
}

type Event struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     [][]byte `json:"topics"`
	Data       []byte   `json:"data"`
}
