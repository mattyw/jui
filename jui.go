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
	origin Service
	end    Service
}

type GoRect struct {
	qml.Object
	Relations []Relation
}

func (r *GoRect) Paint(p *qml.Painter) {
	gl.LineWidth(2.5)
	gl.Color4f(0.0, 0.0, 0.0, 1.0)
	gl.Begin(gl.LINES)
	for _, s := range r.Relations {
		fmt.Println(s)
		ox := gl.Float(s.origin.x)
		oy := gl.Float(s.origin.x)
		ex := gl.Float(s.end.y)
		ey := gl.Float(s.end.y)

		gl.Vertex2f(ex, ey)
		gl.Vertex2f(ox, oy)
		gl.Vertex2f(ox, ey)
		gl.Vertex2f(ex, oy)
	}
	gl.End()
}

type Service struct {
	Name string
	ctx  *qml.Context
	obj  qml.Object
	x    int
	y    int
}

func newService(name string, engine *qml.Engine, rect qml.Object) Service {
	s := Service{Name: name}
	s.ctx = engine.Context().Spawn()
	s.obj = rect.Create(s.ctx)
	return s
}

func (s *Service) Draw(rect qml.Object, win *qml.Window) {
	s.ctx.SetVar("service", s)
	s.obj.Set("parent", win.Root())
}

func (s *Service) NewPos(x, y int) {
	fmt.Printf("new pos %v %v\n", x, y)
	s.x = x
	s.y = y
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

	s1 := newService("a", engine, rect)
	s2 := newService("b", engine, rect)
	services := []Service{s1, s2}
	relation := Relation{s1, s2}
	relations := []Relation{relation}
	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{

		Init: func(r *GoRect, obj qml.Object) {
			r.Object = obj
			r.Relations = relations
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
