package main

import (
    "fmt"
    "image"
    "image/color"
    // "image/draw"
    "image/png"
    "os"
    "strconv"
    "gopkg.in/cheggaaa/pb.v1"
)

var mandelColor color.Color = color.RGBA{0, 204, 102, 153}
var mag float64 = 120.0 // todo move to ARgs

func drawMandelbrot(mw, mh int, img *image.RGBA) {
	var cR, cX, zR, zX, swp float64
	var r int
	var mandelbrot []bool
	mandelbrot = make([]bool, mw*mh)
	var bar *pb.ProgressBar
	bar = pb.StartNew(mw*mh)

	for y:=0.0; y<float64(mh); y++ {
		for x:=0.0; x<float64(mw); x++ {
			cR = x/mag
			cX = y/mag
			zR = 0.0
			zX = 0.0
			for r=0; r < 50 && zR < 2; r++ {
			    swp = zR
				zR = zR*zR + -1*(zX*zX) + cR
				zX = 2*swp*zX + cX
				r += 1
			}
			if zR > 2 {
			    // x, y is in set
			    // array holds rows, so 3d row starts at mw*3+0
				mandelbrot[int(y)*mw+int(x)] = true
		 	} 
		 	bar.Increment()
		}
	}
		
	bar.FinishPrint("Booled this up")

	// second cycle, keeping it neighbour-point-aware
	var clr color.Color
	var trp uint8
	var points [8]int
	var cntf, cntt int
	bar = pb.StartNew(mw*mh)
	for y:=0; y<mh; y++ {
		for x:=0; x<mw; x++ {
			if mandelbrot[y*mw+x] {
				img.Set(x,y,mandelColor)
			} else {
				// transparency [0..255] equals 255*adjacent false/adjacent true
				// 8 points from top left to bottom right
				cntf,cntt = 0, 0
				points[0] = (y-1)*mw+x-1
				points[1] = (y-1)*mw+x
				points[2] = (y-1)*mw+x+1
				points[3] = y*mw+x-1
				points[4] = y*mw+x+1
				points[5] = (y+1)*mw+x-1
				points[6] = (y+1)*mw+x
				points[7] = (y+1)*mw+x+1
				for p:=0; p<len(points); p++ {
					if 0 < points[p] && points[p] < len(mandelbrot)-1 {
						if mandelbrot[points[p]] {
							cntt++
						} else {
							cntf++
						}
					}
				}

				trp = uint8(cntf/(cntt+cntf))
				clr = color.RGBA{trp, 0, 0, 0}
				img.Set(x,y,clr)
				bar.Increment()
			}
		}
	}
	bar.FinishPrint("Colored Up")
	

}

func main() {
	var w, h int
	if len(os.Args)>1 {
		tw, _ := strconv.ParseInt(os.Args[1], 10, 32)
		th, _ := strconv.ParseInt(os.Args[2], 10, 32)
		w = int(tw)
		h = int(th)
	} else {
		w = 250
		h = 250
	}

    file, err := os.Create("someimage.png")

    if err != nil {
        fmt.Errorf("%s", err)
        panic(err)
    }

    // this returns a bloody pointer
    img := image.NewRGBA(image.Rect(0, 0, w, h))

    drawMandelbrot(w, h, img)

    png.Encode(file, img)
}