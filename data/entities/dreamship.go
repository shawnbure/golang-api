package entities

type DreamshipItems struct {
	Id					int64			`json:"id"`
	Name				string			`json:"name"`
	Description 		string			`json:"description"`
	ReleaseStatus		string			`json:"release_status"`
	OperationalStatus	string			`json:"operational_status"`
	ProductionDaysMin 	int64			`json:"production_days_min"`
	ProductionDaysMax	int64			`json:"production_days_max"`
	Colors				[]VariantColor	`json:"colors"`
	Sizes				[]VariantSize	`json:"sizes"`
	PrintAreas			[]interface{}	`json:"print_areas"`
	ShipZones			[]interface{}	`json:"ship_zones"`
	ItemVariants		[]ItemVariants	`json:"item_variants"`
}

type ItemVariants struct {
	Id				int64			`json:"id"`
	Name			string			`json:"name"`
	Color			VariantColor	`json:"color"`
	Size			VariantSize		`json:"size"`
	Availability	string			`json:"availability"`
	Cost			float64			`json:"cost"`
}

type VariantColor struct {
	Id				int64	`json:"id"`
	Name			string	`json:"name"`
	PrimaryHex		string	`json:"primary_hex"`
	SecondaryHex	string	`json:"secondary_hex"`
	Pattern			string	`json:"pattern"`
	Availability	string	`json:"availability"`
}

type VariantSize struct {
	Id				int64	`json:"id"`
	Name			string	`json:"name"`
	Availability	string	`json:"availability"`
}

type ShippingMethodResponse struct {
	Code	string				`json:"code"`
	Name	string				`json:"name"`
	Methods	[]ShippingMethod	`json:"methods"`
}

type ShippingMethod struct {
	Cost			float64	`json:"cost"`
	Method			string	`json:"method"`
	DeliveryDaysMax	uint64	`json:"delivery_days_max"`
	DeliveryDaysMin	uint64	`json:"delivery_days_min"`
}