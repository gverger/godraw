package main

import (
	"encoding/json"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/godraw/models"
)

type CameraHandler struct {
	Camera *rl.Camera2D
}

func NewCameraHandler() CameraHandler {
	return CameraHandler{
		Camera: &rl.Camera2D{
			Zoom: 1.0,
		},
	}
}

func (h *CameraHandler) Update() {
	camera := h.Camera
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1.0/camera.Zoom)
		camera.Target = rl.Vector2Add(camera.Target, delta)
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		// Get the world point that is under the mouse
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), *camera)

		// Set the offset to where the mouse is
		camera.Offset = rl.GetMousePosition()

		// Set the target to match, so that the camera maps the world space point
		// under the cursor to the screen space point under the cursor at any zoom
		camera.Target = mouseWorldPos

		// Zoom increment
		scaleFactor := float32(1.0 + (0.25 * math.Abs(float64(wheel))))
		if wheel < 0 {
			scaleFactor = 1.0 / scaleFactor
		}
		camera.Zoom = rl.Clamp(camera.Zoom*scaleFactor, 0.0125, 1024.0)
	}
}

func drawPoint(p models.Point) {
	rl.DrawCircleV(rl.NewVector2(p.X, p.Y), 5, rl.Red)
}

func drawLine(l models.Line) {
	points := make([]rl.Vector2, len(l.Points))
	for i, p := range l.Points {
		points[i] = rl.NewVector2(p.X, p.Y)
	}
	rl.DrawLineStrip(points, rl.Blue)
}

func drawPolygon(poly models.Polygon) {
	points := make([]rl.Vector2, len(poly.Points))
	for i, p := range poly.Points {
		points[i] = rl.NewVector2(p.X, p.Y)
	}
	if points[len(points)-1] != points[0] {
		points = append(points, points[0])
	}
	rl.DrawLineStrip(points, rl.Green)
}

func draw(shape any) {
	switch s := shape.(type) {
	case *models.Point:
		drawPoint(*s)
	case *models.Line:
		drawLine(*s)
	case *models.Polygon:
		drawPolygon(*s)
	default:
		fmt.Printf("not a correct shape: %+v\n", s)
	}
}

func main() {
	str := `{"items":[{"item":"point","color":"red","x":12,"y":32},{"item":"line","color":"blue","draw_points":true,"points":[{"color":"red","x":2,"y":3.4},{"x":1.3,"y":5}]},{"item":"polygon","color":"blue","fill":"green","points":[{"color":"red","x":20,"y":3.4},{"x":1.3,"y":5}]}]}`

	var drawing models.Drawing
	if err := json.Unmarshal([]byte(str), &drawing); err != nil {
		fmt.Println("error:", err)
		return
	}

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(800, 450, "Shapes Visualization")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	camera := NewCameraHandler()
	camera.Camera.Offset = rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2)

	for !rl.WindowShouldClose() {
		camera.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(*camera.Camera)

		for _, s := range drawing.Items {
			draw(s)
		}

		rl.EndMode2D()
		rl.EndDrawing()
	}
}
