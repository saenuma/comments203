package main

import (
	"image"
	"fmt"
	"strings"
	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/AvraamMavridis/randomcolor"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	ObjCoords    *map[int]g143.Rect	
}

func New2dCtx(wWidth, wHeight int, objCoords *map[int]g143.Rect) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx, 
		ObjCoords: objCoords}
	return ctx
}

func Continue2dCtx(img image.Image, objCoords *map[int]g143.Rect) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx, 
		ObjCoords: objCoords}
	return ctx
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	width, height := toolBoxW, toolBoxH
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+10, float64(originY)+FontSize+10)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	objCoords[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawButtonB(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+20, textH+15
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+10, float64(originY)+FontSize)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	(*ctx.ObjCoords)[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) drawTextInput(inputId, originX, originY, inputWidth, height int, values string) g143.Rect {
	ctx.ggCtx.SetHexColor("#444")
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(inputWidth), float64(height))
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX)+2, float64(originY)+2, float64(inputWidth)-4, float64(height)-4)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: inputWidth, Height: height, OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = entryRect

	if len(strings.TrimSpace(values)) != 0 {
		strs := strings.Split(values, "\n")
		currentY := originY
		for _, str := range strs {
			ctx.ggCtx.SetHexColor("#444")
			ctx.ggCtx.DrawString(str, float64(originX+15), float64(currentY)+FontSize)
			currentY += FontSize + 5
		}
	}
	return entryRect
}

func (ctx *Ctx) drawCommentBox(inputId, originX, originY int) g143.Rect {
	colorInHex := randomcolor.GetRandomColorInHex()
	ctx.ggCtx.SetHexColor(colorInHex)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), 80, 40)
	ctx.ggCtx.Fill()

	cBRect := g143.Rect{Width: 80, Height: 40, OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = cBRect

	idToDisplay := inputId - 1000 - 1
	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawString(fmt.Sprintf("#%d", idToDisplay), float64(originX)+10, float64(originY)+FontSize)

	return cBRect
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}

func nextHorizontalCoords(aRect g143.Rect, margin int) (int, int) {
	nextOriginX := aRect.OriginX + aRect.Width + margin
	nextOriginY := aRect.OriginY
	return nextOriginX, nextOriginY
}

func nextVerticalCoords(aRect g143.Rect, margin int) (int, int) {
	nextOriginX := margin
	nextOriginY := aRect.OriginY + aRect.Height + margin
	return nextOriginX, nextOriginY
}