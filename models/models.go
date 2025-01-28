package models

import (
	"encoding/json"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Drawing struct {
	Items []Drawable `json:"items"`
}

type Drawable interface {
}

type Fillable struct {
	FillColor string `json:"fill,omitempty"`
}

type Colorable struct {
	Color string `json:"color,omitempty"`
}

type DrawPoints struct {
	DrawPoints bool `json:"draw_points,omitempty"`
}

type Point struct {
	Colorable
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Line struct {
	Colorable
	DrawPoints
	Points []Point `json:"points"`
}

type Polygon struct {
	Colorable
	DrawPoints
	Fillable
	Points []Point `json:"points"`
}

// Draw implements Drawable.
func (p Polygon) Draw() {
	points := make([]rl.Vector2, len(p.Points))
	for i, p := range p.Points {
		points[i] = rl.NewVector2(p.X, p.Y)
	}
	if points[len(points)-1] != points[0] {
		points = append(points, points[0])
	}
	rl.DrawLineStrip(points, rl.Green)
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
