package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command {
	Use:   "remove [#]",
	Short: "Short Description",
	Long:  "Long Description",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := buntdb.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Update(func(tx *buntdb.Tx) error {
			key := fmt.Sprintf("spider:%s", args[0])
			_, err := tx.Delete(key)
			return err
		})

		return nil
	},
}
