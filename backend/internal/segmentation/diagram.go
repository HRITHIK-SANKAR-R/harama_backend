package segmentation

import (
	"image"
	"image/color"
	"image/draw"
	"bytes"
	"image/jpeg"
	"image/png"
)

// DiagramDetector handles the identification of non-text regions in exam papers
// Refactored to stub implementation to avoid OpenCV system dependency in this environment
type DiagramDetector struct{}

func NewDiagramDetector() *DiagramDetector {
	return &DiagramDetector{}
}

// DetectRegions identifies potential diagram regions in an image
func (d *DiagramDetector) DetectRegions(imgBytes []byte) ([]image.Rectangle, error) {
	// MOCK implementation: Returns a sample rectangle to simulate detection
	return []image.Rectangle{
		image.Rect(100, 100, 500, 400),
	}, nil
}

// ExtractRegion crops a specific rectangle from the image
func (d *DiagramDetector) ExtractRegion(imgBytes []byte, rect image.Rectangle) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	// Create a new image for the cropped region
	cropped := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(cropped, cropped.Bounds(), img, rect.Min, draw.Src)

	var buf bytes.Buffer
	err = png.Encode(&buf, cropped)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper to draw detected regions (useful for debugging/UI)
func (d *DiagramDetector) DrawRegions(imgBytes []byte, regions []image.Rectangle) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	// Convert to RGBA for drawing
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Draw red rectangles
	red := color.RGBA{255, 0, 0, 255}
	for _, rect := range regions {
		drawRect(rgba, rect, red)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawRect(img *image.RGBA, rect image.Rectangle, col color.Color) {
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		img.Set(x, rect.Min.Y, col)
		img.Set(x, rect.Max.Y, col)
	}
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		img.Set(rect.Min.X, y, col)
		img.Set(rect.Max.X, y, col)
	}
}