package models_test

import (
	"encoding/json"
	"testing"

	"github.com/gverger/godraw/models"
	"github.com/matryer/is"
)

func TestAll(t *testing.T) {
	is := is.New(t)

	jsonString := `{
		"items": [
	{ "item": "point", "color": "red", "x": 12, "y": 32 },
	{ "item": "line", "draw_points": true, "color": "blue", "points": [{"color": "red", "x": 2, "y": 3.4}, {"x": 1.3, "y": 5}]},
	{ "item": "polygon", "draw_points": false, "fill": "green", "color": "blue", "points": [{"color": "red", "x": 2, "y": 3.4}, {"x": 1.3, "y": 5}]}
		]
	}`

	var drawing models.Drawing
	err := json.Unmarshal([]byte(jsonString), &drawing)

	is.NoErr(err)
	is.Equal(3, len(drawing.Items))
	is.Equal(&models.Point{Colorable: models.Colorable{Color: "red"}, X: 12, Y: 32}, drawing.Items[0])
	is.Equal(&models.Line{
		PointsDrawable: models.PointsDrawable{true},
		Colorable:      models.Colorable{Color: "blue"},
		Points: []models.Point{
			{Colorable: models.Colorable{Color: "red"}, X: 2, Y: 3.4},
			{X: 1.3, Y: 5},
		},
	}, drawing.Items[1])
	is.Equal(&models.Polygon{
		PointsDrawable: models.PointsDrawable{false},
		Colorable:      models.Colorable{Color: "blue"},
		Fillable:       models.Fillable{FillColor: "green"},
		Points: []models.Point{
			{Colorable: models.Colorable{Color: "red"}, X: 2, Y: 3.4},
			{X: 1.3, Y: 5},
		},
	}, drawing.Items[2])

	str, err := json.Marshal(drawing)
	is.NoErr(err)
	is.Equal(string(str), `{"items":[{"item":"point","color":"red","x":12,"y":32},{"item":"line","color":"blue","draw_points":true,"points":[{"color":"red","x":2,"y":3.4},{"x":1.3,"y":5}]},{"item":"polygon","color":"blue","fill":"green","points":[{"color":"red","x":2,"y":3.4},{"x":1.3,"y":5}]}]}`)
}
