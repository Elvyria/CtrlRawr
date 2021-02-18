package cmd

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"

	"CtrlRawr/crawlers"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command {
	Use:   "list",
	Short: "Short Description",
	Long:  "Long Description",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := buntdb.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		urls := [][]string{{"#", "URL", "Last Update"}}

		db.CreateIndex("urls", "spider:*", buntdb.IndexJSON("host"))

		db.View(func(tx *buntdb.Tx) error {
			tx.Ascend("urls", func(key, value string) bool {
				spider := crawlers.Spider{}
				if err := json.Unmarshal([]byte(value), &spider); err != nil {
					log.Println(err)
					return true
				}

				i := strings.Index(key, ":")
				urls = append(urls, []string{key[i + 1:], spider.Url, spider.LastUpdated.Format("15:04 02/01/2006")})

				return true
			})

			return nil
		})

		fmt.Println()
		pterm.DefaultTable.WithHasHeader().WithData(urls).Render()

		return err
	},
}

func reindex(db *buntdb.DB) error {

	return nil
}
