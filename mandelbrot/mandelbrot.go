package main

import (
    "fmt"
    "image"
    "image/color"
    // "image/draw"
    "image/png"
    "flag"
    // "strconv"
    "gopkg.in/cheggaaa/pb.v1"
    "time"
    "os"
)

var mandelColor color.Color = color.RGBA{255, 255, 255, 255}
var mh, mw int
var mag int
var mx, my int

type MandelPoint struct {
	x int      // img x
	y int 	   // img y
	zR float64 // final complex plain value
	zX float64 // final real plain value
	r int 	   // number of iterations by which cR skyrocketed to INF; always 255 if within the set
}

func calculatePoint(x, y, imx, imy, imag int, retch *chan MandelPoint) {
	// x, y are points on a image plain
	// mx, my are by how much the mandelbrot set is moved
	// mag is mandelbrot set magnification.
	// mx, my, mag parms are there to move the mandelbrot set on the real and complex plains and magnify it
	var cR, cX, zR, zX, swp float64
	var r int
	cR = float64(x+imx)/float64(imag)
	cX = float64(y+imy)/float64(imag)
	zR = 0.0
	zX = 0.0
	for r=0; r < 255 && zR < 2; r++ { // 255 iterations go into uint8 and are enough 
		swp = zR
		zR = zR*zR + -1*(zX*zX) + cR // calculation gets rid of complex "i" making it possible to do in a program
		zX = 2*swp*zX + cX
	}
	*retch <- MandelPoint{x, y, zR, zX, r}
}

func populateMandelbrot(mandelset *[][]MandelPoint) {
	var bar *pb.ProgressBar
	afterCh := time.After(30 * time.Second)
	var retch chan MandelPoint
	retch = make(chan MandelPoint)
	bar = pb.StartNew(mw*mh)	
	for i, row := range *mandelset {
		for j := range row {
			go calculatePoint(i, j, mx, my, mag, &retch)
			bar.Increment()
		}
	}
	bar.FinishPrint("Ran go-routines")
	var tmp MandelPoint

	bar = pb.StartNew(mw*mh)	
	for i:= 0; i<len(*mandelset); i++ {
		for j:= 0; j<len((*mandelset)[0]); j++ { // presuming exact number of return values
			select {
				case tmp = <-retch:
					(*mandelset)[tmp.y][tmp.x] = tmp
				case <-afterCh:
					panic("TIME'S UP")
			}
			bar.Increment()
		}
	}
	bar.FinishPrint("Gathered results")
}

func drawImg(mandelset *[][]MandelPoint) (*image.RGBA) {
	var clr color.Color
	h := len(*mandelset)
	w := len((*mandelset)[0]) // presuming it exists because we don't do zero-length images
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for _, row := range *mandelset {
		for _, p := range row {
			if p.zR < 2 { // flippable
				clr = mandelColor
			} else { // todo: think how to color this to be pwetty
				clr = color.RGBA{uint8(p.zR), uint8(p.x), uint8(p.y), uint8(p.r)}
			}
			// x-y controls GB, red is the final Real value, transparency is based on the number of iterations it took to go over
			// since the number of iterations goes up to 255, points inside the set will all be fully opaque
			img.Set(p.x, p.y, clr)
		}
	}
	return img
}

func saveImg(img *image.RGBA, fn string) {
	f, err := os.Create(fn)
	if err != nil {
	    panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	var f string

	// this weird flag re-assigning is due to go ignoring flags for some reason :(
	flag.IntVar(&mh, "height", 250, "img height")
	flag.IntVar(&mw, "width", 250, "img width")
	flag.IntVar(&mx, "mx", 0, "set complex shift")
	flag.IntVar(&my, "my", 0, "set real shift") // positive is up all's ok
	flag.IntVar(&mag, "mag", 120, "set magnification")
	flag.StringVar(&f, "fn", "mandel.png", "output image file name")
	flag.Parse()
	mx = -mx // this is done so positive is right

	if mh == 0 || mw == 0 {
		fmt.Println("Complex plains are no reason to generate 0px images, dude")
	} else {
    	var mandelbrot [][]MandelPoint
    	mandelbrot = make([][]MandelPoint, mh)
    	for i:=0;i<mw;i++ {
    		mandelbrot[i] = make([]MandelPoint, mw)
    	}
    	populateMandelbrot(&mandelbrot)
    	img := drawImg(&mandelbrot) // no copying D:
    	saveImg(img, f)
    }
}