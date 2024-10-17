package main

import (
	g143 "github.com/bankole7782/graphics143"
	"image"
)

const (
	fps              = 10
	fontSize         = 20
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
)

// var objCoords map[g143.Rect]any
var objCoords map[int]g143.Rect
var currentWindowFrame image.Image

type CircleSpec struct {
	X int
	Y int
}

var drawnIndicators []CircleSpec
var activeTool int
var currentWorkingImagePath string
var canvasRect g143.Rect