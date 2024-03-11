package cards

import (
	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_features"
)

var PLACEHOLDER = "I'm sorry, the ~~-V-O-I-D-~~ is unable to respond."
var PLACEHOLDER_SMALL = "~~-V-O-I-D-~~"

var Feature = features.Feature{
	Name: "cards",
	Desc: "Invites your economy server to participate in trading custom cards.",
	Cmds: []*cmds.Cmd{
		cmds.NewCmd("cards", "This doesn't do anything... (yet)", handleCards),
	},
}

func handleCards(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return cmds.NewCmdRespEmbed("Cards", PLACEHOLDER)
}
