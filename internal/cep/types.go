package cep

type CepLocation struct {
	CEP          string   `json:"cep"`
	State        string   `json:"state"`
	City         string   `json:"city"`
	Neighborhood string   `json:"neighborhood"`
	Street       string   `json:"street"`
	Service      string   `json:"service"`
	IBGE         IBGE     `json:"ibge"`
	Timezone     string   `json:"timezone"`
	Location     Location `json:"location"`
}

type IBGE struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type Coordinates struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type Location struct {
	Coordinates Coordinates `json:"coordinates"`
}
