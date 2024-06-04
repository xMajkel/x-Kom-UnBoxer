package roller

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xMajkel/x-kom-unboxer/pkg/utility"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/config"
	"github.com/xMajkel/x-kom-unboxer/pkg/xkom"
)

type Roller struct {
	Account *xkom.Account
}

func (roller *Roller) Start() {
	var err error

	roller.Account, err = xkom.NewAccount(
		config.GlobalConfig.Email,
		config.GlobalConfig.Password,
		"",
	)
	if err != nil {
		panic(err)
	}

	h, m := utility.ParsePreferredRollTime(config.GlobalConfig.PreferredRollTime)

	nextRollTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), h, m, 0, 0, time.Local)

	for {
		if time.Now().After(nextRollTime) {
			nextRollTime = nextRollTime.Add(24 * time.Hour)
		}

		duration := time.Until(nextRollTime)

		log.Printf("[WAIT] Next roll: %s\n", nextRollTime.Format("2006-01-02 15:04:05"))
		time.Sleep(duration)
		log.Print("[WAIT] The time has come!\n")

		roller.RollBoxes()
		time.Sleep(1 * time.Second)

		nextRollTime = nextRollTime.Add(24 * time.Hour)
	}

}

func (roller *Roller) RollBoxes() {
	var err error
	var boxItem xkom.BoxItem

	err = roller.Account.Login()
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	boxes, err := roller.Account.GetBoxes()
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		if config.GlobalConfig.WebhookURL != "" {
			go roller.Account.SendErrorWebhook("-1", err.Error(), config.GlobalConfig.WebhookURL)
		}
		return
	}

boxLoop:
	for _, box := range boxes {
		boxIdS := fmt.Sprintf("%d", box.BoxId)
		available := true

		// Check if all requirements for box are met
		for _, requirement := range box.Requirements {
			if !requirement.IsMatched {
				available = false
				break
			}
		}

		if !available {
			log.Printf("[WARNING] requirements for %s not met\n", xkom.BoxNames[boxIdS])
			continue
		}

		// Check "NextBoxOpeningPossibleDate"
		nextOpenTime, err := time.Parse(time.RFC3339, box.NextBoxOpeningPossibleDate)
		if err == nil {
			nextOpenTime = nextOpenTime.In(time.Local)
			duration := time.Until(nextOpenTime)

			// If roll time is in more than 5 minutes skip it
			if duration > time.Minute*5 {
				log.Printf("[SKIP] %s roll not available for another: %ss\n", xkom.BoxNames[boxIdS], strings.Split(duration.String(), ".")[0])
				continue boxLoop
			}

			// Otherwise wait for it to be available
			if duration > 0 {
				log.Printf("[WAIT] %s waiting additional %ss\n", xkom.BoxNames[boxIdS], strings.Split(duration.String(), ".")[0])
				time.Sleep(duration)
			}
		}

	roll:
		for i := 1; i <= 3; i++ {
			log.Printf("[ROLL] %s roll attempt (%d/3)\n", xkom.BoxNames[boxIdS], i)

			boxItem, err = roller.Account.RollBox(boxIdS)

			if err != nil {
				log.Printf("[ERROR] %+v\n", err)

				time.Sleep(1 * time.Second)
				continue roll
			}
			break roll
		}
		if err != nil {
			log.Printf("[FAILED] %s roll: %s\n", xkom.BoxNames[boxIdS], err.Error())
			if config.GlobalConfig.WebhookURL != "" {
				go roller.Account.SendErrorWebhook(boxIdS, err.Error(), config.GlobalConfig.WebhookURL)
			}
			continue boxLoop
		}

		if config.GlobalConfig.WebhookURL != "" {
			go roller.Account.SendWebhook(boxIdS, boxItem, config.GlobalConfig.WebhookURL)
		}

		promoPercent := int((boxItem.Item.CatalogPrice - boxItem.BoxPrice) / boxItem.Item.CatalogPrice * 100)

		log.Printf("[SUCCESS] %s rolled item: \n %s\n %.2f zł -%.2f zł | -%d%% : %.2f zł\n",
			xkom.BoxNames[boxIdS],
			boxItem.Item.Name,
			boxItem.Item.CatalogPrice,
			boxItem.PromotionGain.Value,
			promoPercent,
			boxItem.BoxPrice,
		)
	}
}
