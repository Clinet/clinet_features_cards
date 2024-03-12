package cards

import (
	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_services"
	"github.com/Clinet/clinet_storage"
)

//Card holds details about a card
type Card struct {
	Model	string	`json:"model"`	//Model number of this card, so it can be identified for edits
	Name	string	`json:"name"`	//Name of this card
	Series  string  `json:"series"` //Series of this card
	Color	int		`json:"color"`	//Color of this card
	Front	string	`json:"front"`	//URL to the front image of this card
	Back	string	`json:"back"`	//URL to the back image of this card
}

func (c *Card) render() *cmds.CmdResp {
	card := services.NewMessage().
		SetTitle(c.Name).
		SetContent("Series #" + c.Series).
		SetFooter("Model: " + c.Model).
		SetColor(c.Color)
	return cmds.CmdRespFromMsg(card)
}

func (c *Card) RenderFront() *cmds.CmdResp {
	front := c.render()
	front.Image = c.Front
	front.Thumbnail = c.Back
	return front
}

func (c *Card) RenderBack() *cmds.CmdResp {
	back := c.render()
	back.Image = c.Back
	back.Thumbnail = c.Front
	return back
}

//GetCardsFromStorageUser returns a list of card models from a storage's user deck
func GetCardsFromStorageUser(storage *storage.Storage, user string) []string {
	cardsPtr, err := storage.UserGet(user, "deck")
	if err != nil {
		return make([]string, 0)
	}
	cards := make([]string, 0)
	switch cardsPtr.(type) {
	case []interface{}:
		cardsInter := cardsPtr.([]interface{})
		for i := 0; i < len(cardsInter); i++ {
			cards = append(cards, cardsInter[i].(string))
		}
	case []string:
		cards = cardsPtr.([]string)
	}
	return cards
}

//DeleteCardsFromStorageUser deletes all instances, or a specific amount, of a card model from a storage's user deck
func DeleteCardsFromStorageUser(storage *storage.Storage, user, model string, count int) int {
	srcCards := GetCardsFromStorageUser(storage, user)
	if len(srcCards) == 0 {
		return 0
	}

	toDelete := make([]int, 0)
	for i := 0; i < len(srcCards); i++ {
		if srcCards[i] == model {
			toDelete = append(toDelete, i)
		}
		if count > 0 && len(toDelete) >= count {
			break
		}
	}

	if len(toDelete) <= 0 {
		return 0
	}

	if len(srcCards) == len(toDelete) {
		srcCards = make([]string, 0)
	} else {
		for i := 0; i < len(toDelete); i++ {
			toRemove := toDelete[i]-1
			srcCards = append(srcCards[:toRemove], srcCards[toRemove+1:]...)
		}
	}

	storage.UserSet(user, "deck", srcCards)
	return len(toDelete)
}

//GetCardsFromStorageServer returns a list of cards from a storage's server
func GetCardsFromStorageServer(storage *storage.Storage, server string) []*Card {
	cardsPtr, err := storage.ServerGet(server, "cards")
	if err != nil {
		return make([]*Card, 0)
	}
	cards := make([]*Card, 0)
	switch cardsPtr.(type) {
	case []interface{}:
		cardsInter := cardsPtr.([]interface{})
		for i := 0; i < len(cardsInter); i++ {
			cardMap := cardsInter[i].(map[string]interface{})
			cards = append(cards, &Card{
				Model: cardMap["model"].(string),
				Name: cardMap["name"].(string),
				Series: cardMap["series"].(string),
				Color: int(cardMap["color"].(float64)),
				Front: cardMap["front"].(string),
				Back: cardMap["back"].(string),
			})
		}
	case []*Card:
		cards = cardsPtr.([]*Card)
	}
	return cards
}