// Package service provides image resizing functionality
package services

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/babilu-online/common/context"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"

	// Register decoders for additional image formats
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/webp"
)

const (
	// ServiceID is the unique identifier for the resize service
	ServiceID = "resize_svc"
	// DefaultJPEGQuality is the quality setting for JPEG encoding
	DefaultJPEGQuality = 100
)

// ResizeService handles image resizing operations
type ResizeService struct {
	context.DefaultService
}

// ID returns the service identifier
func (svc ResizeService) ID() string {
	return ServiceID
}

// Start initializes the resize service
func (svc *ResizeService) Start() error {
	return nil
}

// Resize scales an image to the specified size while maintaining aspect ratio
// size parameter represents the target height in pixels
func (svc *ResizeService) Resize(data []byte, out io.Writer, size int) error {
	if len(data) == 0 {
		return fmt.Errorf("empty image data")
	}
	if size <= 0 {
		return fmt.Errorf("invalid size: %d", size)
	}

	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if format == "gif" {
		return svc.handleGIF(data, out, size)
	}

	resized := resize.Resize(0, uint(size), src, resize.MitchellNetravali)
	return svc.encodeImage(resized, format, out)
}

// encodeImage writes the resized image to the output writer in the specified format
func (svc *ResizeService) encodeImage(img image.Image, format string, out io.Writer) error {
	switch format {
	case "png":
		return png.Encode(out, img)
	case "jpeg", "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: DefaultJPEGQuality})
	default:
		return jpeg.Encode(out, img, &jpeg.Options{Quality: DefaultJPEGQuality})
	}
}

// handleGIF processes and resizes animated GIF images
func (svc *ResizeService) handleGIF(data []byte, out io.Writer, size int) error {
	gifImg, err := svc.resizeGIF(data, 0, size/2)
	if err != nil {
		return fmt.Errorf("failed to resize GIF: %w", err)
	}
	return gif.EncodeAll(out, gifImg)
}

// resizeGIF resizes all frames in a GIF image
func (svc *ResizeService) resizeGIF(data []byte, width, height int) (*gif.GIF, error) {
	img, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode GIF: %w", err)
	}

	if width == 0 {
		width = int(float64(img.Config.Width) * float64(height) / float64(img.Config.Height))
	} else if height == 0 {
		height = int(float64(img.Config.Height) * float64(width) / float64(img.Config.Width))
	}

	img.Config.Width = width
	img.Config.Height = height

	buffer := image.NewRGBA(img.Image[0].Bounds())
	for i, frame := range img.Image {
		bounds := frame.Bounds()
		draw.Draw(buffer, bounds, frame, bounds.Min, draw.Over)
		img.Image[i] = svc.convertToPaletted(resize.Resize(uint(width), uint(height), buffer, resize.MitchellNetravali))
	}

	return img, nil
}

// convertToPaletted converts any image to a paletted image using Floyd-Steinberg dithering
func (svc *ResizeService) convertToPaletted(img image.Image) *image.Paletted {
	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.Point{})
	return paletted
}
