package main

import (
	"strings"
	"runtime"
	"time"
	"os/exec"
	"fmt"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/fogleman/gg"
	"github.com/disintegration/imaging"
)

func main() {
	_, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	runtime.LockOSThread()

	window := g143.NewWindow(1450, 900, "comments203: commenting on an image program.", false)
	drawMainWindow(window)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(24) - time.Since(t))
	}
}



func drawMainWindow(window *glfw.Window) {
	objCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight)

	// background rectangle
	theCtx.ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	theCtx.ggCtx.SetHexColor("#ddd")
	theCtx.ggCtx.Fill()

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


	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if g143.InRect(canvasRect, xPosInt, yPosInt) {
		widgetRS = canvasRect
		widgetCode = CanvasWidget
	}

	if widgetCode == 0 {
		return
	}

	rootPath, _ := GetRootPath()
	switch widgetCode {
	case SelectImageTool:
		imagePath := pickFileUbuntu("png|jpg")
		if strings.TrimSpace(imagePath) == "" {
			return
		}

		currentWorkingImagePath = imagePath

		img, err := imaging.Open(imagePath)
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

		theCtx := Continue2dCtx(currentWindowFrame)

		theCtx.ggCtx.SetHexColor("#ddd")
		theCtx.ggCtx.DrawRectangle(220, 0, float64(wWidth), float64(wHeight))
		theCtx.ggCtx.Fill()

		theCtx.ggCtx.DrawImage(img, 220, 20)

		canvasRect = g143.NewRect(220, 20, img.Bounds().Dx(), img.Bounds().Dy())

		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = theCtx.ggCtx.Image()


	case AddCommentTool, DeleteCommentTool:

		ggCtx := gg.NewContextForImage(currentWindowFrame)

		activeTool = widgetCode

		// clear indicators
		for _, cs := range drawnIndicators {
			ggCtx.SetHexColor("#dddddd")
			ggCtx.DrawCircle(float64(cs.X), float64(cs.Y), indicatorCircleR+2)
			ggCtx.Fill()
		}
		// draw an indicator on the active tool
		ggCtx.SetHexColor("#DAC166")
		ggCtx.DrawCircle(float64(widgetRS.OriginX+widgetRS.Width-20), float64(widgetRS.OriginY+20), 10)
		ggCtx.Fill()
		drawnIndicators = append(drawnIndicators, CircleSpec{X: widgetRS.OriginX + widgetRS.Width - 20, Y: widgetRS.OriginY + 20})

		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = ggCtx.Image()

	case SaveWorkTool:


	case CanvasWidget:
		fmt.Println("in canvas")


	case OpenFolderTool:
		exec.Command("xdg-open", rootPath).Run()
	default:

	}
}