package models

import (
	"encoding/json"
	"fmt"
)

type Drawing struct {
	Items []Drawable `json:"items"`
}

type Drawable interface {
	Draw()
}

type Colorable struct {
	Color string `json:"color"`
}

type Point struct {
	Color string  `json:"color"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
}

// Draw implements Drawable.
func (p Point) Draw() {
	panic("unimplemented")
}

type Line struct {
	Colorable
	Points []Point `json:"points"`
}

// Draw implements Drawable.
func (l Line) Draw() {
	panic("unimplemented")
}

type Polygon struct {
	Colorable
	Points []Point `json:"points"`
}

// Draw implements Drawable.
func (p Polygon) Draw() {
	panic("unimplemented")
}

type rawDrawing struct {
	Items []json.RawMessage `json:"items"`
}

type Typed struct {
	Type string `json:"item"`
}

func (d *Drawing) UnmarshallJSON(b []byte) error {
	var rawDrawing rawDrawing

	if err := json.Unmarshal(b, &rawDrawing); err != nil {
		return err
	}

	for _, raw := range rawDrawing.Items {
		var typed Typed
		if err := json.Unmarshal(raw, &typed); err != nil {
			return err
		}

		var item Drawable

		switch typed.Type {
		case "point":
			item = Point{}
		case "line":
			item = Line{}
		case "polygon":
			item = Polygon{}
		default:
			return fmt.Errorf("unknown item type: %s", typed.Type)
		}

		if err := json.Unmarshal(raw, item); err != nil {
			return err
		}
		d.Items = append(d.Items, item)
	}

	return nil
}
