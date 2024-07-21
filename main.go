package main

import (
	"ascii_art/ascii"
	"ascii_art/utils"
	"errors"
	"flag"
	"image"
    "fmt"
	"os"
	"strconv"
	"strings"
)

func getInputSize(input string) (*utils.Size, error) {
    split := strings.Split(input, "x")
    if len(split) != 2 {
        return nil, errors.New("wrong input size")
    }

    w, err := strconv.Atoi(split[0]) 
    if err != nil {
        return nil, err
    }

    h, err := strconv.Atoi(split[1])
    if err != nil {
        return nil, err
    }

    if w > 3840 {
        return nil, errors.New("maximum width allowed = 3840")
    }
    if h > 2160 {
        return nil, errors.New("maximum height allowed = 2160")
    }

    return &utils.Size{Width : uint(w), Height : uint(h)}, nil
}

func main() {
    var outputFile = flag.String("o", "img_output", "filepath to output")
    var inputFile = flag.String("f", "", "filepath to output")
    var inputSize = flag.String("size", "720x360", "size of image desired (i.e. 720x360)")
//    var debug = flag.Bool("d", false, "debugging only")
    var filterBool = flag.Bool("filter", false, "choose true or false to apply the filters")
    flag.Parse()

    if *inputFile == "" {
        panic("no input file provided")
    }
    
    size, err := getInputSize(*inputSize)

    if err != nil {
        panic(err) 
    }

    img, err := utils.LoadImage(*inputFile, *size)

    if err != nil {
        panic("could not load image")
    }

    asciiImg := ""

    if *filterBool == true {
        fmt.Println("Generating Ascii Art with filters")
        asciiImg = execFilter(img, *size)
    } else {
        asciiImg = execNormal(img, *size)
    }
    fmt.Println("Done Generating! Writing to file..")

    file_err := os.WriteFile(*outputFile, []byte(asciiImg), 0644)
    if file_err != nil {
        panic("could not write the new image")
    }
    fmt.Println("File written. Exiting")
}

func execFilter(img image.Image, size utils.Size) string {
    fmt.Println("Resizing image to 2560x1440 for better quality")
    resized4k := utils.Resize(img, utils.Size{Width : 2560, Height : 1440})
    fmt.Println("Applying Difference of Gaussians")
    newImg := ascii.DifferenceOfGaussians(resized4k, 1.0, 10.0)
    fmt.Printf("Resizing image to %dx%d\n", size.Width, size.Height)
    filtered := utils.Resize(newImg, size)
    fmt.Println("Applying Sobel Filter")
    filteredImg, angleMap := ascii.Filter(filtered)

    fmt.Println("Calculating edges")
    sobeledString := ascii.GenerateAsciiFiltered(filteredImg, make([]rune, 0), angleMap)
    asciiImg := make([]rune, 0)
    
    fmt.Printf("Resizing input image to %dx%d\n", size.Width, size.Height)
    resizedImg := utils.Resize(img, size)
    fmt.Println("Generating Ascii Art")
    asciiImg = ascii.GenerateAsciiFiltered(resizedImg, sobeledString, angleMap)
    return string(asciiImg)
}

func execNormal(img image.Image, size utils.Size) string {
    fmt.Printf("Resizing input image to %dx%d\n", size.Width, size.Height)
    resized := utils.Resize(img, size)
    fmt.Println("Generating Ascii Art")
    asciiImg := ascii.GenerateAscii(resized)
    return asciiImg
}
