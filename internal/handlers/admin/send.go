package admin

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

type Notification struct {
	ChatID  int64
	Message *botmodels.Message
}

type NotificationBatch struct {
	Message    *botmodels.Message
	Successful int
	ErrorCount int
	ErrorList  []int64
}

func newBatch(message *botmodels.Message) *NotificationBatch {
	return &NotificationBatch{
		Message:    message,
		Successful: 0,
		ErrorCount: 0,
		ErrorList:  make([]int64, 0),
	}
}

type NotificationService struct {
	bot    *bot.Bot
	queue  chan Notification
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	wg    sync.WaitGroup
	mu    sync.Mutex
	batch *NotificationBatch
}

func NewNotificationService(bot *bot.Bot, workerCount int) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationService{
		bot:    bot,
		queue:  make(chan Notification, 1000),
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

func (ns *NotificationService) SendHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	// FIXME: too lazy to do proper parsing
	words := strings.Split(update.Message.Text, " ")
	if len(words) == 1 {
		b.SendMessage(ctx, shared.Params(update, "Не указан ID чата либо <pre>all</pre>"))
	} else {
		reply := update.Message.ReplyToMessage
		if reply == nil {
			b.SendMessage(ctx, shared.Params(update, "нужно ответить на сообщение"))
		} else {
			if words[1] == "all" {
				ns.sendMsgToAll(update, reply)

				return
			} else {
				id, _ := strconv.ParseInt(words[1], 10, 64)
				err := shared.CopyMessage(ctx, b, reply, id)
				if err != nil {
					b.SendMessage(ctx, shared.Params(update, err.Error()))
				} else {
					b.SendMessage(ctx, shared.Params(update, "Сообщение отослано"))
				}
			}
		}
	}
}

func (ns *NotificationService) sendMsgToAll(update *botmodels.Update, msg *botmodels.Message) {
	ns.Start()
	defer ns.Stop()

	chats, err := db.Chats()
	if err != nil {
		ns.bot.SendMessage(ns.ctx, shared.Params(update, "Не удалость отослать сообщение"))
		return
	}
	ns.batch = newBatch(msg)

	for _, chat := range chats {
		ns.Enqueue(chat.ID, msg)
	}

	// FIXME: context cancelling works weird and currently is not used
	<-ns.done
	ns.bot.SendMessage(context.Background(), shared.Params(update, fmt.Sprintf("%d успешно, %d ошибок.", ns.batch.Successful, ns.batch.ErrorCount)))
}

func (ns *NotificationService) Start() {
	logging.Debug("Starting notification service")
	go ns.worker()
}

func (ns *NotificationService) Stop() {
	logging.Debug("Stopping notification service")
	close(ns.queue)
	close(ns.done)
	ns.cancel()
	ns.wg.Wait()
	logging.Info("%d successful, %d errors", ns.batch.Successful, ns.batch.ErrorCount)
	logging.Info("Errors: %#v", ns.batch.ErrorList)
	ns.clearBatch()
}

func (ns *NotificationService) worker() {
	ns.wg.Add(1)
	defer ns.wg.Done()
	// NOTE: telegram has 25 messages / sec limiting, but in case
	// user sends a command during notification, ticker is a little bit slower
	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()

	for {
		select {
		case <-ns.ctx.Done():
			logging.Info("Context is cancelled, stopping notif service")
			return
		case notif, ok := <-ns.queue:
			if !ok {
				logging.Info("Queue is closed, stopping notif service")
				return
			}
			<-ticker.C
			go ns.sendNotification(&notif)
		// FIXME: kinda ugly, but i cant figure out how to do it without waiting
		case <-time.After(time.Second):
			ns.done <- struct{}{}
		}
	}
}

func (ns *NotificationService) Enqueue(chatID int64, message *botmodels.Message) {
	select {
	case ns.queue <- Notification{ChatID: chatID, Message: message}:
	case <-time.After(5 * time.Second):
		logging.Warning("Notification queue full, dropping message for %d", chatID)
	}
}

func (ns *NotificationService) sendNotification(notif *Notification) {
	ns.wg.Add(1)
	defer ns.wg.Done()
	logging.Trace("Sending message to %d", notif.ChatID)
	err := shared.CopyMessage(ns.ctx, ns.bot, notif.Message, notif.ChatID)
	if err != nil {
		ns.recordFailure(notif.ChatID)
	} else {
		ns.recordSuccess()
	}
}

func (ns *NotificationService) recordSuccess() {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.batch.Successful++
}

func (ns *NotificationService) recordFailure(chatID int64) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.batch.ErrorCount++
	ns.batch.ErrorList = append(ns.batch.ErrorList, chatID)
}

func (ns *NotificationService) clearBatch() {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.batch = nil
}
