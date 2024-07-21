package utils

import (
    "image"
    "os"
    "github.com/nfnt/resize"
    "image/color"
)

type Size struct {
    Width uint
    Height uint
}

func Resize(img image.Image, rect Size) image.Image {
    resized := resize.Resize(rect.Width, rect.Height, img, resize.Bilinear)
    return resized
}

// function that loads the image from a path and returns it resized
func LoadImage(img_path string, rect Size) (image.Image, error){
    reader, file_err := os.Open(img_path)
    if file_err != nil {
        return nil, file_err
    }
    defer reader.Close()
    
    img, _, err := image.Decode(reader)
    
    if err != nil {
        return nil, err
    }

    return img, nil
}

// convert image to grayscale
func ToGrayScale(img image.Image) *image.Gray {
    newGray := image.NewGray(img.Bounds())
    for x:=0; x<img.Bounds().Max.X; x++ {
        for y:=0; y<img.Bounds().Max.Y; y++ {
            gray := color.GrayModel.Convert(img.At(x, y))
            newGray.Set(x, y, gray)
        }
    }
    return newGray
}
