package cmd

import (
	"encoding/json"
	"log"

	"CtrlRawr/crawlers"

	"github.com/go-toast/toast"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var fHost   string
var fIndex  uint
var fForced bool

func init() {
	updateCmd.Flags().StringVar(&fHost,  "host", "", "-h animixplay.to")
	updateCmd.Flags().UintVarP(&fIndex,  "index", "i", 0, "-i 1")
	updateCmd.Flags().BoolVarP(&fForced, "forced", "f", false, "")

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command {
	Use:   "update",
	Short: "Short Description",
	Long:  "Long Description",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := buntdb.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Update(func(tx *buntdb.Tx) error {
			var updated []string

			err = tx.AscendKeys("spider:*", func(key, value string) bool {
				spider := crawlers.Spider{}
				err := json.Unmarshal([]byte(value), &spider)
				if err != nil {
					log.Println(err.Error())
					return true
				}

				hash := spider.Hash

				info, err := spider.Crawl()
				if err != nil {
					log.Println(err.Error())
					return true
				}

				if hash != spider.Hash {
					notify(info)
				}

				b, err := json.Marshal(spider)
				if err != nil {
					log.Println(err.Error())
					return true
				}

				updated = append(updated, key, string(b))

				return true
			})

			for i := 0; i < len(updated); i += 2 {
				tx.Set(updated[i], updated[i + 1], nil)
			}

			return err
		})

		return err
	},
}

func notify(info *crawlers.Info) error {
	notification := toast.Notification {
		AppID: "CtrlRawr",
		Title: info.Title,
		Message: info.Message,
		// ActivationArguments: url,
	}

	return notification.Push()
}
