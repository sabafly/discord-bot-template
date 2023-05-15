package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sabafly/discord-bot-template/bot/db"
	botlib "github.com/sabafly/sabafly-lib/bot"
	"github.com/sabafly/sabafly-lib/handler"
)

func Ping(b *botlib.Bot[*db.DB]) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "ping",
			Description: "pong!",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": pingCommandHandler(b),
		},
	}
}

func pingCommandHandler(b *botlib.Bot[*db.DB]) handler.CommandHandler {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		message := discord.NewMessageCreateBuilder()
		message.SetContent("pong!")
		if err := event.CreateMessage(message.Build()); err != nil {
			return botlib.ReturnErr(event, err)
		}
		return nil
	}
}
