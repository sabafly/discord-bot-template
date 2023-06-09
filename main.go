package main

import (
	"os"

	"github.com/sabafly/discord-bot-template/bot"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gobot",
	Short: "とても便利でおいしいディスコードボット",
	// TODO: 書く
	Long: `後で書く`,
}

func main() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	botCmd.Flags().String("config", "config.json", "config file of bot")
	botCmd.Flag("config").Shorthand = "c"
	botCmd.Flags().String("lang", "lang", "lang file path")
	botCmd.Flag("lang").Shorthand = "l"
	root.AddCommand(botCmd)
}

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "botを起動する",
	Long: `botの説明
	後で書く`,
	ValidArgs: []string{
		"config",
	},
	Run: func(cmd *cobra.Command, args []string) {
		bot.Run(cmd.Flag("config").Value.String(), cmd.Flag("lang").Value.String())
	},
}
