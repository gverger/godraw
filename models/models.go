package models

import (
	"encoding/json"
	"fmt"
)

type Drawing struct {
	Items []Drawable `json:"items"`
}

type Drawable interface {
	AllPoints() []Point
}

type Fillable struct {
	FillColor string `json:"fill,omitempty"`
}

type Colorable struct {
	Color string `json:"color,omitempty"`
}

type PointsDrawable struct {
	DrawPoints bool `json:"draw_points,omitempty"`
}

type Point struct {
	Colorable
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// AllPoints implements Drawable.
func (p Point) AllPoints() []Point {
	return []Point{p}
}

type Line struct {
	Colorable
	PointsDrawable
	Points []Point `json:"points"`
}

// AllPoints implements Drawable.
func (l Line) AllPoints() []Point {
	return l.Points
}

type Polygon struct {
	Colorable
	PointsDrawable
	Fillable
	Points []Point `json:"points"`
}

// AllPoints implements Drawable.
func (p Polygon) AllPoints() []Point {
	return p.Points
}

type rawDrawing struct {
	Items []json.RawMessage `json:"items"`
}

type Typed struct {
	Type string `json:"item"`
}

func (d *Drawing) UnmarshalJSON(b []byte) error {
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
			item = &Point{}
		case "line":
			item = &Line{}
		case "polygon":
			item = &Polygon{}
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

type withType struct {
	Typed
	Drawable `json:",inline"`
}

func (d Drawing) MarshalJSON() ([]byte, error) {
	var rawDrawing rawDrawing

	for _, v := range d.Items {
		var t string
		switch v.(type) {
		case *Point, Point:
			t = "point"
		case *Line, Line:
			t = "line"
		case *Polygon, Polygon:
			t = "polygon"
		default:
			return nil, fmt.Errorf("unknown drawable: %+v", v)
		}

		typed, err := json.Marshal(Typed{Type: t})
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		str := make([]byte, len(typed)+len(b)-1)
		copy(str, typed)
		copy(str[len(typed)-1:], b)
		str[len(typed)-1] = ','

		fmt.Println(string(str))

		rawDrawing.Items = append(rawDrawing.Items, str)
	}
	return json.Marshal(rawDrawing)
}
