package ascii

import (
    "image"
    _ "image/jpeg"
    _ "image/png"
    "image/draw"
    "math"
)

type Range struct {
    min int
    max int
}

type RangeFloat struct {
    min float64
    max float64
}

var asciiChars = map[Range]rune {
    {230,256} : '@',
    {200,230} : '#',
    {180,200} : '8',
    {160,180} : '?',
    {145,160} : 'P',
    {130,145} : 'o',
    {100,130} : ':',
    {60,100} : '*',
    {30, 60} : '.',
    {0, 30} : ' ',
}

var asciiEdges = map[RangeFloat]rune {
    {0.0,0.25} : '/',
    {0.25,0.5} : '-',
    {0.5,0.75} : '\\',
    {0.75,1.0} : '|',
}

func getAsciiChar(index int) rune {
    var res rune = '\n'
    for k, v := range asciiChars {
        if index >= k.min && index < k.max {
            res = v 
        }
    }
    return res
}


func getEdge(index float64) rune {
    res := '\n'
    for k, v := range asciiEdges {
        if index >= k.min && index < k.max {
            res = v 
        }
    }
    return res
}

func GenerateAsciiFiltered(img image.Image, filtered []rune, angleMap map[Position]float64) []rune {
    bounds := img.Bounds()

    size := bounds.Size()
    result := make([]rune, (size.X+1)*size.Y)
    pos := 0
    for y:=bounds.Min.Y; y<bounds.Max.Y; y++ {
        for x:=bounds.Min.X; x<bounds.Max.X; x++ {
            var ascii_char rune
            angledChar := getEdge(angleMap[Position{x,y}])
            if len(filtered) > 0 && pos < len(filtered) && filtered[pos] != ' ' && angledChar != '\n'{
                ascii_char = angledChar
            } else {
                r, g, b, _ := img.At(x, y).RGBA()
                // weighted gray scale
                gray_value := int((float32(r) * 0.2126) + (float32(g) * 0.7152) + (float32(b) * 0.0722))
                ascii_char = getAsciiChar(int(gray_value)/256)
            }
            result[pos] = ascii_char 
            pos++
        }
        result[pos] = '\n'
        pos++
    }
    return result 
}

// testing if the Gray texture renders better than the colored
func GenerateAsciiGray(img image.Image) string {
    bounds := img.Bounds()

    
    size := bounds.Size()
    result := make([]rune, (size.X+1)*size.Y)
    pos := 0
    newGray := image.NewGray(bounds)
    draw.Draw(newGray, bounds, img, bounds.Min, draw.Src)
    for y:=bounds.Min.Y; y<bounds.Max.Y; y++ {
        for x:=bounds.Min.X; x<bounds.Max.X; x++ {
            r, _, _, _ := newGray.At(x, y).RGBA()
            luminance := math.Floor(float64(r) * 10) / 10 
            ascii_char := getAsciiChar(int(luminance)/256)
            result[pos] = ascii_char
            pos++
        }
        result[pos] = '\n'
        pos++
    }

    return string(result)
}

func GenerateAscii(img image.Image) string {
    bounds := img.Bounds()

    size := bounds.Size()
    result := make([]rune, (size.X+1)*size.Y)
    pos := 0
    for y:=bounds.Min.Y; y<bounds.Max.Y; y++ {
        for x:=bounds.Min.X; x<bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            // weighted gray scale
            gray_value := int((float32(r) * 0.2126) + (float32(g) * 0.7152) + (float32(b) * 0.0722))
            result[pos] = getAsciiChar(int(gray_value)/256)
            pos++
        }
        result[pos] = '\n'
        pos++
    }
    return string(result)
}
