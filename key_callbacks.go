package main

import (
	"strings"
	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)


func CDCharCallback(window *glfw.Window, char rune) {
	wWidth, wHeight := window.GetSize()
	tmp := enteredTxt + string(char)
	tmpParts := strings.Split(tmp, "\n")
	if len(tmpParts[len(tmpParts)-1]) > 45 {
		enteredTxt = enteredTxt + "\n" + string(char)
	} else {
		enteredTxt = tmp
	}


	sIRect := CDObjCoords[CD_CommentInput]
	theCtx := Continue2dCtx(currentWindowFrame, &CDObjCoords)
	theCtx.drawTextInput(CD_CommentInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func CDKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	val := enteredTxt
	if key == glfw.KeyBackspace && len(enteredTxt) != 0 {
		enteredTxt = val[:len(val)-1]
	}	

	if action != glfw.Release {
		return
	}
	wWidth, wHeight := window.GetSize()

	if key == glfw.KeyEnter {
		enteredTxt = val + "\n"
	}

	sIRect := CDObjCoords[CD_CommentInput]
	theCtx := Continue2dCtx(currentWindowFrame, &CDObjCoords)
	theCtx.drawTextInput(CD_CommentInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height,  enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}