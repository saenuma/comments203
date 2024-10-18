package main

import (
	g143 "github.com/bankole7782/graphics143"
	"image"
)

const (
	FontSize         = 20
	toolBoxW         = 170
	toolBoxH         = 40
	indicatorCircleR = 8
	canvasWidth      = 1200
	canvasHeight     = 600

	SelectImageTool = 101
	AddCommentTool = 102
	DeleteCommentTool = 103
	SaveWorkTool = 104
	OpenFolderTool = 105
	CanvasWidget = 106

	CD_AddBtn = 107
	CD_CloseBtn = 108
	CD_CommentInput = 109
)

type CircleSpec struct {
	X int
	Y int
}

type Comment struct {
	X int
	Y int
	Comment string
}

var (
	objCoords map[int]g143.Rect
	CDObjCoords map[int]g143.Rect
	currentWindowFrame image.Image


	drawnIndicators []CircleSpec
	activeTool int
	currentWorkingImagePath string
	canvasRect g143.Rect
	commentsCount int	

	isUpdateDialog bool
	cursorEventsCount int

	activeX, activeY int
	enteredTxt string
	comments []Comment = make([]Comment, 0)
)

