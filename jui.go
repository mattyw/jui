package main

import (
	"flag"
	"fmt"
	"github.com/mattyw/jui/status"
	"gopkg.in/qml.v0"
	"gopkg.in/qml.v0/gl"
	"io/ioutil"
	"os"
	"strings"
)

var (
	deployerFile = flag.String("file", "", "The deployer file to read from")
	bundle       = flag.String("bundle", "", "The bundle to use")
)

func main() {
	flag.Parse()
	if err := run(*deployerFile, *bundle); err != nil {
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
	gl.LineWidth(2.5)
	gl.Color4f(0.0, 0.0, 0.0, 1.0)
	gl.Begin(gl.LINES)
	for _, s := range r.Relations {
		ox := gl.Float(s.origin.x)
		oy := gl.Float(s.origin.y)
		ex := gl.Float(s.end.x)
		ey := gl.Float(s.end.y)

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
	s.y = 600 //opengl and qml y seems to be flipped
	return &s
}

func (s *Service) Draw(rect qml.Object, win *qml.Window) {
	s.ctx.SetVar("service", s)
	s.obj.Set("parent", win.Root())
}

func (s *Service) NewPos(x, y int) {
	// Magic number for the sweet spot for drawing lines
	s.x = x + 40
	s.y = s.canvas.Int("height") - (y + 40)
	s.canvas.Call("update")
}

func (s *Service) Coords() (gl.Float, gl.Float) {
	x := gl.Float(s.obj.Int("x"))
	y := gl.Float(s.obj.Int("y"))
	return x, y
}

func run(deployerFile, bundle string) error {
	qml.Init(nil)

	engine := qml.NewEngine()

	rect, err := engine.LoadFile("rect.qml")
	if err != nil {
		return err
	}

	services := map[string]*Service{}
	relations := []Relation{}

	if deployerFile == "" {
		fmt.Println("juju status")
		var env status.JujuStatus
		env, err = status.GetStatus()
		if err != nil {
			return err
		}
		// build the services
		for name, _ := range env.Services {
			s := newService(name, engine, rect)
			services[name] = s
		}

		// build the relations
		for sName, service := range env.Services {
			for _, r := range service.Relations {
				for _, name := range r {
					relations = append(relations, Relation{services[sName], services[name]})
				}
			}
		}
	} else {
		fmt.Println("deployer")
		data, err := ioutil.ReadFile(deployerFile)
		if err != nil {
			return err
		}
		env, err := status.StatusFromDeployer(bundle, data)
		fmt.Println(env)
		if err != nil {
			return err
		}
		// build the services
		for name, _ := range env.Services {
			s := newService(name, engine, rect)
			services[name] = s
		}

		// build the relations
		fmt.Printf("FFFFFF %v\n", env.Relations)
		for _, r := range env.Relations {
			fmt.Printf("%v %v\n", r[0], r[1])
			r[0] = strings.Split(r[0], ":")[0] //HACK
			r[1] = strings.Split(r[1], ":")[0] //HACK
			rA, ok := services[r[0]]
			if !ok {
				continue
			}
			rB, ok := services[r[1]]
			if !ok {
				continue
			}
			relations = append(relations, Relation{rA, rB})
		}
	}

	var canvas *GoRect

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{

		Init: func(r *GoRect, obj qml.Object) {
			fmt.Println("registering")
			r.Object = obj
			r.Relations = relations
			canvas = r
			// attatch the canvas for updating - yuck!
			for _, service := range services {
				service.canvas = canvas
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
