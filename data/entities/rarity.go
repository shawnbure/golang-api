package entities

type TokenRarity struct {
	UsedTraitsCount   int     `json:"usedTraitsCount"`
	StatRarity        float64 `json:"statRarity"`
	RarityScore       float64 `json:"rarityScore"`
	AvgRarity         float64 `json:"avgRarity"`
	RarityScoreNormed float64 `json:"rarityScoreNormed"`
}
