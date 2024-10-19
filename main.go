package main

import (
	"runtime"
	"time"
	"fmt"
	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/disintegration/imaging"
	"image/draw"
	"image/color"
	"image"
	"os"
	"strings"
	"path/filepath"
	"encoding/json"
	// "log"
)

func main() {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 {
		if ! strings.HasSuffix(os.Args[1], ".c203f") {
			fmt.Println("Invalid file extension.")
			os.Exit(1)
		}

		workingPath := filepath.Join(rootPath, ".tmp", UntestedRandomString(3))
		toClearTmp := workingPath

		files := unpackTar(os.Args[1], workingPath)
		for _, f := range files {
			if strings.HasSuffix(f, ".json") {
				rawJSON, _ := os.ReadFile(filepath.Join(workingPath, f))
				var objs []Comment
				json.Unmarshal(rawJSON, &objs)

				comments = objs
			} else {
				currentWorkingImagePath = filepath.Join(workingPath, f)
			}
		}
	}

	runtime.LockOSThread()

	window := g143.NewWindow(1450, 900, "comments203: commenting on an image program.", false)
	drawMainWindow(window)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// quick hover effect
	window.SetCursorPosCallback(getHoverCB(objCoords))
	// clear tmp on close
	window.SetCloseCallback(CloseCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(24) - time.Since(t))
	}
}



func drawMainWindow(window *glfw.Window) {
	objCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &objCoords)

	// draw tools box
	theCtx.ggCtx.SetHexColor("#DAC166")
	theCtx.ggCtx.DrawRoundedRectangle(10, 20, toolBoxW+20, 270, 10)
	theCtx.ggCtx.Fill()

	// draw tools
	sIRect := theCtx.drawButtonA(SelectImageTool, 20, 30, "Select Image", "#444", "#ddd")
	_, aCTY := nextVerticalCoords(sIRect, 10)
	aCRect := theCtx.drawButtonA(AddCommentTool, 20, aCTY, "+ Comment", "#444", "#ddd")
	_, rCTY := nextVerticalCoords(aCRect, 10)
	rCRect := theCtx.drawButtonA(DeleteCommentTool, 20, rCTY, "- Comment", "#444", "#ddd")
	_, sWTY := nextVerticalCoords(rCRect, 10)
	sWRect := theCtx.drawButtonA(SaveWorkTool, 20, sWTY, "Save Work", "#444", "#ddd")
	_, oWDY := nextVerticalCoords(sWRect, 10)
	theCtx.drawButtonA(OpenFolderTool, 20, oWDY, "Open Folder", "#444", "#ddd")

	if currentWorkingImagePath != "" {
		img, err := imaging.Open(currentWorkingImagePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		availWidth, availHeight := wWidth - 220, wHeight - 100
		if img.Bounds().Dx() > availWidth {
			newWidth := availWidth - 100
			newHeightf64 := (float64(newWidth)/float64(img.Bounds().Dx())) * float64(img.Bounds().Dy())
			img = imaging.Fit(img, newWidth, int(newHeightf64), imaging.Lanczos)
		}

		if img.Bounds().Dy() > availHeight {
			newHeight := availHeight - 100
			newWidthF64 := (float64(newHeight)/float64(img.Bounds().Dy())) * float64(img.Bounds().Dx())
			img = imaging.Fit(img, int(newWidthF64), newHeight, imaging.Lanczos)
		}

		theCtx.ggCtx.SetHexColor("#fff")
		theCtx.ggCtx.DrawRectangle(220, 0, float64(wWidth), float64(wHeight))
		theCtx.ggCtx.Fill()

		finalDst := image.NewRGBA(theCtx.ggCtx.Image().Bounds())
		draw.Draw(finalDst, finalDst.Bounds(), theCtx.ggCtx.Image(), image.Point{}, draw.Src)
		alphaImg := imaging.New(img.Bounds().Dx(), img.Bounds().Dy(), color.Alpha{A: 150})

		draw.DrawMask(finalDst, image.Rect(220, 20, img.Bounds().Dx()+220, img.Bounds().Dy()+20), 
			img, image.Point{}, alphaImg, image.Point{}, draw.Over)

		theCtx.ggCtx.DrawImage(finalDst, 0, 0)

		canvasRect = g143.NewRect(220, 20, img.Bounds().Dx(), img.Bounds().Dy())

		for i, commentObj := range comments {
			inputId := 1000 + i + 1
			theCtx.drawCommentBox(inputId, commentObj.X, commentObj.Y)
		}		
	}

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}


// clear temporary files
func CloseCallback(w *glfw.Window) {
	if toClearTmp != "" {
		os.RemoveAll(toClearTmp)
	}
}