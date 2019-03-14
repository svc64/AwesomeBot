package main
import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)
var token string // should be provided at build time
func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println("ERRRRRR")
		os.Exit(21)
	}
	b.Handle("/hello", func(m *tb.Message) {
		_, err := b.Send(m.Chat, "HI THERE")
		handleError(err, *m)
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			_, err := b.Send(m.Chat, "ERRRRR")
			handleError(err, *m)
		}
		err := b.Ban(m.Chat, user)
		if err != nil {
			_, err := b.Send(m.Chat, "ERRRR")
			handleError(err, *m)
		}
	})
	// handle bans
	b.Handle("/ban", func(m *tb.Message) {
		replied := m.ReplyTo
		sender, _ := b.ChatMemberOf(m.Chat, m.Sender)
		user, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if user.CanRestrictMembers || // check if the sender is an admin or the group creator.
			tb.Creator == sender.Role {
			user, _ := b.ChatMemberOf(m.Chat, replied.Sender)
			if err != nil {
				_, err := b.Send(m.Chat, "ERRRRR")
				handleError(err, *m)
			}
			err := b.Ban(m.Chat, user)
			if err != nil {
				_, err := b.Send(m.Chat, "ERRRR")
				handleError(err, *m)
			}
		} else { // runs when the sender doesn't have permission to ban users
			_, err := b.Reply(m, "You don't have permission to ban users!")
			handleError(err, *m)
		}
	})
	b.Start()
}
