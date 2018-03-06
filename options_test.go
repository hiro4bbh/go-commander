package gocommander

import (
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestOptionBool(t *testing.T) {
	opt := NewOptionBool(false)
	goassert.New(t, false).Equal(opt.Get())
	goassert.New(t, "false").Equal(opt.String())
	goassert.New(t).SucceedWithoutError(opt.Set("-"))
	goassert.New(t, false).Equal(opt.Get())
	goassert.New(t).SucceedWithoutError(opt.Set(""))
	goassert.New(t, true).Equal(opt.Get())
	goassert.New(t).SucceedWithoutError(opt.Set("+"))
	goassert.New(t, true).Equal(opt.Get())
	goassert.New(t, "true").Equal(opt.String())
	goassert.New(t, `illegal OptionBool value: unknown`).ExpectError(opt.Set("unknown"))
	goassert.New(t, "[+-]").Equal(opt.ValueFormat())
}

func TestOptionString(t *testing.T) {
	opt := NewOptionString("default value")
	goassert.New(t, "default value").Equal(opt.Get())
	goassert.New(t, "\"default value\"").Equal(opt.String())
	goassert.New(t).SucceedWithoutError(opt.Set("=new value"))
	goassert.New(t, "new value").Equal(opt.Get())
	goassert.New(t).SucceedWithoutError(opt.Set(""))
	goassert.New(t, "").Equal(opt.Get())
	goassert.New(t, `illegal OptionString value: \+`).ExpectError(opt.Set("+"))
	goassert.New(t, "").Equal(opt.Get())
	goassert.New(t, `illegal OptionString value: -`).ExpectError(opt.Set("-"))
	goassert.New(t, "").Equal(opt.Get())
	goassert.New(t, "=VALUE").Equal(opt.ValueFormat())
}
