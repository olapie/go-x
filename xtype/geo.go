package xtype

type Point struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

type Place struct {
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	Coordinate *Point `json:"coordinate,omitempty"`
}
