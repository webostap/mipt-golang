package timepng

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"time"
)

// TimePNG записывает в `out` картинку в формате png с текущим временем
func TimePNG(out io.Writer, t time.Time, c color.Color, scale int) {
	png.Encode(out, buildTimeImage(t, c, scale))
}

// buildTimeImage создает новое изображение с временем `t`
func buildTimeImage(t time.Time, c color.Color, scale int) *image.RGBA {
	timeStr := t.Format("15:04")

	charWidth := scale * 3
	charHeight := scale * 5

	img := image.NewRGBA(image.Rect(0, 0, charWidth*len(timeStr)+scale*(len(timeStr)-1), charHeight))

	for i, r := range timeStr {
		sub := img.SubImage(image.Rect(i*charWidth+i*scale, 0, (i+1)*charWidth+i*scale, charHeight)).(*image.RGBA)
		sub.Rect = image.Rect(0, 0, charWidth, charHeight)
		fillWithMask(sub, nums[r], c, scale)
	}

	return img
}

// fillWithMask заполняет изображение `img` цветом `c` по маске `mask`. Маска `mask`
// должна иметь пропорциональные размеры `img` с учетом фактора `scale`
// NOTE: Так как это вспомогательная функция, можно считать, что mask имеет размер (3x5)
func fillWithMask(img *image.RGBA, mask []int, c color.Color, scale int) {
	for i, val := range mask {
		x, y := i%3, i/3
		if val == 1 {
			for xi := 0; xi < scale; xi++ {
				for yi := 0; yi < scale; yi++ {
					img.Set(scale*x+xi, scale*y+yi, c)
				}
			}
		}
	}
}

var nums = map[rune][]int{
	'0': {
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'1': {
		0, 1, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'2': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	'3': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'4': {
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'5': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'6': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'7': {
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	'8': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'9': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	':': {
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
	},
}
