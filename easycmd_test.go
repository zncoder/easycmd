package easycmd

import (
	"strings"
	"testing"

	"github.com/zncoder/easytest"
)

func TestEasycmd(t *testing.T) {
	tt := easytest.New(t)

	var a1b1c1, a1b1, a2 int
	fn1 := func() { a1b1c1++ }
	fn2 := func() { a1b1++ }
	fn3 := func() { a2++ }

	ci := &cmdInfo{children: make(map[string]*cmdInfo)}
	err := addCmd(ci, []string{"a1", "b1", "c1"}, fn1, "a1 b1 c1")
	tt.Nil(err)
	err = addCmd(ci, []string{"a1", "b1"}, fn2, "a1 b1")
	tt.Nil(err)
	err = addCmd(ci, []string{"a2"}, fn3, "a2")
	tt.Nil(err)

	err = addCmd(ci, []string{"a1", "b1"}, func() {}, "dup a1 b1")
	tt.True(err != nil)

	cur, fns, chain := findCmd(ci, []string{"test", "foo"})
	tt.True(len(fns) == 0)
	tt.True(cur == ci)
	tt.DeepEqual([]string{"test"}, chain)

	cur, fns, chain = findCmd(ci, []string{"test", "a1", "b1", "c1"})
	tt.True(len(fns) == 2)
	tt.True(len(cur.children) == 0)
	tt.DeepEqual([]string{"test", "a1", "b1", "c1"}, chain)

	ok := runCmd(cur, fns, strings.Join(chain, " "))
	tt.True(ok)
	tt.True(a1b1 == 1)
	tt.True(a1b1c1 == 1)
}
