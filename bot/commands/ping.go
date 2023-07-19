package commands

import (
	"github.com/sabafly/discord-bot-template/bot/client"
	"github.com/sabafly/disgo/discord"
	"github.com/sabafly/disgo/events"
	botlib "github.com/sabafly/sabafly-lib/v2/bot"
	"github.com/sabafly/sabafly-lib/v2/handler"
)

func Ping(b *botlib.Bot[*client.Client]) handler.Command {
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

func pingCommandHandler(b *botlib.Bot[*client.Client]) handler.CommandHandler {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		message := discord.NewMessageCreateBuilder()
		message.SetContent("pong!")
		if err := event.CreateMessage(message.Build()); err != nil {
			return botlib.ReturnErr(event, err)
		}
		return nil
	}
}
