package main

import (
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	_, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	runtime.LockOSThread()

	window := g143.NewWindow(1450, 800, "comments203: commenting on an image program.", false)
	drawMainWindow(window)

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
	theCtx.ggCtx.DrawRoundedRectangle(10, 10, toolBoxW+20, 270, 10)
	theCtx.ggCtx.Fill()

	// draw tools
	sIRect := theCtx.drawButtonA(PencilWidget, 20, 20, "Select Image", "#444", "#ddd")
	_, aCTY := nextVerticalCoords(sIRect, 10)
	aCRect := theCtx.drawButtonA(SymmLineWidget, 20, aCTY, "Add Comment", "#444", "#ddd")
	_, rCTY := nextVerticalCoords(aCRect, 10)
	rCRect := theCtx.drawButtonA(LeftSymmWidget, 20, rCTY, "Delete Comment", "#444", "#ddd")
	_, sWTY := nextVerticalCoords(rCRect, 10)
	sWRect := theCtx.drawButtonA(SaveWidget, 20, sWTY, "Save Work", "#444", "#ddd")
	_, oWDY := nextVerticalCoords(sWRect, 10)
	theCtx.drawButtonA(OpenWDWidget, 20, oWDY, "Open Folder", "#444", "#ddd")


	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
