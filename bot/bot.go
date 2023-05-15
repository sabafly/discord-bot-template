package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/dislog"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	"github.com/sabafly/discord-bot-template/bot/commands"
	"github.com/sabafly/discord-bot-template/bot/db"
	botlib "github.com/sabafly/sabafly-lib/bot"
	"github.com/sabafly/sabafly-lib/handler"
	"github.com/sabafly/sabafly-lib/translate"
	"github.com/sirupsen/logrus"
)

var (
	version string = "dev"
)

func init() {
	botlib.BotName = "template-bot"
	botlib.Color = 0x252525
}

func Run(config_path, lang_path string) {
	if _, err := translate.LoadTranslations(lang_path); err != nil {
		panic(err)
	}
	cfg, err := botlib.LoadConfig(config_path)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.ReportCaller = cfg.DevMode
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
	logger.SetOutput(colorable.NewColorableStdout())
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lvl)
	if cfg.Dislog.Enabled {
		dlog, err := dislog.New(
			dislog.WithLogLevels(dislog.TraceLevelAndAbove...),
			dislog.WithWebhookIDToken(cfg.Dislog.WebhookID, cfg.Dislog.WebhookToken),
		)
		if err != nil {
			logger.Fatal("error initializing dislog: ", err)
		}
		defer dlog.Close(context.TODO())
	}
	logger.Infof("Starting bot version: %s", version)
	logger.Infof("Syncing commands? %t", cfg.ShouldSyncCommands)

	b := botlib.New[*db.DB](logger, version, *cfg)

	b.DB, err = db.SetupDatabase(cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	defer func() { _ = b.DB.Close() }()

	b.Handler.AddCommands(
		commands.Ping(b),
	)

	b.Handler.AddComponents()

	b.Handler.AddModals()

	b.Handler.AddMemberJoins(
		handler.MemberJoin{
			UUID: uuid.New(),
			Handler: func(event *events.GuildMemberJoin) error {
				b.OnGuildMemberJoin(event)
				return nil
			},
		},
	)

	b.Handler.AddMemberLeaves(
		handler.MemberLeave{
			UUID: uuid.New(),
			Handler: func(event *events.GuildMemberLeave) error {
				b.OnGuildMemberLeave(event)
				return nil
			},
		},
	)

	b.Handler.AddReady(func(r *events.Ready) {
		r.Client().Logger().Infof("Ready! %s", r.User.Tag())
	})

	b.SetupBot(bot.NewListenerFunc(b.Handler.OnEvent))
	b.Client.EventManager().AddEventListeners(&events.ListenerAdapter{
		OnGuildJoin:  b.OnGuildJoin,
		OnGuildLeave: b.OnGuildLeave,
	})

	if cfg.ShouldSyncCommands {
		var guilds []snowflake.ID
		if cfg.DevOnly {
			guilds = b.Config.DevGuildIDs
		}
		b.Handler.SyncCommands(b.Client, guilds...)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := b.Client.OpenGateway(ctx); err != nil {
		b.Logger.Fatalf("failed to connect gateway: %s", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer b.Client.Close(ctx)

	b.Logger.Info("Bot is running, Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)
	<-s
	b.Logger.Info("Shutting down....")
}
