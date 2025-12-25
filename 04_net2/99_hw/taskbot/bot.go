package main

// сюда писать код

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам токен. Если захотите потыкать своего бота через телегу - используйте именно его
	BotToken = "XXX"

	// Урл, в который будет стучаться телега при получении сообщения от пользователя.
	// Может быть как айпишником личной виртуалки, так и просто выдан сервисом для деплоя
	WebhookURL = "https://525f2cb5.ngrok.io"
)

type Task struct {
	ID       int
	Title    string
	Author   *tgbotapi.User
	Assignee *tgbotapi.User
}

type TaskBot struct {
	bot    *tgbotapi.BotAPI
	tasks  []*Task
	nextID int
	mu     sync.RWMutex
	server *http.Server
}

func NewTaskBot() (*TaskBot, error) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	return &TaskBot{
		bot:    bot,
		tasks:  make([]*Task, 0),
		nextID: 1,
	}, nil
}

func (tb *TaskBot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := tb.bot.Send(msg)
	return err
}

func (tb *TaskBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	user := update.Message.From
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if !strings.HasPrefix(text, "/") {
		return
	}

	parts := strings.Fields(text)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch {
	case command == "/tasks":
		tb.handleTasks(chatID, user)
	case command == "/my":
		tb.handleMy(chatID, user)
	case command == "/owner":
		tb.handleOwner(chatID, user)
	case strings.HasPrefix(command, "/new"):
		tb.handleNew(chatID, user, text)
	case strings.HasPrefix(command, "/assign_"):
		tb.handleAssign(chatID, user, command)
	case strings.HasPrefix(command, "/unassign_"):
		tb.handleUnassign(chatID, user, command)
	case strings.HasPrefix(command, "/resolve_"):
		tb.handleResolve(chatID, user, command)
	}
}

func (tb *TaskBot) handleTasks(chatID int64, user *tgbotapi.User) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	if len(tb.tasks) == 0 {
		tb.sendMessage(chatID, "Нет задач")
		return
	}

	var result strings.Builder
	for i, task := range tb.tasks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(fmt.Sprintf("%d. %s by @%s", task.ID, task.Title, task.Author.UserName))
		if task.Assignee != nil {
			if task.Assignee.ID == user.ID {
				result.WriteString("\nassignee: я")
			} else {
				result.WriteString(fmt.Sprintf("\nassignee: @%s", task.Assignee.UserName))
			}
		}
		if task.Assignee == nil {
			result.WriteString(fmt.Sprintf("\n/assign_%d", task.ID))
		} else if task.Assignee.ID == user.ID {
			result.WriteString(fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID))
		}
	}

	tb.sendMessage(chatID, result.String())
}

func (tb *TaskBot) handleNew(chatID int64, user *tgbotapi.User, text string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		return
	}
	title := parts[1]

	tb.mu.Lock()
	task := &Task{
		ID:     tb.nextID,
		Title:  title,
		Author: user,
	}
	tb.tasks = append(tb.tasks, task)
	taskID := tb.nextID
	tb.nextID++
	tb.mu.Unlock()

	tb.sendMessage(chatID, fmt.Sprintf(`Задача "%s" создана, id=%d`, title, taskID))
}

func (tb *TaskBot) handleAssign(chatID int64, user *tgbotapi.User, command string) {
	idStr := strings.TrimPrefix(command, "/assign_")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	var task *Task
	for _, t := range tb.tasks {
		if t.ID == taskID {
			task = t
			break
		}
	}

	if task == nil {
		return
	}

	oldAssignee := task.Assignee
	task.Assignee = user

	tb.sendMessage(chatID, fmt.Sprintf(`Задача "%s" назначена на вас`, task.Title))

	if oldAssignee != nil && oldAssignee.ID != user.ID {
		tb.sendMessage(oldAssignee.ID, fmt.Sprintf(`Задача "%s" назначена на @%s`, task.Title, user.UserName))
	}

	if oldAssignee == nil && task.Author.ID != user.ID {
		tb.sendMessage(task.Author.ID, fmt.Sprintf(`Задача "%s" назначена на @%s`, task.Title, user.UserName))
	}
}

