package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"wikistream/internal/stream"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
	lang    string
}

func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		session: session,
		lang:    "en",
	}

	session.AddHandler(bot.messageHandler)
	return bot, nil
}

func (b *Bot) Run() error {
	if err := b.session.Open(); err != nil {
		return err
	}

	log.Println("Bot is running. Press CTRL+C to exit.")

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	return b.session.Close()
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "!recent":
		go b.handleRecent(s, m)
	case "!setLang":
		if len(args) > 1 {
			b.handleSetLang(s, m, args[1])
		} else {
			s.ChannelMessageSend(m.ChannelID, "Usage: `!setLang [language_code]`")
		}
	}
}

func (b *Bot) handleRecent(s *discordgo.Session, m *discordgo.MessageCreate) {
	changes, err := stream.GetRecentChanges(b.lang, 5)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching changes")
		return
	}

	if len(changes) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No changes found for `%s`", b.lang))
		return
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("**Last %d Wikipedia Changes for `%s`:**\n", len(changes), b.lang))

	for _, change := range changes {
		message.WriteString(fmt.Sprintf("\n**Page:** %s\n", change.Title))
		message.WriteString(fmt.Sprintf("**URL:** https://%s/wiki/%s\n",
			change.ServerName,
			strings.ReplaceAll(change.Title, " ", "_")))
		message.WriteString(fmt.Sprintf("**User:** %s\n", change.User))
		message.WriteString(fmt.Sprintf("**Timestamp:** %s\n",
			stream.FormatTimestamp(change.Timestamp)))
		message.WriteString("----------------------------------------\n")
	}

	s.ChannelMessageSend(m.ChannelID, message.String())
}

func (b *Bot) handleSetLang(s *discordgo.Session, m *discordgo.MessageCreate, newLang string) {
	b.lang = newLang
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Language set to `%s`", b.lang))
}
