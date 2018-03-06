package gocommander

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hiro4bbh/go-assert"
	"github.com/hiro4bbh/go-log"
)

type command1 struct {
	opt1, opt2, opt3 bool
}

func (cmd1 *command1) Description() string {
	return "command 1"
}

func (cmd1 *command1) Init(ctx *Context) {
	ctx.AddOption("opt1", NewOptionBool(false), "option 1")
	ctx.AddOption("opt2", NewOptionBool(true), "option 2")
	ctx.AddOption("opt3", NewOptionBool(true), "option 3")
}

func (cmd1 *command1) Run(ctx *Context) error {
	ctx.Logger().Infof("command1 is started")
	cmd1.opt1 = ctx.GetOption("opt1").(*OptionBool).Get()
	cmd1.opt2 = ctx.GetOption("opt2").(*OptionBool).Get()
	cmd1.opt3 = ctx.GetOption("opt3").(*OptionBool).Get()
	return nil
}

type command2 struct{}

func (cmd2 *command2) Description() string {
	return "command 2"
}

func (cmd2 *command2) Init(ctx *Context) {
}

func (cmd2 *command2) Run(ctx *Context) error {
	return fmt.Errorf("always fail")
}

func TestCommander(t *testing.T) {
	// The case that the argument settings of New is nil will be tested in error test cases.
	commander := New(&Settings{})
	goassert.New(t, DefaultName).Equal(commander.Name())
	goassert.New(t, DefaultCopyright).Equal(commander.Copyright())
	cmd1 := &command1{}
	ctx1 := commander.Add("cmd1", cmd1)
	goassert.New(t, ctx1).Equal(commander.Get("cmd1"))
	goassert.New(t, 3).Equal(goassert.New(t).SucceedNew(commander.Parse([]string{"@cmd1", "opt1", "opt2-"})).(int))
	goassert.New(t).SucceedWithoutError(commander.Run())
	goassert.New(t, &command1{
		opt1: true,
		opt2: false,
		opt3: true,
	}).Equal(cmd1)
}

func TestCommanderHelp(t *testing.T) {
	var buf bytes.Buffer
	commander := New(&Settings{Logger: golog.New(&buf, &golog.Parameters{MinLevel: golog.DEBUG, TimeFormat: " "})})
	commander.Add("cmd1", &command1{})
	commander.Add("cmd2", &command2{})
	goassert.New(t).SucceedNew(commander.Parse([]string{"@cmd2", "@help", "@cmd1", "help"}))
	goassert.New(t).SucceedWithoutError(commander.Run())
	goassert.New(t, "An go-commander application\nCopyright 2018- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\ncommands:\n  @help\tShow this help and exit\n  @cmd1\tcommand 1\n  @cmd2\tcommand 2\n").Equal(buf.String())

	(&buf).Reset()
	goassert.New(t).SucceedNew(commander.Parse([]string{"@cmd2", "@cmd1", "help"}))
	goassert.New(t).SucceedWithoutError(commander.Run())
	goassert.New(t, "An go-commander application\nCopyright 2018- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\n@cmd1: command 1\noptions:\n  help\tShow this help and exit\n  opt1[+-]\toption 1\n  opt2[+-]\toption 2 (default true)\n  opt3[+-]\toption 3 (default true)\n").Equal(buf.String())

	(&buf).Reset()
	goassert.New(t).SucceedNew(commander.Parse([]string{"@cmd1"}))
	goassert.New(t).SucceedWithoutError(commander.Run())
	goassert.New(t, " INFO   command1 is started\n").Equal(buf.String())
}

func TestCommanderAddError(t *testing.T) {
	var caughtPanic interface{}
	commander := New(nil)
	commander.Add("cmd1", &command1{})
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		commander.Add("cmd1", &command1{})
	}()
	goassert.New(t, fmt.Errorf("commander has already command cmd1")).Equal(caughtPanic)
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		commander.Add("help", &command1{})
	}()
	goassert.New(t, fmt.Errorf("illegal command name: help")).Equal(caughtPanic)
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		commander.Add("cmd1+", &command1{})
	}()
	goassert.New(t, fmt.Errorf("illegal command name: cmd1+")).Equal(caughtPanic)
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		commander.Add("cmd1-", &command1{})
	}()
	goassert.New(t, fmt.Errorf("illegal command name: cmd1-")).Equal(caughtPanic)
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		commander.Add("cmd1=", &command1{})
	}()
	goassert.New(t, fmt.Errorf("illegal command name: cmd1=")).Equal(caughtPanic)
}

func TestCommanderParseError(t *testing.T) {
	commander := New(nil)
	cmd1 := &command1{}
	commander.Add("cmd1", cmd1)
	goassert.New(t, `expected command name, but got: cmd`).ExpectError(commander.Parse([]string{"cmd", "opt1", "opt2-"}))
	goassert.New(t, `unknown command: @cmd`).ExpectError(commander.Parse([]string{"@cmd", "opt1", "opt2-"}))
	goassert.New(t, `cannot run @cmd1 multiple times`).ExpectError(commander.Parse([]string{"@cmd1", "opt1", "opt2-", "@cmd1"}))
	goassert.New(t, `@cmd1: opt2: illegal OptionBool value: =X`).ExpectError(commander.Parse([]string{"@cmd1", "opt1", "opt2=X"}))
}

func TestCommanderRunError(t *testing.T) {
	commander := New(nil)
	cmd1 := &command1{}
	commander.Add("cmd1", cmd1)
	commander.Add("cmd2", &command2{})
	goassert.New(t, 4).Equal(goassert.New(t).SucceedNew(commander.Parse([]string{"@cmd1", "opt1", "opt2-", "@cmd2"})).(int))
	goassert.New(t, `@cmd2: always fail`).ExpectError(commander.Run())
	goassert.New(t, &command1{
		opt1: true,
		opt2: false,
		opt3: true,
	}).Equal(cmd1)
}
