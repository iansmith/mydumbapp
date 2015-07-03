package main

import (
	_ "fmt"
	"math"

	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
	"github.com/iansmith/tropical/std"
)

var portholeWidth = 80

//
//This demonstrates how to override the "default" behavior for the entire system.
//
//Draw a black X through any interactor that doesn't implement its own drawing function.
func DebugDrawSelf(i tropical.Interactor, canvas tropical.Canvas) {
	canvas.Save()
	canvas.SetStrokeColor("#000000")
	w := i.Width()
	h := i.Height()
	canvas.Rectangle(i.X(), i.Y(), w, h)
	canvas.Stroke()
	canvas.MoveTo(0, 0)
	canvas.LineTo(w, h)
	canvas.Stroke()
	canvas.MoveTo(w, 0)
	canvas.LineTo(0, h)
	canvas.Stroke()
	canvas.Restore()

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

	NewDumbLeaf(parent, "#4169e1", 10, 20, 40, 50)   //royalblue
	NewDumbLeaf(parent, "#1e90ff", 120, 130, 20, 60) //dodgerblue

	//porthole with a child that is read
	porthole := NewPorthole(parent)
	NewIconLeaf(porthole, 0, 0, portholeWidth, portholeWidth, "hudson", 0, 0)

	//random trivial children
	NewTrivialLeaf(parent, 275, 2, 32, 32)
	NewTrivialLeaf(parent, 270, 260, 32, 32)

	//force a drawing pass
	root.Draw()

	//in the background, process things coming through the channel
	//this is the event stream, but it's useful to select on other
	//things as well, like timers, network data, etc
	go func() {
		for {
			select {
			case event := <-ch:
				std.MouseDispatch.Process(event, root)
				root.Draw()
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
	c.SetStrokeColor("#cdc8b1")
	c.Rectangle(0, 0, d.Width(), d.Height())
	c.Stroke()
	std.Default.DrawChildren(d, c)
	c.Restore()
}

//
// DumbLeaf is a simple rectangle with a fill color.  It's entire space is
// covered.
//
type DumbLeaf struct {
	stroke                   bool
	color                    string
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.TreeManipulator
}

func NewDumbLeaf(parent tropical.Interactor, color string, x, y, w, h int) *DumbLeaf {
	result := &DumbLeaf{
		color:           color,
		stroke:          false,
		Coords:          std.NewCoords(x, y, w, h), //set w and h
		TreeManipulator: std.NewTreeManipulator(parent),
	}
	parent.AppendChild(result)
	return result
}

func (d *DumbLeaf) DrawSelf(c tropical.Canvas) {
	if d.stroke {
		c.SetStrokeColor(d.color)
		c.Rectangle(0, 0, d.Width(), d.Height())
		c.Stroke()
	} else {
		c.SetFillColor(d.color)
		c.FillRectangle(0, 0, d.Width(), d.Height())
	}
}

func (d *DumbLeaf) Click() {
	d.stroke = !d.stroke
}

//
// Porthole is a parent that expects to have one child. It masks its child
// with a circle. Width and Height are == and set to portholeWidth
//

type Porthole struct {
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.SingleChild
	startDragX, startDragY   int
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

func (p *Porthole) DragStart() {
	p.startDragX = p.X()
	p.startDragY = p.Y()
}
func (p *Porthole) DragEnd() {
}

func (p *Porthole) Drag(x, y int) {
	p.SetX(p.startDragX + x)
	p.SetY(p.startDragY + y)
}

//
// Have to compensate for the porthole effect when picking.
//
func (p *Porthole) PickSelf(e tropical.Event, pl tropical.PickList) bool {
	centerX := e.X() - (p.Width() / 2)
	centerY := e.Y() - (p.Width() / 2) //has to be square!
	dist := int(math.Floor(math.Sqrt(float64(centerX*centerX) + float64(centerY*centerY))))
	if dist > p.Width()/2 {
		return false
	}
	//
	//this code, if uncommented, will pick children that are inside our hole
	//but we have turned it off because we want to drag _any_ children
	//that we have and NOT engage their picking at all.
	//
	/*if len(p.Children()) > 0 {
		child := p.Children()[0]
		e.Translate(child.X(), child.Y())
		picks, ok := child.(tropical.PicksSelf)
		if !ok {
			std.Default.PickSelf(child, e, pl)
		} else {
			picks.PickSelf(e, pl)
		}
		e.Translate(-child.X(), -child.Y())
	}
	*/

	//add ourself to the picklist, if this was not called just as a test of coords
	if pl != nil {
		pl.AddHit(p)
	}
	return true
}

//
// TrivialLeaf is the smallest leaf possible, code-size-wise.  It will get the
// default drawing behavior, which in this app is defined in DebugDrawSelf.
//
type TrivialLeaf struct {
	tropical.Coords          //implementation => std.Coords
	tropical.TreeManipulator //implementation => std.TreeManipulator
}

func NewTrivialLeaf(parent tropical.Interactor, x, y, w, h int) tropical.Interactor {
	result := &TrivialLeaf{
		Coords:          std.NewCoords(x, y, w, h),
		TreeManipulator: std.NewTreeManipulator(parent),
	}
	parent.AppendChild(result)
	return result
}

//
// Icon pulls an image from the page and always draws it at the same location
// within itself.
//
type IconLeaf struct {
	tropical.Coords            //implementation => std.Coords
	tropical.TreeManipulator   //implementation => std.TreeManipulator
	htmlId                     string
	imageOffsetX, imageOffsetY int
}

func NewIconLeaf(parent tropical.Interactor, x, y, w, h int, htmlId string, offX, offY int) tropical.Interactor {
	result := &IconLeaf{
		Coords:          std.NewCoords(x, y, w, h),
		TreeManipulator: std.NewTreeManipulator(parent),
		htmlId:          htmlId,
		imageOffsetX:    offX,
		imageOffsetY:    offY,
	}
	parent.AppendChild(result)
	return result
}

func (i *IconLeaf) DrawSelf(c tropical.Canvas) {
	c.DrawImageById(i.htmlId, i.imageOffsetX, i.imageOffsetY)
}
