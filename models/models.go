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

type Labellable struct {
	Label  string  `json:"label,omitempty"`
	LabelX float32 `json:"label_x,omitempty"`
	LabelY float32 `json:"label_y,omitempty"`
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
	Labellable
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
	Labellable
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
	Labellable
	Points []Point `json:"points"`
}

func drawPoints(shape *PointsDrawable) {
	shape.DrawPoints = true
}

func setColor(shape *Colorable, color string) {
	shape.Color = color
}

func setFill(shape *Fillable, color string) {
	shape.FillColor = color
}

type PolyOpts struct {
}

func (p PolyOpts) DrawPoints() func(*Polygon) {
	return func(p *Polygon) { drawPoints(&p.PointsDrawable) }
}

func (p PolyOpts) Color(color string) func(*Polygon) {
	return func(p *Polygon) { setColor(&p.Colorable, color) }
}

func (p PolyOpts) Fill(color string) func(*Polygon) {
	return func(p *Polygon) { setFill(&p.Fillable, color) }
}

func NewPolygon(points []Point, options ...func(*Polygon)) Polygon {
	p := Polygon{Points: points}

	for _, o := range options {
		o(&p)
	}

	return p
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

	d.Items = make([]Drawable, 0, len(rawDrawing.Items))
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
