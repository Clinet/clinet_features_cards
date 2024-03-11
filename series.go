package cards

//SeriesList is used to (de)serialize the entire list of cards
type SeriesList struct {
	Series map[string]*Series `json:"series"` //Map of card series, where key is the series number
}

type Series struct {
	Released bool    `json:"released"` //Whether or not this series has been released
	Cards    []*Card `json:"cards"`    //List of cards in this series
}

/*
{
	"series": {
		"1": {
			"released": true,
			"cards": [
				{
					"title": "Rick Grimes",
					"color": "purple",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559302047699074/1.png?ex=6600d424&is=65ee5f24&hm=fc5cdb204fdff2e5f60d2afc2599fac4557e4927770fa4fba4a231fcc052a6b1",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559302408536115/1b.png?ex=6600d424&is=65ee5f24&hm=906a8a29cabbdec925b4680543ab383af19628d020de3338896390405fb089a6"
				},
				{
					"title": "The Prophecy",
					"color": "blue",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				},
				{
					"title": "King Does",
					"color": "green",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				},
				{
					"title": "SGT. Greasemixer",
					"color": "red",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				},
				{
					"title": "Kanin The Hare",
					"color": "red",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
						"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				},
				{
					"title": "Cow Booboo",
					"color": "blue",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				},
				{
					"title": "Ayyjay121",
					"color": "red",
					"front": "https://media.discordapp.net/attachments/676312586123214858/1216559319693262848/2.png?ex=6600d428&is=65ee5f28&hm=911e7bc92b451c005448fb20e023869532242eac799b28709ec28686a610741f",
					"back": "https://media.discordapp.net/attachments/676312586123214858/1216559320108367972/2b.png?ex=6600d428&is=65ee5f28&hm=2b1edec881d5351150da12eb2617c7b962971b68eb339c55d32574b2051c9767",
				}
			]
		}
	}
}
*/