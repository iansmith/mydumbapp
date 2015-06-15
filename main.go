package main

import (
	_ "fmt"
	"math"

	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
	"github.com/iansmith/tropical/std"
)

var portholeWidth = 80

//draw a red X through any interactor that doesn't implement its own drawing function
func DebugDrawSelf(i tropical.Interactor, canvas tropical.Canvas) {
	canvas.SetStrokeColor("cd5c5c")
	canvas.DrawLine(i.X(), i.Y(), i.Width(), i.Height())
	canvas.Stroke()
}

//this can called BEFORE the dom is finished loading.  this is a good place
//to do networkish things, but not a good place to actually manipulate the
//screen. we call domReady() when the DOM is ready.
func main() {
	js.Global.Get("document").Call("addEventListener", "DOMContentLoaded", func(event *js.Object) {
		domReady()
	})
}

func domReady() {
	//for easy debugging
	std.Default.DrawSelf = DebugDrawSelf

	//hook to web page
	root, ch := std.NewRootInteractor("canvas", "#f0fff0", nil) //honeydew for root
	//root here can be treated two ways here, as a tropical.Interactor or as a
	//std.RootInteractor depending on how you want to think of it.  It is
	//goish to use the baz.(foo).bar() notation to access the method bar
	//of more detailed type foo when you know baz is a foo.
	//As a side effect of this new call, the dumb
	parent := NewDumbParent(root.(*std.RootInteractor))

	NewDumbLeaf(parent, "4169e1", 10, 20, 40, 50)   //royalblue
	NewDumbLeaf(parent, "1e90ff", 120, 130, 20, 60) //dodgerblue

	porthole := NewPorthole(parent)

	NewDumbLeaf(porthole, "cd5c5c", 0, 0, portholeWidth, portholeWidth)

	//force a drawing pass
	root.(*std.RootInteractor).Draw()

	go func() {
		for {
			select {
			case <-ch:

			}
		}
	}()
}

//
// DumbParent fills his upper left quandrant with cornsilk
//
type DumbParent struct {
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.TreeManipulator
}

func NewDumbParent(root *std.RootInteractor) *DumbParent {
	result := &DumbParent{
		Coords:          std.NewCoords(0, 0, 350, 250), //set w and h
		TreeManipulator: std.NewTreeManipulator(root),
	}
	root.AppendChild(result)
	return result
}

//we implement drawSelf and want the default child drawing behavior
func (d *DumbParent) DrawSelf(c tropical.Canvas) {
	c.Save()
	c.SetFillColor("#cdc8b1") //cornsilk3
	c.FillRectangle(0, 0, d.Width()/2, d.Height()/2)
	std.Default.DrawChildren(d, c)
	c.Restore()
}

//
// DumbLeaf is a simple rectangle with a fill color.  It's entire space is
// covered.
//
type DumbLeaf struct {
	fillColor                string
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.TreeManipulator
}

func NewDumbLeaf(parent tropical.Interactor, fillColor string, x, y, w, h int) *DumbLeaf {
	result := &DumbLeaf{
		fillColor:       fillColor,
		Coords:          std.NewCoords(x, y, w, h), //set w and h
		TreeManipulator: std.NewTreeManipulator(parent),
	}
	parent.AppendChild(result)
	return result
}

func (d *DumbLeaf) DrawSelf(c tropical.Canvas) {
	c.SetFillColor(d.fillColor)
	c.FillRectangle(0, 0, d.Width(), d.Height())
}

//
// Porthole is a parent that expects to have one child. It masks its child
// with a circle.
//

type Porthole struct {
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.SingleChild
}

func NewPorthole(parent tropical.Interactor) *Porthole {
	result := &Porthole{
		Coords:          std.NewCoords(200, 160, portholeWidth, portholeWidth), //set all the params for now, width and height better be same
		TreeManipulator: std.NewSingleChild(parent),
	}
	parent.AppendChild(result)
	return result
}

func (p *Porthole) DrawSelf(c tropical.Canvas) {
	c.Save()
	c.BeginPath()
	c.Arc(p.Width()/2, p.Height()/2, p.Width()/2, 0, math.Pi*2.0)
	c.Clip()

	//now just call the default for our child
	std.Default.DrawChildren(p, c)
	c.Restore()
}
