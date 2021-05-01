package processing

import (
	"image"
	"sort"
)

const cutoffTh = 2
const countTh = 2
const mergeTh = 4

func getRowFreq(arr [][][]uint8) []int {
	freq := make([]int, len(arr))

	for r := range arr {
		for c := range arr[r] {
			if arr[r][c][0] == 0 {
				freq[r]++
			}
		}
	}

	return freq
}
func GetSplitLines(m *image.Gray)([]int, []int) {

	starts, ends := splitLine(GetImgArray(m))
	return starts, ends
}
func splitLine(arr [][][]uint8) ([]int, []int) {
	freq := getRowFreq(arr)

	startMark := make([]int, 0)
	endMark := make([]int, 0)
	wCount := 0
	bCount := 0
	start := 0
	end := 0
	sumText := 0
	sumSpace := 0
	lineCount := 0

	for r := range arr {
		if freq[r] >= countTh {
			// Spot black
			if wCount > 0 {
				if wCount >= cutoffTh {
					start = r
				}

				wCount = 0
			}

			bCount++

		} else {
			// Spot white
			if bCount > 0 {
				if bCount >= cutoffTh {
					end = r - 1

					if len(endMark) > 0 {
						sumSpace += start - endMark[len(endMark)-1]
					}

					startMark = append(startMark, start)
					endMark = append(endMark, end)

					lineCount++
					sumText += end - start

				}

				bCount = 0
			}

			wCount++
		}
	}

	if lineCount > 0 {
		// Need to merge check
		avgSpace := sumSpace / (lineCount - 1)
		avgText := sumText / lineCount

		if avgText/avgSpace > mergeTh {
			// Merge to single line
			startMark = startMark[:1]
			endMark = endMark[len(endMark)-1:]
		} else {
			// Merge to multiple line
			for i := len(startMark) - 1; i > 0; i-- {
				if startMark[i]-endMark[i-1] < avgSpace {
					// Merge line
					endMark[i-1] = endMark[i]

					startMark = startMark[:i+copy(startMark[i:], startMark[i+1:])]
					endMark = endMark[:i+copy(endMark[i:], endMark[i+1:])]
				}
			}
		}
	}

	return startMark, endMark
}

func updateRect(old, new image.Rectangle) image.Rectangle {
	if new.Max.X > old.Max.X {
		old.Max.X = new.Max.X
	}

	if new.Min.X < old.Min.X {
		old.Min.X = new.Min.X
	}

	if new.Max.Y > old.Max.Y {
		old.Max.Y = new.Max.Y
	}

	if new.Min.Y < old.Min.Y {
		old.Min.Y = new.Min.Y
	}

	return old
}

func getMapRoot(maptable map[int]int, val int) int {
	for val != maptable[val] {
		val = maptable[val]
	}

	return val
}
func GetSplitChars(m *image.Gray)[]image.Rectangle {

	rects := GetSegmentChar(GetImgArray(m))
	return rects
}
func GetSegmentChar(imgArr [][][]uint8) []image.Rectangle {
	grass := make([][]int, len(imgArr))
	maptable := make(map[int]int)
	recttable := make(map[int]image.Rectangle)
	num := 0

	// Init grass
	for r := range imgArr {
		grass[r] = make([]int, len(imgArr[r]))

		for c := range imgArr[r] {
			if imgArr[r][c][0] != 255 {
				grass[r][c] = 0
			} else {
				grass[r][c] = -1
			}
		}
	}

	// Start a fire
	for y := range grass {
		for x := range grass[y] {
			if grass[y][x] >= 0 {
				found := make([]int, 0)

				// Contour search
				searchArea := [][]int{{y, x - 1}, {y - 1, x - 1}, {y - 1, x}, {y - 1, x + 1}}

				for s := range searchArea {
					j := searchArea[s][0]
					i := searchArea[s][1]

					if j >= 0 && j < len(grass) && i >= 0 && i < len(grass[y]) && grass[j][i] > 0 {
						found = append(found, grass[j][i])
					}
				}

				if len(found) == 0 {
					// New object
					num++
					grass[y][x] = num
					maptable[num] = num
					recttable[num] = image.Rectangle{image.Point{x, y}, image.Point{x + 1, y + 1}}
				} else {
					// Same object
					sort.Ints(found)

					rootNode := getMapRoot(maptable, found[0])
					grass[y][x] = rootNode
					recttable[rootNode] = updateRect(recttable[rootNode], image.Rectangle{image.Point{x, y}, image.Point{x + 1, y + 1}})

					// Update maptable and recttable
					for k := 1; k < len(found); k++ {
						if newRect, ok := recttable[found[k]]; ok && found[k] != rootNode {
							maptable[found[k]] = rootNode
							recttable[rootNode] = updateRect(recttable[rootNode], newRect)
							delete(recttable, found[k])
						}
					}
				}
			}
		}
	}

	// Map to array
	rectArray := make([]image.Rectangle, 0)

	for _, r := range recttable {
		rectArray = append(rectArray, r)
	}

	sort.Slice(rectArray, func(i, j int) bool { return rectArray[i].Min.X < rectArray[j].Min.X })

	return rectArray
}
