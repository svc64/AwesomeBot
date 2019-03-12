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
		b.Send(m.Sender, "HI THERE")
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			b.Send(m.Chat, "ERRRRR")
		}
		err := b.Ban(m.Chat, user)
		if err != nil {
			b.Send(m.Chat, "ERRRR")
		}
	})

	b.Handle("/ban", func(m *tb.Message) {
		replied := m.ReplyTo
		user, _ := b.ChatMemberOf(m.Chat, replied.OriginalSender)
		if err != nil {
			b.Send(m.Chat, "ERRRRR")
		}
		err := b.Ban(m.Chat, user)
		if err != nil {
			b.Send(m.Chat, "ERRRR")
		}
	})
	b.Start()
}