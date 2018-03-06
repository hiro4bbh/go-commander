package gocommander

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hiro4bbh/go-log"
)

// Command is the command interface.
type Command interface {
	// Description returns the command description.
	Description() string
	// Init is called, and initializes the command options at the command added to a commander.
	Init(ctx *Context)
	// Run is called on the command run by a commander.
	//
	// This function can return an error in running.
	Run(ctx *Context) error
}

// Context is the command instance with a ContextHandlers.
type Context struct {
	commander *Commander
	cmd       Command
	help      bool
	opts      map[string]Option
	descs     map[string]string
}

func newContext(commander *Commander, cmd Command) *Context {
	return &Context{
		commander: commander,
		cmd:       cmd,
		opts:      map[string]Option{},
		descs:     map[string]string{},
	}
}

// AddOption adds the given opt with the given name.
//
// This function calls panic if the given name is already used or the given name is illegal.
// The illegal name are "h" or "help".
func (ctx *Context) AddOption(name string, opt Option, description string) {
	if _, ok := ctx.opts[name]; ok {
		panic(fmt.Errorf("option name %s is already used", name))
	}
	if name == "help" {
		panic(fmt.Errorf("illegal option name: %s", name))
	}
	ctx.opts[name], ctx.descs[name] = opt, description
}

// GetOption returns the Option with the given option name.
func (ctx *Context) GetOption(name string) Option {
	return ctx.opts[name]
}

// Help writes the help message to the writer of the commander's logger.
func (ctx *Context) Help(cmdname string) {
	w := ctx.commander.Logger().Writer()
	fmt.Fprintf(w, "%s\n%s\n\n@%s: %s\noptions:\n", ctx.commander.Name(), ctx.commander.Copyright(), cmdname, ctx.cmd.Description())
	names := make([]string, 0, len(ctx.opts)+1)
	names = append(names, "help")
	for name := range ctx.opts {
		names = append(names, name)
	}
	sort.Strings(names[1:])
	for i, name := range names {
		if i == 0 {
			fmt.Fprintf(w, "  %s\tShow this help and exit\n", name)
		} else {
			opt, desc, defaultPart := ctx.opts[name], ctx.descs[name], ""
			switch defaultStr := opt.String(); defaultStr {
			case "", "false", "0", "0.0", "\"\"":
			default:
				defaultPart = fmt.Sprintf(" (default %s)", defaultStr)
			}
			fmt.Fprintf(w, "  %s%s\t%s%s\n", name, opt.ValueFormat(), desc, defaultPart)
		}
	}
}

// Logger returns the command's logger.
func (ctx *Context) Logger() *golog.Logger {
	return ctx.commander.Logger()
}

// OptionsString returns the string representation of options.
func (ctx *Context) OptionsString() string {
	names := make([]string, 0, len(ctx.opts))
	for name := range ctx.opts {
		names = append(names, name)
	}
	sort.Strings(names)
	str := "["
	for i, name := range names {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf("%s:%s", name, ctx.opts[name])
	}
	return str + "]"
}

// Parse parses the given command line arguments, and returns the next argument index.
//
// This function returns an error in parsing.
func (ctx *Context) Parse(args []string) (int, error) {
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "@") {
			break
		}
		if arg == "help" {
			ctx.help = true
			i++
			continue
		}
		keyValue := strings.SplitN(arg, "=", 2)
		key, value := keyValue[0], ""
		if len(keyValue) < 2 {
			if strings.HasSuffix(key, "+") || strings.HasSuffix(key, "-") {
				key, value = key[:len(key)-1], key[len(key)-1:]
			}
		} else {
			value = "=" + keyValue[1]
		}
		opt := ctx.opts[key]
		if opt == nil {
			return -1, fmt.Errorf("unknown option: %s", key)
		}
		if err := opt.Set(value); err != nil {
			return -1, fmt.Errorf("%s: %s", key, err)
		}
		i++
	}
	return i, nil
}
