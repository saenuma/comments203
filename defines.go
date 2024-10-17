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

	PencilWidget        = 101
	CanvasWidget        = 102
	SymmLineWidget      = 103
	LeftSymmWidget      = 104
	RefLineWidget       = 105
	ClearRefLinesWidget = 106
	SaveWidget          = 107
	OpenWDWidget        = 108
)

// var objCoords map[g143.Rect]any
var objCoords map[int]g143.Rect
var currentWindowFrame image.Image