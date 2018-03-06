package gocommander

import (
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestBoxString(t *testing.T) {
	goassert.New(t, "BoxString(\"str\")").Equal(NewBoxString("str").String())
	goassert.New(t, "str").EqualWithoutError(NewBoxString("str").Unwrap())
	goassert.New(t, "str").Equal(NewBoxString("str").UnwrapOr("defval"))
	goassert.New(t, FilePath("str")).EqualWithoutError(NewBoxString("str").UnwrapFilePath())
	goassert.New(t, FilePath("str")).Equal(NewBoxString("str").UnwrapFilePathOr("defval"))
	goassert.New(t, "BoxString(none)").Equal(NewBoxString(0).String())
	goassert.New(t, `tried to unwrap BoxString\(none\)`).ExpectError(NewBoxString(0).Unwrap())
	goassert.New(t, "defval").Equal(NewBoxString(0).UnwrapOr("defval"))
	goassert.New(t, `tried to unwrap BoxString\(none\)`).ExpectError(NewBoxString(0).UnwrapFilePath())
	goassert.New(t, FilePath("defval")).Equal(NewBoxString(0).UnwrapFilePathOr("defval"))
}

func TestFilePath(t *testing.T) {
	goassert.New(t, FilePath("c")).Equal(FilePath("a/b/c").Base())
	goassert.New(t, FilePath("a/b")).Equal(FilePath("a/b/c").Dir())
	goassert.New(t, ".txt").Equal(FilePath("a/b/c.txt").Ext())
	goassert.New(t, FilePath("a/b/c")).Equal(FilePath("a").Join("b/c"))
}

func TestEnv(t *testing.T) {
	goassert.New(t, "warn").EqualWithoutError(Env("GOLOG_MINLEVEL").Unwrap())
	goassert.New(t, `tried to unwrap BoxString\(none\)`).ExpectError(Env("GOLOG_MINLEVEL0").Unwrap())
}
