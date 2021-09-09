package restaurant

type mapResponse struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

type Result struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
	PlusCode          *PlusCode          `json:"plus_code,omitempty"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Bounds       *Bounds  `json:"bounds,omitempty"`
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Bounds   `json:"viewport"`
}

type Bounds struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}
