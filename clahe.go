package processing

import (
	"image"
	"image/color"
)

type Histogram struct {
	Histogram []int
	numPixel  int
}

// Histogram creates a histogram of a grayscale image.
func GetHistogram(m *image.Gray) *Histogram {

	hist := new(Histogram)

	hist.Histogram = make([]int, 256)
	hist.numPixel = len(m.Pix)
	for i := 0; i < hist.numPixel; i++ {
		hist.Histogram[m.Pix[i]]++
	}
	return hist
}

func (h *Histogram) Equalize(gray *image.Gray) *image.Gray {
	equialized := h.equalize()
	for x := gray.Bounds().Min.X; x < gray.Bounds().Max.X; x++ {
		for y := gray.Bounds().Min.Y; y < gray.Bounds().Max.Y; y++ {
			gray.SetGray(x,y,
				color.Gray{Y:equialized[gray.GrayAt(x,y).Y]})
		}
	}
	return gray
}

func (h *Histogram) equalize() []uint8 {
	if h.numPixel == 0 || len(h.Histogram) == 0{
		return nil
	}
	equalizedHistogram := make([]uint8, 256)
	coef := float64(255) / float64(h.numPixel)
	prev := float64(0)
	for i := 0; i < 256; i++ {
		prev = prev + (float64(h.Histogram[i]) * coef)
		equalizedHistogram[i] = uint8(prev)
	}
	return equalizedHistogram
}
