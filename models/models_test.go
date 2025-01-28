package models_test

import (
	"encoding/json"
	"testing"

	"github.com/gverger/godraw/models"
	"github.com/matryer/is"
)

func TestPoint(t *testing.T) {
	is := is.New(t)

	jsonString := `{
		"items": [
	{ "item": "point", "color": "red", "x": 12, "y": 32 }
		]
	}`

	var drawing models.Drawing
	err := json.Unmarshal([]byte(jsonString), &drawing)

	is.NoErr(err)
	is.Equal(1, len(drawing.Items))
	is.Equal(models.Point{X: 12, Y: 32}, drawing.Items[0])
}
