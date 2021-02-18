package cmd

import (
	"fmt"
	"encoding/json"
	"strconv"

	"CtrlRawr/crawlers"

	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command {
	Use:   "add [url]",
	Short: "Short Description",
	Long:  "Long Description",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := buntdb.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		ref := args[0]
		if contains(db, ref) {
		}

		spider, err := crawlers.Create(ref)
		if err != nil {
			return err
		}

		return addSpider(db, spider)
	},
}

func addSpider(db *buntdb.DB, spider *crawlers.Spider) error {
	err := db.Update(func(tx *buntdb.Tx) error {
		value, err := json.Marshal(spider)
		if err != nil {
			return err
		}

		index := 0

		sIndex, err := tx.Get("index")
		if err == nil {
			index, err = strconv.Atoi(sIndex)
		}

		err = nil
		for err == nil {
			index += 1
			sIndex = fmt.Sprint(index)

			_, err = tx.Get(sIndex)
		}

		tx.Set("index", sIndex, nil)
		key := fmt.Sprintf("spider:%s", sIndex)

		tx.Set(key, string(value), nil)

		return nil
	})

	return err
}

func contains(db *buntdb.DB, ref string) bool {
	result := false

	db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("spiders", func(key, value string) bool {
			if value == ref {
				result = true

				return false
			}

			return true
		})

		return nil
	})

	return result
}
