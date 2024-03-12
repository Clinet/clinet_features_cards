package cards

import (
	"fmt"

	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_features"
	"github.com/Clinet/clinet_services"
	"github.com/Clinet/clinet_storage"
)

var PLACEHOLDER = "I'm sorry, the ~~-V-O-I-D-~~ is unable to respond."
var PLACEHOLDER_SMALL = "~~-V-O-I-D-~~"

var Feature = features.Feature{
	Name: "cards",
	Desc: "Invites your economy server to participate in trading custom cards.",
	Cmds: []*cmds.Cmd{
		cmds.NewCmd("cards", "Card inventory management", nil).AddSubCmds(
			cmds.NewCmd("dumpctx", "Dumps the context for testing", cmdDumpCtx).AddArgs(
				cmds.NewCmdArg("user", "user", cmds.ArgTypeUser),
				cmds.NewCmdArg("int", "int", 0),
				cmds.NewCmdArg("string", "string", ""),
			),
			cmds.NewCmd("list", "Lists every single card", cmdList).AddArgs(
				cmds.NewCmdArg("page", "Nest into the pages of cards", 1),
			),
			cmds.NewCmd("deck", "Lists all cards in one's deck", cmdDeck).AddArgs(
				cmds.NewCmdArg("page", "Nest into the pages of cards", 1),
				cmds.NewCmdArg("user", "Who to view the deck of (default: you)", cmds.ArgTypeUser),
			),
			cmds.NewCmd("create", "Creates a new card", cmdCreate).AddArgs(
				cmds.NewCmdArg("name", "Name of the card", "").SetRequired(),
				cmds.NewCmdArg("color", "Color of the card for embeds", "").SetRequired(),
				cmds.NewCmdArg("series", "Series number, 0 for one-offs", -1).SetRequired(),
				cmds.NewCmdArg("front", "URL to front of card", "").SetRequired(),
				cmds.NewCmdArg("back", "URL to back of card", "").SetRequired(),
			),
			cmds.NewCmd("update", "Updates an existing card", cmdUpdate).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to update", "").SetRequired(),
				cmds.NewCmdArg("name", "Name of the card", ""),
				cmds.NewCmdArg("color", "Color of the card for embeds", ""),
				cmds.NewCmdArg("series", "Series number, 0 for one-offs", -1),
				cmds.NewCmdArg("front", "URL to front of card", ""),
				cmds.NewCmdArg("back", "URL to back of card", ""),
			),
			cmds.NewCmd("delete", "Wipes a card from existence", cmdDelete).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to delete", "").SetRequired(),
			),
			cmds.NewCmd("tear", "Tears one or more cards from one's deck", cmdTear).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to tear", "").SetRequired(),
				cmds.NewCmdArg("count", "How many to tear", 1),
				cmds.NewCmdArg("user", "Who to tear the card from (default: you)", cmds.ArgTypeUser),
				cmds.NewCmdArg("confirm", "Tears the card permanently", false).SetRequired(),
			),
			cmds.NewCmd("view", "Renders the specified card", cmdView).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to render", "").SetRequired(),
				cmds.NewCmdArg("flip", "Flips to the details", false),
				cmds.NewCmdArg("op", "Show it off without owning it", false),
			),
			cmds.NewCmd("give", "Gives a card away", cmdGive).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to give", "").SetRequired(),
				cmds.NewCmdArg("user", "Who to give the card to", cmds.ArgTypeUser).SetRequired(),
				cmds.NewCmdArg("count", "How many to give", 1),
				cmds.NewCmdArg("op", "Pull the card from the void", false),
			),
			cmds.NewCmd("request", "Request to trade for one's card", cmdRequest).AddArgs(
				cmds.NewCmdArg("model", "Model of the card to request", "").SetRequired(),
				cmds.NewCmdArg("user", "Who to request the card from", cmds.ArgTypeUser).SetRequired(),
			),
			cmds.NewCmd("deny", "Deny the unequal trade with someone", cmdDeny).AddArgs(
				cmds.NewCmdArg("user", "Who to deny the trade of", cmds.ArgTypeUser),
			),
			cmds.NewCmd("accept", "Accept the balanced trade with someone", cmdAccept).AddArgs(
				cmds.NewCmdArg("user", "Who to accept the trade of", cmds.ArgTypeUser),
			),
		),
	},
	Init: Init,
}

