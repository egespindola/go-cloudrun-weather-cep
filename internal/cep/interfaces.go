package cep

type CepService interface {
	GetLocationByCep(cep string) (*CepLocation, error)
}

