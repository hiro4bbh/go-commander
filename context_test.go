package gocommander

import (
	"fmt"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestContext(t *testing.T) {
	// Test only option managements by Context, so commander and Command are nil.
	commander := New(nil)
	ctx := newContext(commander, nil)
	opt := NewOptionBool(false)
	ctx.AddOption("opt", opt, "option")
	goassert.New(t, opt).Equal(ctx.GetOption("opt"))
	goassert.New(t, commander.Logger()).Equal(ctx.Logger())
}

func TestContextAddOptionErrors(t *testing.T) {
	commander := New(nil)
	ctx := newContext(commander, nil)
	ctx.AddOption("opt", NewOptionBool(false), "option")
	var caughtPanic interface{}
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		ctx.AddOption("opt", NewOptionBool(true), "option")
	}()
	goassert.New(t, fmt.Errorf("option name opt is already used")).Equal(caughtPanic)
	func() {
		defer func() {
			caughtPanic = recover()
		}()
		ctx.AddOption("help", NewOptionBool(true), "option")
	}()
	goassert.New(t, fmt.Errorf("illegal option name: help")).Equal(caughtPanic)
}

func TestContextParse(t *testing.T) {
	commander := New(nil)
	ctx := newContext(commander, nil)
	ctx.AddOption("opt1", NewOptionBool(false), "option 1")
	ctx.AddOption("opt2", NewOptionBool(false), "option 2")
	ctx.AddOption("opt3", NewOptionBool(true), "option 3")
	ctx.AddOption("opt4", NewOptionBool(true), "option 4")
	goassert.New(t, 3).Equal(goassert.New(t).SucceedNew(ctx.Parse([]string{"opt1", "opt2+", "opt3-"})).(int))
	goassert.New(t, true).Equal(ctx.GetOption("opt1").(*OptionBool).Get())
	goassert.New(t, true).Equal(ctx.GetOption("opt2").(*OptionBool).Get())
	goassert.New(t, false).Equal(ctx.GetOption("opt3").(*OptionBool).Get())
	goassert.New(t, true).Equal(ctx.GetOption("opt4").(*OptionBool).Get())
	goassert.New(t, "[opt1:true opt2:true opt3:false opt4:true]").Equal(ctx.OptionsString())
}

func TestContextParseErrors(t *testing.T) {
	commander := New(nil)
	ctx := newContext(commander, nil)
	ctx.AddOption("opt1", NewOptionBool(false), "option 1")
	goassert.New(t, `unknown option: opt`).ExpectError(ctx.Parse([]string{"opt"}))
	goassert.New(t, `opt1: illegal OptionBool value: =X`).ExpectError(ctx.Parse([]string{"opt1=X"}))
}
