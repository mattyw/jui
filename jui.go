package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"gopkg.in/qml.v0/gl"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type Relation struct {
	origin *Service
	end    *Service
}

type GoRect struct {
	qml.Object
	Relations []Relation
}

func (r *GoRect) Paint(p *qml.Painter) {
	fmt.Println("Painting")
	gl.LineWidth(2.5)
	gl.Color4f(0.0, 0.0, 0.0, 1.0)
	gl.Begin(gl.LINES)
	for _, s := range r.Relations {
		ox := gl.Float(s.origin.x)
		oy := gl.Float(s.origin.y)
		ex := gl.Float(s.end.x)
		ey := gl.Float(s.end.y)
		fmt.Println(ox, oy)
		fmt.Println(ex, ey)

		gl.Vertex2f(ox, oy)
		gl.Vertex2f(ex, ey)
	}
	gl.End()
}

type Service struct {
	Name   string
	ctx    *qml.Context
	obj    qml.Object
	x      int
	y      int
	canvas *GoRect
}

func newService(name string, engine *qml.Engine, rect qml.Object) *Service {
	s := Service{Name: name}
	s.ctx = engine.Context().Spawn()
	s.obj = rect.Create(s.ctx)
	return &s
}

func (s *Service) Draw(rect qml.Object, win *qml.Window) {
	s.ctx.SetVar("service", s)
	s.obj.Set("parent", win.Root())
}

func (s *Service) NewPos(x, y int) {
	fmt.Printf("new pos %v %v %v\n", x, y, s.canvas)
	s.x = x
	s.y = y
	fmt.Println(s.canvas)
	s.canvas.Call("update")
}

func (s *Service) Coords() (gl.Float, gl.Float) {
	x := gl.Float(s.obj.Int("x"))
	y := gl.Float(s.obj.Int("y"))
	return x, y
}

func run() error {
	qml.Init(nil)

	engine := qml.NewEngine()

	rect, err := engine.LoadFile("rect.qml")
	if err != nil {
		return err
	}

	var canvas *GoRect

	s1 := newService("a", engine, rect)
	s2 := newService("b", engine, rect)
	services := []*Service{s1, s2}
	relation := Relation{s1, s2}
	relations := []Relation{relation}
	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{

		Init: func(r *GoRect, obj qml.Object) {
			fmt.Println("registering")
			r.Object = obj
			r.Relations = relations
			canvas = r
			// attatch the canvas for updating - yuck!
			for _, service := range services {
				service.canvas = canvas
				fmt.Println(canvas)
			}
		},
	}})

	base, err := engine.LoadFile("base.qml")
	if err != nil {
		return err
	}

	win := base.CreateWindow(nil)
	for _, s := range services {
		s.Draw(rect, win)
	}

	win.Show()
	win.Wait()

	return nil
}
