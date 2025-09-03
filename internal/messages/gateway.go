package messages

import (
	"context"
	"sync"
	"time"

	"github.com/go-telegram/bot" 
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
)

type Notification struct {
	ChatID  int64
	Message *botmodels.Message
}

type MessageGateway struct {
	bot    *bot.Bot
	queue  chan Notification
	ctx    context.Context
	cancel context.CancelFunc

	wg    sync.WaitGroup
	mu    sync.Mutex
	done  chan struct{}
}

func NewMessageGateway(bot *bot.Bot) *MessageGateway {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageGateway{
		bot:    bot,
		queue:  make(chan Notification, 100),
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

func (mg *MessageGateway) Enqueue(chatID int64, msg *botmodels.Message) {
	select {
	case mg.queue <- Notification{chatID, msg}:
	case <-time.After(time.Second * 5):
		logging.Warning("Notification queue full for 5s, dropping message for %d", chatID)
	}
}

func (mg *MessageGateway) worker() {
	ticker := time.NewTicker(time.Millisecond * 40)
	select {
	case <-mg.ctx.Done():
		logging.Info("Context is cancelled, stopping notif service")
		return
	case msg, ok := <-mg.queue:
		if !ok {
			logging.Info("Queue is closed, stopping gateway service")
			return
		}
		<-ticker.C
		mg.bot.SendMessage(mg.ctx, &bot.SendMessageParams{
			ChatID: msg.ChatID,
			Text:   msg.Message.Text,
		})
	}
}

func (mg *MessageGateway) Start() {
	logging.Debug("Starting notification service")
	go mg.worker()
}

func (mg *MessageGateway) Stop() {
	logging.Debug("Stopping notification service")
	close(mg.queue)
	close(mg.done)
	mg.cancel()
	mg.wg.Wait()
}