func isAdmin(ctx *cmds.CmdCtx) bool {
	ctxPerms, err := ctx.Service.GetUserPerms(ctx.Server.ServerID, ctx.Channel.ChannelID, ctx.User.UserID)
	if err != nil {
		return false
	}
	if ctx.User.UserID == "127184346334494721" { //@joshuadoes on Discord
		return true
	}
	return ctxPerms.CanAdministrate()
}

func cmdList(ctx *cmds.CmdCtx) *cmds.CmdResp {
	if !isAdmin(ctx) {
		return cmds.NewCmdRespEmbed("Operator Denial", "You must be an administrator to use this command!")
	}

	page := ctx.GetArg("page").GetInt()

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(cards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to list yet!")
	}

	pagedCards, totalPages := Paginate(cards, page-1, 5)
	if len(pagedCards) <= 0 {
		return cmds.NewCmdRespEmbed("Error!", "No cards found on that page!")
	}

	msg := services.NewMessage().SetFooter(fmt.Sprintf("Page %d/%d", page, totalPages))
	for i := 0; i < len(pagedCards); i++ {
		card := pagedCards[i]
		if i == 0 {
			msg.SetImage(card.Front)
		}
		if i == 1 {
			msg.SetThumbnail(card.Front)
		}
		msg.AddField(card.Name, "Series #" + card.Series + "\n" + card.Model)
	}
	return cmds.CmdRespFromMsg(msg)
}

func cmdDeck(ctx *cmds.CmdCtx) *cmds.CmdResp {
	page := ctx.GetArg("page").GetInt()
	user := ctx.GetArg("user").GetUser()

	userID := ctx.User.UserID
	if user != nil && user.UserID != "" {
		if !isAdmin(ctx) {
			return cmds.NewCmdRespEmbed("Error!", "You must be an administrator to view another user's deck!")
		}
		userID = user.UserID
	}

	srcDeck := GetCardsFromStorageUser(Storage, ctx.Server.ServerID + userID)
	deck, totalPages := Paginate(srcDeck, page-1, 5)
	if len(deck) <= 0 {
		return cmds.NewCmdRespEmbed("Error!", "No cards found on that page!")
	}

	srcCards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(srcCards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to list yet!")
	}
	pagedCards := make([]*Card, len(deck))

	for i := 0; i < len(deck); i++ {
		for j := 0; j < len(srcCards); j++ {
			if srcCards[j].Model == deck[i] {
				pagedCards[i] = srcCards[j]
			}
		}
	}

	msg := services.NewMessage().SetFooter(fmt.Sprintf("Page %d/%d", page, totalPages))
	for i := 0; i < len(pagedCards); i++ {
		card := pagedCards[i]
		if i == 0 {
			msg.SetImage(card.Front)
		}
		if i == 1 {
			msg.SetThumbnail(card.Front)
		}
		msg.AddField(card.Name, "Series #" + card.Series + "\n" + card.Model)
	}
	return cmds.CmdRespFromMsg(msg)
}

func cmdCreate(ctx *cmds.CmdCtx) *cmds.CmdResp {
	if !isAdmin(ctx) {
		return cmds.NewCmdRespEmbed("Operator Denial", "You must be an administrator to use this command!")
	}

	name := ctx.GetArg("name").GetString()
	series := ctx.GetArg("series").GetString()
	color := GetColor(ctx.GetArg("color").GetString())
	front := ctx.GetArg("front").GetString()
	back := ctx.GetArg("back").GetString()

	card := &Card{
		Model: RandomStringUpper(5),
		Name: name,
		Series: series,
		Color: color,
		Front: front,
		Back: back,
	}

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	cards = append(cards, card)
	Storage.ServerSet(ctx.Server.ServerID, "cards", cards)
	return cmds.NewCmdRespEmbed("Card created!", fmt.Sprintf("You can view it with `/cards view model:%s`", card.Model))
}

