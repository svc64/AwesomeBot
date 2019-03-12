package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"sync"
	"time"
)
func handleAdmin(token string) {
	b, _ := tb.NewBot(tb.Settings{
		Token:  token,
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, err := b.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			b.Send(m.Chat, "ERRRRR")
		}
		err = b.Ban(m.Chat, user)
		if err != nil {
			b.Send(m.Chat, "ERRRR")
		}
	})

	b.Handle("/ban", func(m *tb.Message) {
		replied := m.ReplyTo
		sender, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if tb.Administrator == sender.Role ||
			tb.Creator == sender.Role {
			user, err := b.ChatMemberOf(m.Chat, replied.Sender)
			if err != nil {
				b.Send(m.Chat, "ERRRRR")
			}
			err = b.Ban(m.Chat, user)
			if err != nil {
				b.Send(m.Chat, "ERRRR")
			}
		}
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Start()
	}()
}
