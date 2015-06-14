package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical/std"
)

//this can called BEFORE the dom is finished loading.  this is a good place
//to do networkish things, but not a good place to actually manipulate the
//screen. we call domReady() when the DOM is ready.
func main() {
	js.Global.Get("document").Call("addEventListener", "DOMContentLoaded", func(event *js.Object) {
		domReady()
	})
}

func domReady() {
	root := std.NewRootInteractor("canvas")

	//root here can be treated two ways here, as a tropical.Interactor or as a
	//std.RootInteractor depending on how you want to think of it.  It is
	//goish to use the baz.(foo).bar() notation to access the method bar
	//of more detailed type foo when you know baz is a foo.
	root.(*std.RootInteractor).Draw("#eee9e9")
}
