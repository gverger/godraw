package main

import (
	"encoding/json"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gverger/godraw/comm"
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

func (h *CameraHandler) FocusOn(drawing models.Drawing) {
	points := make([]models.Point, 0)
	for _, s := range drawing.Items {
		points = append(points, s.AllPoints()...)
	}
	// move camera for the whole scene
	if len(points) > 0 {
		minX := points[0].X
		minY := points[0].Y
		maxX := points[0].X
		maxY := points[0].Y

		for _, p := range points[1:] {
			minX = min(minX, p.X)
			minY = min(minY, p.Y)
			maxX = max(maxX, p.X)
			maxY = max(maxY, p.Y)
		}

		dx := (maxX - minX) / float32(rl.GetScreenWidth()-20)
		dy := (maxY - minY) / float32(rl.GetScreenHeight()-20)

		h.Camera.Target = rl.NewVector2((minX+maxX)/2, (minY+maxY)/2)
		h.Camera.Zoom = 1 / max(dx, dy)
	}

	h.Camera.Offset = rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2)
}

func (h CameraHandler) ScreenDistanceToWorld(dist float32) float32 {
	origin := rl.GetScreenToWorld2D(rl.NewVector2(0, 0), *h.Camera)
	x := rl.GetScreenToWorld2D(rl.NewVector2(5, 0), *h.Camera)
	size := x.X - origin.X

	return size
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

func drawPoint(p models.Point, c CameraHandler) {
	rl.DrawCircleV(rl.NewVector2(p.X, p.Y), c.ScreenDistanceToWorld(5), rl.Red)
}

func drawLine(l models.Line, c CameraHandler) {
	points := make([]rl.Vector2, len(l.Points))
	for i, p := range l.Points {
		points[i] = rl.NewVector2(p.X, p.Y)
	}
	rl.DrawLineStrip(points, rl.Blue)

	if l.DrawPoints {
		for _, p := range l.Points {
			drawPoint(p, c)
		}
	}
}

func drawPolygon(poly models.Polygon, c CameraHandler) {
	points := make([]rl.Vector2, len(poly.Points))
	for i, p := range poly.Points {
		points[i] = rl.NewVector2(p.X, p.Y)
	}
	if points[len(points)-1] != points[0] {
		points = append(points, points[0])
	}
	rl.DrawLineStrip(points, rl.Green)
	if poly.DrawPoints {
		for _, p := range poly.Points {
			drawPoint(p, c)
		}
	}
}

func draw(shape any, camera CameraHandler) {
	switch s := shape.(type) {
	case *models.Point:
		drawPoint(*s, camera)
	case *models.Line:
		drawLine(*s, camera)
	case *models.Polygon:
		drawPolygon(*s, camera)
	default:
		fmt.Printf("not a correct shape: %+v\n", s)
	}
}

func main() {
	stream := make(chan models.Drawing)

	address := "tcp://127.0.0.1:40899"
	go comm.Listen(address, stream)

	str := `{"items":[{"item":"point","color":"red","x":12,"y":32},{"item":"line","color":"blue","draw_points":true,"points":[{"color":"red","x":2,"y":3.4},{"x":1.3,"y":5}]},{"item":"polygon","color":"blue","fill":"green","points":[{"color":"red","x":20,"y":3.4},{"x":1.3,"y":5}]}]}`

	var drawing models.Drawing
	if err := json.Unmarshal([]byte(str), &drawing); err != nil {
		fmt.Println("error:", err)
		return
	}

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(1200, 800, "Shapes Visualization")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	camera := NewCameraHandler()

	camera.FocusOn(drawing)

	for !rl.WindowShouldClose() {
		select {
		case d := <-stream:
			drawing = d
			camera.FocusOn(drawing)
		default:
		}

		camera.Update()

		if rl.IsKeyPressed(rl.KeySpace) {
			camera.FocusOn(drawing)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(*camera.Camera)

		for _, s := range drawing.Items {
			draw(s, camera)
		}

		rl.EndMode2D()
		rl.EndDrawing()
	}
}
