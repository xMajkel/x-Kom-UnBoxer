package xkom

import (
	"fmt"
	"strings"

	"github.com/xMajkel/godiscord"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/shared"
)

var BoxColors = map[string]string{
	"Standard":  "#adadad",
	"Unique":    "#3970ce",
	"Legendary": "#7956db",
}

var BoxNames = map[string]string{
	"-1": "Error",
	"1":  ".Box",
	"2":  "mega.Box",
	"3":  "giga.Box",
}

var BoxAvatars = map[string]string{
	"-1": "https://cdn.x-kom.pl/img/media/box/box_sygnet_standard.png",
	"1":  "https://cdn.x-kom.pl/img/media/box/box_sygnet_standard.png",
	"2":  "https://cdn.x-kom.pl/img/media/box/box_sygnet_unique.png",
	"3":  "https://cdn.x-kom.pl/img/media/box/box_sygnet_legendary.png",
}

func (acc *Account) SendWebhook(boxId string, box BoxItem, webhookUrl string) error {
	if webhookUrl == "" {
		return shared.ErrNoWebhookUrl
	}
	embed := godiscord.NewEmbed(
		box.Item.Name,
		"",
		"https://x-kom.pl/"+box.WebUrl)

	embed.SetUser("Un.Boxer", "https://assets.x-kom.pl/public-spa/xkom/404a00afb6f162d3.png")

	if box.BoxRarity.Id != "Standard" {
		embed.SetContent("@here")
	}

	embed.SetAuthor(BoxNames[boxId], "", BoxAvatars[boxId])

	embed.SetImage(box.Item.Photo.Url)
	embed.SetColor(BoxColors[box.BoxRarity.Id])

	promoPercent := int((box.Item.CatalogPrice - box.BoxPrice) / box.Item.CatalogPrice * 100)

	embed.AddField("Standard Price", fmt.Sprintf("%.2f zł", box.Item.CatalogPrice), true)
	embed.AddField("Discount", fmt.Sprintf("-%.2f zł | **-%d**%%", box.PromotionGain.Value, promoPercent), true)
	embed.AddField("Box Price", fmt.Sprintf("__**%.2f zł**__", box.BoxPrice), true)
	embed.AddField("Rating", fmt.Sprintf("%.2f/6 %s (%d)", box.ProductCommentsStatistics.AverageRating, AvgRatingToEmojis(box.ProductCommentsStatistics.AverageRating), box.ProductCommentsStatistics.TotalCount), true)
	embed.AddField("Email", fmt.Sprintf("||%s||", acc.Email), true)

	embed.SetTimestamp()
	embed.SendToWebhook(webhookUrl)

	return nil
}

func (acc *Account) SendErrorWebhook(boxId string, errorMessage string, webhookUrl string) error {
	if webhookUrl == "" {
		return shared.ErrNoWebhookUrl
	}
	embed := godiscord.NewEmbed(
		"ERROR",
		"",
		"")

	embed.SetUser("Un.Boxer", "https://assets.x-kom.pl/public-spa/xkom/404a00afb6f162d3.png")

	embed.SetContent("@here")

	embed.SetAuthor(BoxNames[boxId], "", BoxAvatars[boxId])

	embed.SetColor("#ff0000")

	embed.AddField("Details", errorMessage, true)

	embed.AddField("Email", fmt.Sprintf("||%s||", acc.Email), false)

	embed.SetTimestamp()
	embed.SendToWebhook(webhookUrl)

	return nil
}

// 🌑🌘🌗🌖🌕
var intToEmoji map[int]string = map[int]string{
	0:  ":new_moon:",             //🌑
	1:  ":new_moon:",             //🌑
	2:  ":waning_crescent_moon:", //🌘
	3:  ":waning_crescent_moon:", //🌘
	4:  ":last_quarter_moon:",    //🌗
	5:  ":last_quarter_moon:",    //🌗
	6:  ":last_quarter_moon:",    //🌗
	7:  ":waning_gibbous_moon:",  //🌖
	8:  ":waning_gibbous_moon:",  //🌖
	9:  ":waning_gibbous_moon:",  //🌖
	10: ":full_moon:",            //🌕
}

// It returns a string of 6 Discords moon emotes that represent the user rating
//
// eg. 6.0 returns: 🌕🌕🌕🌕🌕🌕
// 3.5 returns: 🌕🌕🌕🌗🌑🌑
// 2.8 returns:  🌕🌕🌖🌑🌑🌑
func AvgRatingToEmojis(avgRating float64) string {
	if avgRating >= 5.95 {
		return strings.Repeat(intToEmoji[10], 6)
	}
	ratingOutOf60 := int(avgRating * 10)
	decimal := ratingOutOf60 / 10
	fractional := ratingOutOf60 % 10
	return strings.Repeat(intToEmoji[10], decimal) + intToEmoji[fractional] + strings.Repeat(intToEmoji[0], 5-decimal)
}
