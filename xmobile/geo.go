package xmobile

import (
	"go.olapie.com/x/xtype"
)

type Point xtype.Point

func NewPoint() *Point {
	return new(Point)
}

type Place xtype.Place

func NewPlace() *Place {
	return new(Place)
}

func (p *Place) SetCoordinate(c *Point) {
	p.Coordinate = (*xtype.Point)(c)
}

func (p *Place) GetCoordinate() *Point {
	return (*Point)(p.Coordinate)
}
