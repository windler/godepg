package appcontext

import (
	"github.com/windler/cli"
)

//AppContext provides app flags
type AppContext struct {
	Context      *cli.Context
	Strings      map[string]string
	Ints         map[string]int
	Stringslices map[string][]string
	Bools        map[string]bool
}

//GetStringFlag gets the value of a string flag
func (ac AppContext) GetStringFlag(flag string) string {
	if res, found := ac.Strings[flag]; found {
		return res
	}
	return ac.Context.String(flag)
}

//GetStringSliceFlag gets all values for a slice flag
func (ac AppContext) GetStringSliceFlag(flag string) []string {
	if res, found := ac.Stringslices[flag]; found {
		return res
	}
	return ac.Context.StringSlice(flag)
}

//GetIntFlag gets an int-value for a flag
func (ac AppContext) GetIntFlag(flag string) int {
	if res, found := ac.Ints[flag]; found {
		return res
	}
	return ac.Context.Int(flag)
}

//GetBoolFlag gets a bool-value for a flag
func (ac AppContext) GetBoolFlag(flag string) bool {
	if res, found := ac.Bools[flag]; found {
		return res
	}
	return ac.Context.Bool(flag)
}

//SetStringFlag sets a string flag
func (ac AppContext) SetStringFlag(flag, value string) {
	ac.Strings[flag] = value
}

//SetStringSliceFlag sets a stringslice flag
func (ac AppContext) SetStringSliceFlag(flag string, value []string) {
	ac.Stringslices[flag] = value
}

//SetIntFlag sets a int flag
func (ac AppContext) SetIntFlag(flag string, value int) {
	ac.Ints[flag] = value
}

//SetBoolFlag sets a bool flag
func (ac AppContext) SetBoolFlag(flag string, value bool) {
	ac.Bools[flag] = value
}
