package ken

import (
	"encoding/json"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

// CommandInfo contains the parsed application command
// structure of a command as well as additional method
// implementations. This also includes external
// implementations aside from the Command interface.
type CommandInfo struct {
	ApplicationCommand *discordgo.ApplicationCommand `json:"application_command"`
	Implementations    map[string][]interface{}      `json:"implementations"`
}

// String returns the parsed JSON data of the
// CommandInfo.
func (c CommandInfo) String() string {
	return mustToJson(c)
}

// CommandInfoList is a slice of CommandInfo
// elements.
type CommandInfoList []*CommandInfo

// String returns the parsed JSON data of the
// CommandInfoList.
func (c CommandInfoList) String() string {
	return mustToJson(c)
}

func mustToJson(v interface{}) string {
	d, _ := json.MarshalIndent(v, "", "  ")
	return string(d)
}

type KeyTransformerFunc func(string) string

// GetCommandInfo returns a list with information about all
// registered commands.
//
// This call is defaultly cached after first execution
// because it uses reflection to inspect external
// implementations. Because this can be performance
// straining when the method is called frequently,
// the result is cached until the number of commands
// changes.
//
// If you want to disable this behavior, you can set
// Config.DisableCommandInfoCache to true on intializing
// Ken.
func (k *Ken) GetCommandInfo(keyTransformer ...KeyTransformerFunc) (cis CommandInfoList) {
	kt := func(v string) string {
		return v
	}
	if len(keyTransformer) != 0 {
		kt = keyTransformer[0]
	}

	if len(k.cmdInfoCache) == len(k.cmds) {
		return k.cmdInfoCache
	}

	cis = k.collectCommandInfo(kt)

	if !k.opt.DisableCommandInfoCache {
		k.cmdInfoCache = cis
	}

	return
}

func (k *Ken) collectCommandInfo(kt KeyTransformerFunc) (cis CommandInfoList) {
	cis = make(CommandInfoList, 0, len(k.cmds))
	for _, cmd := range k.cmds {
		typ := reflect.TypeOf(cmd)
		impl := make(map[string][]interface{})
		for i := 0; i < typ.NumMethod(); i++ {
			meth := typ.Method(i)
			if meth.IsExported() && meth.Type.NumIn() == 1 {
				vals := meth.Func.Call([]reflect.Value{reflect.ValueOf(cmd)})
				iVals := make([]interface{}, len(vals))
				for i, v := range vals {
					iVals[i] = v.Interface()
				}
				impl[kt(meth.Name)] = iVals
			}
		}
		cis = append(cis, &CommandInfo{
			ApplicationCommand: toApplicationCommand(cmd),
			Implementations:    impl,
		})
	}
	return
}
