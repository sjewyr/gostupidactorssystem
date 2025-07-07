// Basic Actors package, main subject - FuckSystem ("Supervisor") managing Actors ("Gays"), to create a new supervisor call NewFuckSystem()
package fuck

import "fmt"

type Message = [4]byte

const (
	DummyMessage byte = iota
	AnalSex
	AnalSexReciever
)

type GayUndefinedHandlerBehaviour = byte

const (
	Print GayUndefinedHandlerBehaviour = iota
	Drop
	Panic
)

type FuckSystem struct {
	Gays  []*Gay
	state [4]byte
}

// Creates a system of Actors ("Gays"), powered by Handlers, and Messages, create new Actors ("Gays") by AddGay method
func NewFuckSystem() FuckSystem {
	return FuckSystem{}
}

// Adds a Gay in a System, setting it's default Handler to Print (by default)
// Overriding it is possible for all new actors with base Print, Drop and Panic (via ChangeDefaultBehaviourForNewGays method) or to an individual actors by registering a custom Handler to 0 message
// Default message is considered 0
// Registering Handlers for Actors ("Gays") is possible via RegisterHandler method
func (f *FuckSystem) AddGay() byte {
	id := f.state[0]
	gay := Gay{id: id, anal: make(chan Message), sys: f, handlers: make(map[byte]GayHandler)}
	f.Gays = append(f.Gays, &gay)
	switch f.state[1] {
	case Print:
		f.RegisterHandler(id, 0, HandlerFunc(BaseGayHandlerPrint))
	case Drop:
		f.RegisterHandler(id, 0, HandlerFunc(BaseGayHandlerDrop))
	case Panic:
		f.RegisterHandler(id, 0, HandlerFunc(BaseGayHandlerPanic))
	}
	go gay.Serve()
	f.state[0] += 1
	return id
}

// Changes default behaviour for all new actors in a given system
// Variants are Print (default), Drop and Panic
func (f *FuckSystem) ChangeDefaultBehaviourForNewGays(gid GayUndefinedHandlerBehaviour) {
	f.state[1] = gid
}
func (f *FuckSystem) SendMessage(ia byte, m Message) {
	f.Gays[ia].anal <- m
}

// Registers a handler to an actor ia (must be present in a system f)
// message referring to id of a handled message id, h is a func wrapped in a HandlerFunc
// It is possible to define custom default handler by overriding 0 message
func (f *FuckSystem) RegisterHandler(ia byte, message byte, h GayHandler) {
	f.Gays[ia].handlers[message] = h
}

// Is a basic Actor, powered by Handlers, contains own state within state 4 bytes-array, Handler can be defined by wrapping func in HandlerFunc
// Message 0 is considered default message, overriding default behaviour is possible
type Gay struct {
	id       byte
	anal     chan Message
	sys      *FuckSystem
	state    [4]byte
	handlers map[byte]GayHandler
}

type GayHandler interface {
	Handle(m Message, g *Gay)
}

// Is a wrapper for functions to be a handler for Actors ("Gays")
type HandlerFunc func(msg Message, g *Gay)

func (f HandlerFunc) Handle(msg Message, g *Gay) {
	f(msg, g)
}

func BaseGayHandlerPanic(msg Message, g *Gay) {
	panic(fmt.Sprintf("Not handled message id: %d for gay: %d\n", msg[0], g.id))
}
func BaseGayHandlerPrint(msg Message, g *Gay) {
	fmt.Printf("Not handled message id: %d for gay: %d\n", msg[0], g.id)
}
func BaseGayHandlerDrop(msg Message, g *Gay) {
}

func (g *Gay) Serve() {
	for {
		msg := <-g.anal
		typ := msg[0]
		handler, found := g.handlers[typ]
		if !found {
			handler = g.handlers[0]
			msg[0] = 0
		}
		handler.Handle(msg, g)
	}
}
