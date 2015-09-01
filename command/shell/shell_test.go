package shell

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"../command"
)

func TestNew(t *testing.T) {
	settingText := `
[[shells]]
name = "test"
command = "option"
workDir = "./"

[[shells]]
name = "testecho"
command = "echo aaa"
workDir = "./"

[[shells]]
name = "catecho"
command = "echo \"%s\""
workDir = "./"
`
	Convey("New", t, func() {
		Convey("correct", func() {
			cmd := New(settingText)

			Convey("load", func() {
				So(cmd, ShouldNotBeNil)
				So(len(cmd.commands), ShouldEqual, 3)
			})

			Convey("command", func() {
				scmd := cmd.commands["test"]
				So(scmd, ShouldNotBeNil)
				So(scmd.Name, ShouldEqual, "test")
				So(scmd.Command, ShouldEqual, "option")
				So(scmd.WorkDir, ShouldEqual, "./")

				scmd2 := cmd.commands["testecho"]
				So(scmd2, ShouldNotBeNil)
				So(scmd2.Name, ShouldEqual, "testecho")
				So(scmd2.Command, ShouldEqual, "echo aaa")
				So(scmd2.WorkDir, ShouldEqual, "./")
			})

			Convey("execute", func() {
				Convey("no option", func() {
					order := &command.Order{}
					order.Data = "testecho"
					ret := cmd.Execute(*order)
					So(ret, ShouldEqual, "aaa")
				})

				Convey("with option", func() {
					order := &command.Order{}
					order.Data = "catecho execute test"
					ret := cmd.Execute(*order)
					So(ret, ShouldEqual, "execute test")
				})

				Convey("no command", func() {
					order := &command.Order{}
					order.Data = "nothing"
					ret := cmd.Execute(*order)
					So(ret, ShouldEqual, "no shell command")
				})
			})
		})
	})
}
