package entities


type UserOrders struct {
	UserAddress	string					`json:"user_address"`
	Order		[]DreamshipOrderItems	`json:"order" gorm:"foreignKey:ReferenceId"`
	OrderStatus	string					`json:"order_status"`
	Payment		Transaction				`json:"payment"`
}

type ItemWebhook struct {
	Id			string	`json:"id"`
    Status 		string	`json:"status"`
    ReferenceId	string	`json:"reference_id"`
    Cost		string	`json:"cost"`
    TestOrder	bool	`json:"test_order"`
}

type Fulfillments struct{
	Id			string	`json:"id"`
	LineItems	[]LineItem	`json:"line_items"`
	Trackings	[]Tracking	`json:"trackings"`
}

type Tracking struct {
	TrackingNumber 	string	`json:"tracking_number"`
	Carrier			string	`json:"carrier"`
	CarrierUrl		string	`json:"carrier_url"`
	Status			string	`json:"status"`
}

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

type DreamshipOrderItems struct {
	ReferenceId		string		`json:"reference_id"`
	TestOrder		bool		`json:"test_order"`
	ShippingMethod	string		`json:"shipping_method"`
	LineItems		[]LineItem	`json:"line_items"`
	Address			Address	`json:"address"`
}

type Address struct {
	FirstName				string	`json:"first_name"`
	LastName				string	`json:"last_name"`
	Company					string	`json:"company"`
	Phone					string	`json:"phone"`
	Street1					string	`json:"street1"`
	Street2					string	`json:"street2"`			
	City					string	`json:"city"`
	State					string	`json:"state"`
	Country					string	`json:"country"`
	Zip						string	`json:"zip"`
	ForceVerifiedDelivery	bool	`json:"force_verified_delivery"`
}

type LineItem struct {
	PrintAreas		[]PrintArea	`json:"print_areas"`
	ReferenceId		string		`json:"reference_id"`
	Quantity		int64		`json:"quantity"`
	ItemVariant		int64		`json:"item_variant"`
}

type PrintArea struct {
	Key			string	`json:"key"`
	Url			string	`json:"url"`
	Position	string	`json:"position"`
	Resize		string	`json:"resize"`
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