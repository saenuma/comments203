package main

import (
	"image"
	"runtime"
	"fmt"
	"strings"
	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func getHoverCB(state map[int]g143.Rect) glfw.CursorPosCallback {
	return func(window *glfw.Window, xpos, ypos float64) {
		if runtime.GOOS == "linux" {
			// linux fires too many events
			cursorEventsCount += 1
			if cursorEventsCount != 10 {
				return
			} else {
				cursorEventsCount = 0
			}
		}

		wWidth, wHeight := window.GetSize()

		var widgetRS g143.Rect
		var widgetCode int

		xPosInt := int(xpos)
		yPosInt := int(ypos)
		for code, RS := range state {
			if g143.InRect(RS, xPosInt, yPosInt) {
				widgetRS = RS
				widgetCode = code
				break
			}
		}

		if widgetCode == 0 {
			// send the last drawn frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, currentWindowFrame, windowRS)
			window.SwapBuffers()
			return
		}

		rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
			widgetRS.OriginX+widgetRS.Width,
			widgetRS.OriginY+widgetRS.Height)

		pieceOfCurrentFrame := imaging.Crop(currentWindowFrame, rectA)
		invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

		ggCtx := gg.NewContextForImage(currentWindowFrame)
		ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

		// show comments
		if widgetCode > 1000 && widgetCode < 2000 && activeTool != DeleteCommentTool {
			instrId := widgetCode - 1000 - 1
			commentObj := comments[instrId]
			cORect := objCoords[widgetCode]
			ggCtx.SetHexColor("#888")
			ggCtx.DrawRoundedRectangle(float64(cORect.OriginX)+10, float64(cORect.OriginY)+10, 460, 300, 10)
			ggCtx.Fill()

			// load font
			fontPath := getDefaultFontPath()
			ggCtx.LoadFontFace(fontPath, 20)

			// header
			cHLY := cORect.OriginY+10+10

			ggCtx.SetHexColor("#fff")
			ggCtx.DrawString(fmt.Sprintf("Comment #%d", instrId+1), float64(cORect.OriginX)+10+20, 
				float64(cHLY)+FontSize)

			ggCtx.DrawRoundedRectangle(float64(cORect.OriginX)+30, float64(cORect.OriginY)+30+FontSize, 430, 1, 1)
			ggCtx.Fill()

			strs := strings.Split(commentObj.Comment, "\n")
			currentY := cHLY + 25
			for _, str := range strs {
				ggCtx.SetHexColor("#fff")
				ggCtx.DrawString(str, float64(cORect.OriginX)+10+20, float64(currentY)+10+FontSize)
				currentY += FontSize + 5
			}

		}

		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()
	}
}