func (tb *TaskBot) handleUnassign(chatID int64, user *tgbotapi.User, command string) {
	idStr := strings.TrimPrefix(command, "/unassign_")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()
	var task *Task
	for _, t := range tb.tasks {
		if t.ID == taskID {
			task = t
			break
		}
	}

	if task == nil {
		return
	}

	if task.Assignee == nil || task.Assignee.ID != user.ID {
		tb.sendMessage(chatID, "Задача не на вас")
		return
	}

	task.Assignee = nil
	tb.sendMessage(chatID, "Принято")

	if task.Author.ID != user.ID {
		tb.sendMessage(task.Author.ID, fmt.Sprintf(`Задача "%s" осталась без исполнителя`, task.Title))
	}
}

func (tb *TaskBot) handleResolve(chatID int64, user *tgbotapi.User, command string) {
	idStr := strings.TrimPrefix(command, "/resolve_")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	tb.mu.Lock()
	var task *Task
	var taskIndex int = -1
	for i, t := range tb.tasks {
		if t.ID == taskID {
			task = t
			taskIndex = i
			break
		}
	}

	if task == nil {
		tb.mu.Unlock()
		return
	}

	if taskIndex >= 0 {
		tb.tasks = append(tb.tasks[:taskIndex], tb.tasks[taskIndex+1:]...)
	}
	tb.mu.Unlock()

	tb.sendMessage(chatID, fmt.Sprintf(`Задача "%s" выполнена`, task.Title))

	if task.Author.ID != user.ID {
		tb.sendMessage(task.Author.ID, fmt.Sprintf(`Задача "%s" выполнена @%s`, task.Title, user.UserName))
	}
}

func (tb *TaskBot) handleMy(chatID int64, user *tgbotapi.User) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	var myTasks []*Task
	for _, task := range tb.tasks {
		if task.Assignee != nil && task.Assignee.ID == user.ID {
			myTasks = append(myTasks, task)
		}
	}

	if len(myTasks) == 0 {
		tb.sendMessage(chatID, "Нет задач")
		return
	}

	var result strings.Builder
	for i, task := range myTasks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(fmt.Sprintf("%d. %s by @%s", task.ID, task.Title, task.Author.UserName))
		result.WriteString(fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID))
	}

	tb.sendMessage(chatID, result.String())
}

func (tb *TaskBot) handleOwner(chatID int64, user *tgbotapi.User) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	var ownerTasks []*Task
	for _, task := range tb.tasks {
		if task.Author.ID == user.ID {
			ownerTasks = append(ownerTasks, task)
		}
	}

	if len(ownerTasks) == 0 {
		tb.sendMessage(chatID, "Нет задач")
		return
	}

	var result strings.Builder
	for i, task := range ownerTasks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(fmt.Sprintf("%d. %s by @%s", task.ID, task.Title, task.Author.UserName))
		if task.Assignee == nil {
			result.WriteString(fmt.Sprintf("\n/assign_%d", task.ID))
		}
	}

	tb.sendMessage(chatID, result.String())
}

func (tb *TaskBot) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go tb.handleUpdate(update)
	w.WriteHeader(http.StatusOK)
}

func (tb *TaskBot) startServer(ctx context.Context) error {
	port := "8081"
	if strings.Contains(WebhookURL, ":") {
		parts := strings.Split(WebhookURL, ":")
		if len(parts) >= 3 {
			portPart := parts[2]
			portPart = strings.Split(portPart, "/")[0]
			if portPart != "" {
				port = portPart
			}
		} else if len(parts) == 2 {
			portPart := strings.TrimPrefix(parts[1], "//")
			portPart = strings.Split(portPart, "/")[0]
			if portPart != "" {
				port = portPart
			}
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", tb.webhookHandler)

	tb.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := tb.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Server error: %v", err)
	}

	return tb.server.Shutdown(context.Background())
}

func startTaskBot(ctx context.Context) error {
	bot, err := NewTaskBot()
	if err != nil {
		return err
	}

	return bot.startServer(ctx)
}

type tplParams struct {
	URL     string
	Browser string
}

const EXAMPLE = `
		Browser {{.Browser}}

		you at {{.URL}}
	`

var tmpl = template.New("123")

func handle(w http.ResponseWriter, r *http.Request) {
	params := tplParams{
		URL: r.URL.String(),
		Browser: r.UserAgent(),
	}

	tmpl.Execute(w, params)
}



func main() {
	tmpl, _ = tmpl.Parse(EXAMPLE)

	http.HandleFunc("/", handle)

	fmt.Println("server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
