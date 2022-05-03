package models

type Image struct {
	Origin, Destination string
}

func (i Image) GetOrigin() string {
	return i.Origin
}

func (i Image) GetDestination() string {
	return i.Destination
}
