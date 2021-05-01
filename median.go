package processing

import (
	"image"
	"sort"
)

func GetMedianImage(m *image.Gray, rad int) *image.Gray{

	return GetImgMat(MedianFilter(GetImgArray(m), rad))
}
func GetImgArray(img *image.Gray) [][][]uint8 {
	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X
	channels := 1
	i := 0

	flatArr := img.Pix

	arr := make([][][]uint8, height)

	for row := 0; row < height; row++ {
		arr[row] = make([][]uint8, width)

		for col := 0; col < width; col++ {
			arr[row][col] = make([]uint8, channels)

			for ch := 0; ch < channels; ch++ {
				arr[row][col][ch] = flatArr[i]
				i++
			}
		}
	}

	return arr
}
func getAbs(num int) int {
	if num < 0 {
		num = -num
	}

	return num
}

func getMedian(list []uint8) uint8 {
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })

	return list[(len(list)-1)/2]
}

func getMirror(p, bound int) int {
	if p < 0 {
		p = getAbs(p + 1)
	} else if p >= bound {
		p = bound*2 - p - 1
	}

	return p
}

// MedianFilter - 3x3 median filter
func MedianFilter(imgArr [][][]uint8, rad int) [][][]uint8 {
	height := len(imgArr)
	width := len(imgArr[0])

	// Init new array
	newArr := make([][][]uint8, height)

	for r := 0; r < height; r++ {
		newArr[r] = make([][]uint8, width)

		for c := 0; c < width; c++ {
			newArr[r][c] = make([]uint8, 1)

			list := make([]uint8, 0)

			// Matrix loop
			for j := -rad; j <= rad; j++ {
				for i := -rad; i <= rad; i++ {
					list = append(list, imgArr[getMirror(r+j, height)][getMirror(c+i, width)][0])
				}
			}

			newArr[r][c][0] = getMedian(list)
		}
	}

	return newArr
}

// GetImgMat - Convert cvMat to 3-dimension array
func GetImgMat(arr [][][]uint8) *image.Gray {
	height := len(arr)
	width := len(arr[0])
	channels := len(arr[0][0])
	i := 0
	gray := image.NewGray(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: width, Y: height}})

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			for ch := 0; ch < channels; ch++ {
				gray.Pix[i] = arr[row][col][ch]
				i++
			}
		}
	}

	return gray

}
