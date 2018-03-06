package gocommander

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hiro4bbh/go-log"
)

// Settings has commander settings.
type Settings struct {
	// Name is the commander's name.
	Name string
	// Copyright is the commander's copyright.
	Copyright string
	// Logger is the commander's logger.
	Logger *golog.Logger
}

var (
	// DefaultName is the default value of CommanderSettings.Name.
	DefaultName = "An go-commander application"
	// DefaultCopyright is the default value of CommanderSettings.Copyright.
	DefaultCopyright = "Copyright 2018- Tatsuhiro Aoshima (hiro4bbh@gmail.com)."
	// The default value of golog.Logger will be golog.Null.
)

// Commander is a manager of commands.
type Commander struct {
	settings *Settings
	ctxs     map[string]*Context
	help     bool
	queue    []string
}

// New returns a new Commander with the given CommanderSettings.
// If settings is nil, then the fields are filled with the default values.
// If a field of settings is zero value, then the field is filled with the corresponding default value.
func New(settings *Settings) *Commander {
	if settings == nil {
		settings = &Settings{
			Name:      DefaultName,
			Copyright: DefaultCopyright,
			Logger:    golog.Null,
		}
	}
	if settings.Name == "" {
		settings.Name = DefaultName
	}
	if settings.Copyright == "" {
		settings.Copyright = DefaultCopyright
	}
	if settings.Logger == nil {
		settings.Logger = golog.Null
	}
	return &Commander{
		settings: settings,
		ctxs:     map[string]*Context{},
	}
}

// Add adds a new command with the given command name and command handlers, and returns its Context.
//
// This function calls panic if the given command name is used or the given command name is illegal.
// The illegal command names are "h", "help", or one ending with "+" or "-" or "=".
func (commander *Commander) Add(name string, cmd Command) *Context {
	if _, ok := commander.ctxs[name]; ok {
		panic(fmt.Errorf("commander has already command %s", name))
	}
	if name == "help" || strings.HasSuffix(name, "+") || strings.HasSuffix(name, "-") || strings.HasSuffix(name, "=") {
		panic(fmt.Errorf("illegal command name: %s", name))
	}
	ctx := newContext(commander, cmd)
	ctx.cmd.Init(ctx)
	commander.ctxs[name] = ctx
	return ctx
}

// Copyright returns the commander's copyright.
func (commander *Commander) Copyright() string {
	return commander.settings.Copyright
}

// Get returns the Context with the given command name.
func (commander *Commander) Get(name string) *Context {
	return commander.ctxs[name]
}

// Logger returns the commander's logger.
func (commander *Commander) Logger() *golog.Logger {
	return commander.settings.Logger
}

// Name returns the commander's name.
func (commander *Commander) Name() string {
	return commander.settings.Name
}

// Help writes the help message to the writer of the commander's logger.
func (commander *Commander) Help() {
	w := commander.Logger().Writer()
	fmt.Fprintf(w, "%s\n%s\n\ncommands:\n", commander.Name(), commander.Copyright())
	names := make([]string, 0, len(commander.ctxs)+1)
	names = append(names, "@help")
	for name, _ := range commander.ctxs {
		names = append(names, name)
	}
	sort.Strings(names[1:])
	for i, name := range names {
		if i == 0 {
			fmt.Fprintf(w, "  %s\tShow this help and exit\n", name)
		} else {
			fmt.Fprintf(w, "  @%s\t%s\n", name, commander.ctxs[name].cmd.Description())
		}
	}
}

// Parse parses the given command line arguments, and returns the next argument index.
//
// This function returns an error in parsing.
func (commander *Commander) Parse(args []string) (int, error) {
	commander.Reset()
	i := 0
	for i < len(args) {
		arg := args[i]
		if !strings.HasPrefix(arg, "@") {
			return -1, fmt.Errorf("expected command name, but got: %s", arg)
		}
		name := arg[1:]
		if name == "help" {
			commander.help = true
			i++
			continue
		}
		ctx := commander.Get(name)
		if ctx == nil {
			return -1, fmt.Errorf("unknown command: @%s", name)
		}
		found := false
		for _, n := range commander.queue {
			if name == n {
				found = true
				break
			}
		}
		if found {
			return -1, fmt.Errorf("cannot run @%s multiple times", name)
		}
		j, err := ctx.Parse(args[i+1:])
		if err != nil {
			return -1, fmt.Errorf("@%s: %s", name, err)
		}
		commander.queue = append(commander.queue, name)
		i += j + 1
	}
	return i, nil
}

// Reset resets the commander and its command states.
func (commander *Commander) Reset() {
	for name, ctx := range commander.ctxs {
		newctx := newContext(commander, ctx.cmd)
		newctx.cmd.Init(newctx)
		commander.ctxs[name] = newctx
	}
	commander.help = false
	commander.queue = []string{}
}

// Run runs the commands in order.
//
// This function stops the execution as soon as a command returns and error, then returns it.
func (commander *Commander) Run() error {
	if commander.help {
		commander.Help()
		return nil
	}
	for _, name := range commander.queue {
		ctx := commander.ctxs[name]
		if ctx.help {
			ctx.Help(name)
			return nil
		}
	}
	for _, name := range commander.queue {
		ctx := commander.ctxs[name]
		if err := ctx.cmd.Run(ctx); err != nil {
			return fmt.Errorf("@%s: %s", name, err)
		}
	}
	return nil
}
