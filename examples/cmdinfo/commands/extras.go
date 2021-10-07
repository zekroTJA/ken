package commands

type CmdWithExtras interface {
	ExtraString() string
	ExtraInt() int
}