func cmdUpdate(ctx *cmds.CmdCtx) *cmds.CmdResp {
	if !isAdmin(ctx) {
		return cmds.NewCmdRespEmbed("Operator Denial", "You must be an administrator to use this command!")
	}

	model := ctx.GetArg("model").GetString()
	name := ctx.GetArg("name").GetString()
	series := ctx.GetArg("series").GetString()
	colorPtr := ctx.GetArg("color").GetString()
	front := ctx.GetArg("front").GetString()
	back := ctx.GetArg("back").GetString()

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(cards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to update yet!")
	}

	for i := 0; i < len(cards); i++ {
		if cards[i].Model == model {
			if name != "" {
				cards[i].Name = name
			}
			if series != "" {
				cards[i].Series = series
			}
			if colorPtr != "" {
				cards[i].Color = GetColor(colorPtr)
			}
			if front != "" {
				cards[i].Front = front
			}
			if back != "" {
				cards[i].Back = back
			}
			Storage.ServerSet(ctx.Server.ServerID, "cards", cards)
			return cmds.NewCmdRespEmbed("Card updated!", fmt.Sprintf("You can view it with `/cards view model:%s`", model))
		}
	}
	return cmds.NewCmdRespEmbed("Error!", "The card model you specified was not found!")
}

func cmdDelete(ctx *cmds.CmdCtx) *cmds.CmdResp {
	if !isAdmin(ctx) {
		return cmds.NewCmdRespEmbed("Operator Denial", "You must be an administrator to use this command!")
	}

	model := ctx.GetArg("model").GetString()

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(cards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to delete yet!")
	}

	index := -1
	for i := 0; i < len(cards); i++ {
		if cards[i].Model == model {
			index = i
			break
		}
	}
	if index <= -1 {
		return cmds.NewCmdRespEmbed("Error!", "The card model you specified was not found!")
	}

	cards = append(cards[:index], cards[index+1:]...)
	Storage.ServerSet(ctx.Server.ServerID, "cards", cards)

	Storage.ForEachUser(func(key string, val *storage.StorageObject) bool {
		//TODO: Convert (serverID+userID) to (serverID+":"+userID) for easy splitting, this is horrendous and I made a fatal mistake
		//Potentially also switch to userID:serverID instead? Or serviceName:userID:serverID? Rewrite the methods to support it, don't make it optional
		serverID := string(key[:len(ctx.Server.ServerID)])
		if serverID != ctx.Server.ServerID {
			return false
		}
		userID := string(key[len(serverID):])
		if userID == "" {
			return false
		}

		DeleteCardsFromStorageUser(Storage, serverID + userID, model)
		return false
	})

	return cmds.NewCmdRespEmbed("Card deleted!", "It can no longer be viewed, given or traded")
}

func cmdView(ctx *cmds.CmdCtx) *cmds.CmdResp {
	model := ctx.GetArg("model").GetString()
	flip := ctx.GetArg("flip").GetBool()
	op := ctx.GetArg("op").GetBool()

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(cards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to view yet!")
	}

	deck := GetCardsFromStorageUser(Storage, ctx.Server.ServerID + ctx.User.UserID)
	hasModel := false
	for i := 0; i < len(deck); i++ {
		if deck[i] == model {
			hasModel = true
			break
		}
	}
	if !hasModel && !op {
		return cmds.NewCmdRespEmbed("Error!", "You cannot view a card without owning it!")
	}

	for i := 0; i < len(cards); i++ {
		if cards[i].Model == model {
			if flip {
				return cards[i].RenderBack()
			}
			return cards[i].RenderFront()
		}
	}
	return cmds.NewCmdRespEmbed("Error!", "The card model you specified was not found!")
}

