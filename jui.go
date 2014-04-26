package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type Service struct {
	Name string
}

type Relations struct {
}

func (r *Relations) PaintRelations(canvas *qml.Common) {
	fmt.Println("painting")
	_ = canvas.Call("getContext", "2D")
	//ctx.Call("moveTo", 0, 0)

}

func newService(engine *qml.Engine, rect qml.Object, win *qml.Window, base qml.Object, name string) {
	ctx := engine.Context().Spawn()
	s := Service{Name: name}
	ctx.SetVar("service", &s)
	obj := rect.Create(ctx)
	obj.Set("parent", win.Root())
}

func run() error {
	qml.Init(nil)

	engine := qml.NewEngine()

	base, err := engine.LoadFile("base.qml")
	if err != nil {
		return err
	}
	rect, err := engine.LoadFile("rect.qml")
	if err != nil {
		return err
	}

	win := base.CreateWindow(nil)
	for i := 0; i < 2; i++ {
		newService(engine, rect, win, base, "a")
	}
	ctx := engine.Context()
	relations := Relations{}
	ctx.SetVar("relations", &relations)

	win.Show()
	win.Wait()

	return nil
}
