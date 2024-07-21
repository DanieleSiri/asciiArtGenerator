package ascii

import (
	"image"
	"image/color"
	"math"
    "ascii_art/utils"
)

var FilterX = [3][3]int {
    {-1, 0, 1},
    {-2, 0, 2},
    {-1, 0, 1},
}

var FilterY = [3][3]int {
    {-1, -2, -1},
    {0, 0, 0},
    {1, 2, 1},
}

type Position struct {
    x int
    y int 
}

func Filter(img image.Image) (image.Image, map[Position]float64) {
    size := img.Bounds()
    imgGray := utils.ToGrayScale(img)  
    angleMap := make(map[Position]float64, 0)

    newImg := image.NewRGBA(image.Rect(size.Max.X,size.Max.Y,size.Min.X, size.Min.Y))
    var newColor color.Color
    for x:=0; x<size.Max.X; x++ {
        for y:=0; y<size.Max.Y; y++ {
            Gp, angle := applyFilter(imgGray, x, y)
            newColor = color.RGBA {
                R:uint8(Gp),
                G:uint8(Gp),
                B:uint8(Gp),
                A:0,
            }
            angleMap[Position{x, y}] = angle
            newImg.Set(x, y, newColor)
        }
    }
    return newImg, angleMap
}

func applyFilter(img image.Image, x, y int) (uint8, float64){
    var Gx, Gy int = 0, 0

    for i:=0; i<3; i++ {
        for j:=0; j<3; j++ {
            curX := x+j-1
            curY := y+i-1
            pixel, _, _, _ := img.At(curX, curY).RGBA()
            Gx += FilterX[i][j] * int(uint8(pixel))
            Gy += FilterY[i][j] * int(uint8(pixel))
        }
    }
    G := uint8(math.Abs(math.Sqrt(float64((Gx*Gx) + (Gy*Gy)))))
    angle := (math.Atan2(float64(Gy), float64(Gx)) / math.Pi) * 0.5 + 0.5
    return G, angle
}

