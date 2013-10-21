package engine

import (
	"fmt"
	"io"
	"os"
)


type Handler func(*Job) string

var globalHandlers map[string]Handler

func Register(name string, handler Handler) error {
	if globalHandlers == nil {
		globalHandlers = make(map[string]Handler)
	}
	globalHandlers[name] = handler
	return nil
}

// The Engine is the core of Docker.
// It acts as a store for *containers*, and allows manipulation of these
// containers by executing *jobs*.
type Engine struct {
	root		string
	handlers	map[string]Handler
}

// New initializes a new engine managing the directory specified at `root`.
// `root` is used to store containers and any other state private to the engine.
// Changing the contents of the root without executing a job will cause unspecified
// behavior.
func New(root string) (*Engine, error) {
	if err := os.MkdirAll(root, 0700); err != nil && !os.IsExist(err) {
		return nil, err
	}
	eng := &Engine{
		root:		root,
		handlers:	globalHandlers,
	}
	return eng, nil
}

// Job creates a new job which can later be executed.
// This function mimics `Command` from the standard os/exec package.
func (eng *Engine) Job(name string, args ...string) (*Job, error) {
	handler, exists := eng.handlers[name]
	if !exists || handler == nil {
		return nil, fmt.Errorf("Undefined command; %s", name)
	}
	job := &Job{
		eng:		eng,
		Name:		name,
		Args:		args,
		Stdin:		os.Stdin,
		Stdout:		os.Stdout,
		Stderr:		os.Stderr,
		handler:	handler,
	}
	return job, nil
}