func cmdGive(ctx *cmds.CmdCtx) *cmds.CmdResp {
	model := ctx.GetArg("model").GetString()
	user := ctx.GetArg("user").GetUser()
	count := ctx.GetArg("count").GetInt()
	op := ctx.GetArg("op").GetBool()

	if count < 1 {
		return cmds.NewCmdRespEmbed("Error!", "You must specify one or more cards!")
	}

	if op && isAdmin(ctx) {
		srcCards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
		found := false
		for i := 0; i < len(srcCards); i++ {
			if srcCards[i].Model == model {
				found = true
				break
			}
		}
		if !found {
			return cmds.NewCmdRespEmbed("Error!", "That card model doesn't exist to be given!")
		}

		dstCards := GetCardsFromStorageUser(Storage, ctx.Server.ServerID + user.UserID)
		for i := 0; i < count; i++ {
			dstCards = append(dstCards, model)
		}
		Storage.UserSet(ctx.Server.ServerID + user.UserID, "deck", dstCards)
		plural := "Card"
		if count > 1 {
			plural += "s"
		}
		return cmds.NewCmdRespEmbed(plural + " given!", "You drew them from the void.")
	}

	srcCards := GetCardsFromStorageUser(Storage, ctx.Server.ServerID + ctx.User.UserID)
	if len(srcCards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to give yet!")
	}

	toGive := make([]int, 0)
	counter := 0
	for i := 0; i < len(srcCards); i++ {
		if srcCards[i] == model {
			toGive = append(toGive, i)
			counter++
			if counter == count {
				break
			}
		}
	}
	if counter < count {
		return cmds.NewCmdRespEmbed("Error!", "Was there a typo? You don't have enough of that card to give yet!")
	}

	dstCards := GetCardsFromStorageUser(Storage, ctx.Server.ServerID + user.UserID)
	for i := 0; i < len(toGive); i++ {
		dstCards = append(dstCards, srcCards[toGive[i]-i])
		toRemove := toGive[i]-i
		srcCards = append(srcCards[:toRemove], srcCards[toRemove+1:]...)
	}

	Storage.UserSet(ctx.Server.ServerID + ctx.User.UserID, "deck", srcCards)
	Storage.UserSet(ctx.Server.ServerID + user.UserID, "deck", dstCards)
	plural := "Card"
	if count > 1 {
		plural += "s"
	}
	return cmds.NewCmdRespEmbed(plural + " given!", "You forfeited your possession.")
}

func cmdTear(ctx *cmds.CmdCtx) *cmds.CmdResp {
	model := ctx.GetArg("model").GetString()
	count := ctx.GetArg("count").GetInt()
	user := ctx.GetArg("user").GetUser()
	confirm := ctx.GetArg("confirm").GetBool()

	if count < 1 {
		return cmds.NewCmdRespEmbed("Error!", "You must specify one or more cards!")
	}

	cards := GetCardsFromStorageServer(Storage, ctx.Server.ServerID)
	if len(cards) == 0 {
		return cmds.NewCmdRespEmbed("Error!", "You don't have any cards to tear yet!")
	}

	index := -1
	for i := 0; i < len(cards); i++ {
		if cards[i].Model == model {
			index = i
			break
		}
	}
	if index <= -1 {
		return cmds.NewCmdRespEmbed("Error!", "The card model you specified was not found!")
	}

	userID := ctx.Server.ServerID
	if user != nil {
		if !isAdmin(ctx) {
			return cmds.NewCmdRespEmbed("Error!", "You must be an administrator to tear someone else's card!")
		}
		userID += user.ID
	} else {
		userID += ctx.User.UserID
	}

	if !confirm {
		plural := "it"
		if count > 1 {
			plural = "them"
		}
		return cmds.NewCmdRespEmbed("Tear", fmt.Sprintf("Are you sure? Add `confirm:True` to the tear command to tear %s for good!", plural))
	}

	deleted := DeleteCardsFromStorageUser(Storage, userID, model, count)
	if deleted <= 0 {
		return cmds.NewCmdRespEmbed("Error!", "There are no cards in this deck to be torn.")
	}
	plural := "that card"
	if count > 1 {
		plural = "those cards"
	}
	return cmds.NewCmdRespEmbed("Tear", fmt.Sprintf("You tore those cards up permanently!", plural))
}

func cmdRequest(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return cmds.NewCmdRespEmbed("Request", PLACEHOLDER)
}

func cmdDeny(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return cmds.NewCmdRespEmbed("Deny", PLACEHOLDER)
}

func cmdAccept(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return cmds.NewCmdRespEmbed("Accept", PLACEHOLDER)
}

func cmdDumpCtx(ctx *cmds.CmdCtx) *cmds.CmdResp {
	dump := services.NewMessage().
		SetContent(fmt.Sprintf("Dump of ctx (*cmds.CmdCtx):\n```JSON\n%s\n```", ctx)).
		SetColor(0xFF0000)
	return cmds.CmdRespFromMsg(dump).SetReady(true)
}