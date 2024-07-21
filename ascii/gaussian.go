package ascii

import (
	"image"
	"image/color"
	"math"
    "ascii_art/utils"
	"gonum.org/v1/gonum/floats"
)

func applyGaussian(img *image.Gray, sigma float64) *image.Gray {
    size := img.Bounds()
    kernel := createGaussianKernel(sigma)
    blurred := image.NewGray(size)
    width := size.Dx()
    height := size.Dy()

    for y:=0; y<height; y++ {
        for x:=0; x<width; x++ {
            sum, weightSum := 0.0, 0.0
            
            for ky := -len(kernel)/2; ky<=len(kernel)/2; ky++ {
                for kx := -len(kernel)/2; kx<=len(kernel)/2; kx++ {
                    ix := x + kx
                    iy := y + ky
                    if ix >= 0 && ix < width && iy >= 0 && iy < height {
						weight := kernel[ky+len(kernel)/2] * kernel[kx+len(kernel)/2]
						sum += float64(img.GrayAt(ix, iy).Y) * weight
						weightSum += weight
					}
				}
			}
            blurred.SetGray(x, y, color.Gray{uint8(sum/weightSum)})
        }
    }
    return blurred
}

func createGaussianKernel(sigma float64) []float64 {
	size := int(math.Ceil(6*sigma + 1))
	if size%2 == 0 {
		size++
	}
	kernel := make([]float64, size)
	mid := size / 2

	for i := -mid; i <= mid; i++ {
		kernel[i+mid] = math.Exp(-float64(i*i) / (2 * sigma * sigma))
	}
	floats.Scale(1/floats.Sum(kernel), kernel)
	return kernel
}

func DifferenceOfGaussians(img image.Image, sigma1, sigma2 float64) image.Image {
    newGray := utils.ToGrayScale(img)
    blurred1 := applyGaussian(newGray, sigma1)
    blurred2 := applyGaussian(newGray, sigma2)

    // subtract the images
    dogImg := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c1 := blurred1.GrayAt(x, y).Y
			c2 := blurred2.GrayAt(x, y).Y
			dogColor := int16(c1) - int16(c2)
			if dogColor < 0 {
				dogColor = 0
			}
			if dogColor > 255 {
				dogColor = 255
			}
			dogImg.SetGray(x, y, color.Gray{uint8(dogColor)})
		}
	}

    return dogImg
}
