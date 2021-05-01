package deskew

import (
	"image"
	"math"
	"sort"
)

const (
	StepsPerDegree  = 10
	MaxSkewToDetect = 30
	LocalPeakRadius = 4
	IntensityMaxCount = 5
	IntensityThreshold = 0.5
)

type SkewAngle struct {
	stepsPerDegree  int
	houghHeight     int
	thetaStep       float64
	maxSkewToDetect float64
	sinMap          []float64
	cosMap          []float64
}

type HoughLine struct {
	theta             float64
	radius            int16
	intensity         int16
	relativeIntensity float64
}

func New() *SkewAngle {
	sa := &SkewAngle{}
	sa.houghHeight = 2 * MaxSkewToDetect * StepsPerDegree
	sa.thetaStep = (float64(2) * float64(MaxSkewToDetect) * math.Pi / 180) / float64(sa.houghHeight)
	minTheta := 90.0
	for i := 0; i < sa.houghHeight; i++ {
		val := (minTheta * math.Pi / 180) + float64(float64(i)*sa.thetaStep)
		sa.sinMap = append(sa.sinMap, math.Sin(val))
		sa.cosMap = append(sa.cosMap, math.Cos(val))
	}
	return sa
}
func (sa *SkewAngle) GetSkewAngle(image *image.Gray) float64 {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y
	halfWidth := width / 2
	halfHeight := height / 2

	startX := - halfWidth
	startY := - halfHeight

	stopX := width - halfWidth
	stopY := height - halfHeight - 1

	offset := image.Stride - width

	halfHoughWidth := int(math.Sqrt(float64(halfWidth*halfWidth + halfHeight*halfHeight)))
	houghWidth := halfHoughWidth * 2

	houghMap := sa.initHoughMap()
	src := 0
	srcBelow := image.Stride
	var maxIntensity int16
	var maxRadius int
	var maxTheta int
	for y := startY; y < stopY; y++ {
		for x := startX; x < stopX && srcBelow < len(image.Pix); x, src, srcBelow = x+1, src+1, srcBelow+1 {
			if (image.Pix[src] < 128) && (image.Pix[srcBelow] >= 128) {
				for theta := 0; theta < sa.houghHeight; theta++ {
					radius := int(sa.cosMap[theta]*float64(x)-sa.sinMap[theta]*float64(y)) + halfWidth
					if radius < 0 || radius >= houghWidth {
						continue
					}
					houghMap[theta][radius]++
					if maxIntensity < houghMap[theta][radius] {
						maxIntensity = houghMap[theta][radius]
					}
					if maxRadius < radius {
						maxRadius = radius
					}
					if maxTheta < theta {
						maxTheta = theta
					}
				}
			}
			src = src + offset
			srcBelow = srcBelow + offset
		}
	}
	houghLines := sa.collectLine(houghMap, maxIntensity, int16(width/10), maxTheta, maxRadius)
	if houghLines == nil {
		return 0
	}
	var skewAngle float64
	var sumIntensity float64
	for _, hl := range houghLines {
		if hl.relativeIntensity > IntensityThreshold {
			skewAngle = skewAngle + hl.theta * hl.relativeIntensity
			sumIntensity = sumIntensity + hl.relativeIntensity
		}
	}
	if sumIntensity > 0 {
		skewAngle = skewAngle / sumIntensity
	}
	return skewAngle - 60.0

}
func (sa *SkewAngle) collectLine(houghMap []map[int]int16, maxIntensity, minLineIntensity int16, maxTheta, maxRadius int) []*HoughLine {
	var intensity int16

	var houghLines []*HoughLine

	halfHoughWidth := maxRadius >> 1
	for theta := 0; theta < maxTheta; theta++ {
		for radius := 0; radius < maxRadius; radius++ {
			intensity = houghMap[theta][radius]
			if intensity < minLineIntensity {
				//continue
			}
			for tt, ttMax := theta-LocalPeakRadius, theta+LocalPeakRadius; tt < ttMax; tt++ {
				if tt < 0 {
					continue
				}
				if tt > maxTheta {
					break
				}
				for tr, trMax := radius-LocalPeakRadius, radius+LocalPeakRadius; tr < trMax; tr++ {
					if tr < 0 {
						continue
					}
					if tr >= maxRadius {
						break
					}
					if houghMap[tt][tr] > intensity {
						goto END
					}
				}
			}
			houghLines = append(houghLines,
				&HoughLine{
					theta:             90.0 - MaxSkewToDetect + float64(theta)/float64(StepsPerDegree),
					radius:            int16(radius - halfHoughWidth),
					intensity:         intensity,
					relativeIntensity: float64(intensity) / float64(maxIntensity)})
			END:
		}

	}
	sort.SliceStable(houghLines, func(i, j int) bool {
		return houghLines[i].intensity > houghLines[j].intensity
	})
	if (len(houghLines) == 0){
		return nil
	}
	return houghLines[:int(math.Max(float64(IntensityMaxCount), float64(len(houghLines)-1)))]
}
func (sa *SkewAngle) initHoughMap() []map[int]int16 {
	hough := make([]map[int]int16, sa.houghHeight)
	for i := 0; i < sa.houghHeight; i++ {
		hough[i] = make(map[int]int16)
	}
	return hough
}
