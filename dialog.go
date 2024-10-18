package main

import (
	"image"
	"fmt"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)


func drawCommentDialog(window *glfw.Window, currentFrame image.Image) {
	CDObjCoords = make(map[int]g143.Rect)

	wWidth, wHeight := window.GetSize()
	// background image
	img := imaging.AdjustBrightness(currentFrame, -40)
	theCtx := Continue2dCtx(img, &CDObjCoords)

	// dialog rectangle
	dialogWidth := 400
	dialogHeight := 300

	dialogOriginX := (wWidth - dialogWidth) / 2
	dialogOriginY := (wHeight - dialogHeight) / 2

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRectangle(float64(dialogOriginX), float64(dialogOriginY), float64(dialogWidth),
		float64(dialogHeight))
	theCtx.ggCtx.Fill()

	// Add Form
	aFLX, aFLY := dialogOriginX+20, dialogOriginY+20
	theCtx.ggCtx.SetHexColor("#444")
	str1 := fmt.Sprintf("Add Comment #%d", commentsCount+1)
	if isUpdateDialog {
		str1 = "Edit Comment"
	}
	theCtx.ggCtx.DrawString(str1, float64(aFLX), float64(aFLY)+FontSize)

	addBtnOriginX := dialogWidth + dialogOriginX - 160
	addBtnRect := theCtx.drawButtonB(CD_AddBtn, addBtnOriginX, dialogOriginY+20, "Add", "#fff", "#56845A")
	closeBtnX, _ := nextHorizontalCoords(addBtnRect, 10)
	theCtx.drawButtonB(CD_CloseBtn, closeBtnX, addBtnRect.OriginY, "Close", "#fff", "#B75F5F")

	// draw comment box
	_, cIY := nextVerticalCoords(addBtnRect, 20)
	theCtx.drawTextInput(CD_CommentInput, aFLX, cIY, dialogWidth-40, 200, "")

	
	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}