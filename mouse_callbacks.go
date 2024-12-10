package main

import (
	"strings"
	"fmt"
	"os"
	"math"
	"path/filepath"
	"slices"
	"encoding/json"
	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/fogleman/gg"
	"github.com/disintegration/imaging"
)

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt, yPosInt := int(xPos), int(yPos)

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
		if currentWorkingImagePath != "" {
			return
		}
		
		imagePath := PickImageFile()
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

		theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

		theCtx.ggCtx.SetHexColor("#fff")
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

	case CanvasWidget:
		if activeTool == AddCommentTool {
			
			for _, obj := range comments {
				if g143.InRect(obj.getRect(), xPosInt, yPosInt) {
					return
				}
			}

			activeX, activeY = xPosInt, yPosInt
			drawCommentDialog(window, currentWindowFrame)
			window.SetMouseButtonCallback(CDMouseBtnCallback)
			window.SetCursorPosCallback(getHoverCB(CDObjCoords))
			window.SetKeyCallback(CDKeyCallback)
			window.SetCharCallback(CDCharCallback)

		} else if activeTool == DeleteCommentTool {

			var tmp []Comment
			for i, obj := range comments {
				if g143.InRect(obj.getRect(), xPosInt, yPosInt) {
					tmp = slices.Delete(comments, i, i+1)
					break
				}
			}
			comments = tmp

			drawMainWindow(window)
			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// quick hover effect
			window.SetCursorPosCallback(getHoverCB(objCoords))
			window.SetKeyCallback(nil)
			window.SetCharCallback(nil)

		}

	case SaveWorkTool:
		if currentWorkingImagePath == "" {
			return
		}

		imgPath := filepath.Join(os.TempDir(), "img"+filepath.Ext(currentWorkingImagePath))
		rawImg, _ := os.ReadFile(currentWorkingImagePath)
		os.WriteFile(imgPath, rawImg, 0777)

		rawJSON, _ := json.Marshal(comments)
		jsonPath := filepath.Join(os.TempDir(), "comments.json")
		os.WriteFile(jsonPath, rawJSON, 0777)

		finalArchive := filepath.Join(rootPath, UntestedRandomString(7) + ".c203f")
		writeTar([]string{imgPath, jsonPath}, finalArchive)

	case OpenFolderTool:
		ExternalLaunch(rootPath)

	default:

	}
}


func CDMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt, yPosInt := int(xPos), int(yPos)

	wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range CDObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	switch widgetCode {
	case CD_CloseBtn:
		drawMainWindow(window)
		// register the ViewMain mouse callback
		window.SetMouseButtonCallback(mouseBtnCallback)
		// unregister the keyCallback
		window.SetKeyCallback(nil)
		window.SetCharCallback(nil)
		// window.SetScrollCallback(FirstUIScrollCallback)
		window.SetCursorPosCallback(getHoverCB(objCoords))

		activeX, activeY = -1, -1

	case CD_AddBtn:
		if enteredTxt == "" {
			return
		}

		comments = append(comments, Comment{activeX, activeY, enteredTxt})
		enteredTxt = ""

		drawMainWindow(window)
		// register the ViewMain mouse callback
		window.SetMouseButtonCallback(mouseBtnCallback)
		// unregister the keyCallback
		window.SetKeyCallback(nil)
		window.SetCharCallback(nil)
		// window.SetScrollCallback(FirstUIScrollCallback)
		window.SetCursorPosCallback(getHoverCB(objCoords))

		activeX, activeY = -1, -1

	case CD_PasteBtn:
		theCtx := Continue2dCtx(currentWindowFrame, &CDObjCoords)	
		ctrlState := window.GetKey(glfw.KeyLeftControl)

		if ctrlState == glfw.Release {
			toPasteStr := glfw.GetClipboardString()
			if len(strings.TrimSpace(toPasteStr)) == 0 {
				return
			}

			toPasteStr = strings.ReplaceAll(toPasteStr, "\n", " ")
			lines := make([]string, 0)
			toPasteStrLen := len(toPasteStr)
			rem := math.Mod(float64(toPasteStrLen), float64(45))
			for i := range(toPasteStrLen/45) {
				tmp := toPasteStr[(i*45): (i+1)*45]
				lines = append(lines, tmp)
			}

			if int(rem) != 0 {
				lines = append(lines, toPasteStr[toPasteStrLen-int(rem):])
			}

			enteredTxt = strings.Join(lines, "\n")

			sIRect := CDObjCoords[CD_CommentInput]
			theCtx.drawTextInput(CD_CommentInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

		} else if ctrlState == glfw.Press {
			enteredTxt = ""
			sIRect := CDObjCoords[CD_CommentInput]
			theCtx.drawTextInput(CD_CommentInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)			
		}

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = theCtx.ggCtx.Image()
	}

}
