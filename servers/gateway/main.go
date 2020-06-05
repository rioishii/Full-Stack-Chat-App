package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/models/users"
	"github.com/streadway/amqp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/handlers"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/sessions"
)

type director func(r *http.Request)

func customDirector(targets []string, ctx *handlers.HandlerCtx) director {
	var counter int32
	counter = 0

	return func(r *http.Request) {
		r.Header.Del("X-User")
		sessionState := &handlers.SessionState{}
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
		if sessionState.User != nil && err == nil {
			json, _ := json.Marshal(sessionState.User)
			log.Println(string(json))
			r.Header.Add("X-User", string(json))
		}
		i32 := int32(len(targets))
		targ := targets[counter%i32]
		atomic.AddInt32(&counter, 1)
		r.Host = targ
		r.URL.Host = targ
		r.URL.Scheme = "http"
	}
}

// main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	summaryAddr := os.Getenv("SUMMARY")
	if len(summaryAddr) == 0 {
		summaryAddr = "summary:4000"
	}
	summaryURL := []string{summaryAddr}

	messagingAddr := strings.Split(os.Getenv("CHAT"), ",")
	messageURL := make([]string, 0)
	for _, addr := range messagingAddr {
		messageURL = append(messageURL, addr)
	}

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Print("Path to TLS public certificate and/or path to the associated private key is not set")
		os.Exit(1)
	}

	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		sessionKey = "signing key"
	}

	redisaddr := os.Getenv("REDISADDR")
	if len(redisaddr) == 0 {
		redisaddr = "redisServer:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	sessionStore := sessions.NewRedisStore(client, time.Hour)

	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		dsn = fmt.Sprintf("root:%s@tcp(mysqlServer:3306)/userDB", os.Getenv("MYSQL_ROOT_PASSWORD"))
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	sqlStore := users.NewSQLStore(db)

	trie, err := sqlStore.GetAllUsers()
	if err != nil {
		fmt.Printf("error creating trie: %v", err)
	}

	rabbitAddr := os.Getenv("RABBITADDR")
	if len(rabbitAddr) == 0 {
		rabbitAddr = "amqp://guest:guest@rabbitmq:5672/"
	}
	conn, err := amqp.Dial(rabbitAddr)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening a channel: %s", err)
	}
	defer ch.Close()

	q, _ := ch.QueueDeclare(
		"test", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	notifier := handlers.NewNotifier()

	ctx := handlers.NewHandlerContext(sessionKey, sessionStore, sqlStore, trie, notifier)

	go ctx.Notifier.NotifyWebSockets(msgs)

	mux := mux.NewRouter()

	summaryProxy := &httputil.ReverseProxy{Director: customDirector(summaryURL, ctx)}
	messagingProxy := &httputil.ReverseProxy{Director: customDirector(messageURL, ctx)}

	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/channels", messagingProxy)
	mux.Handle("/v1/channels/{channelID}", messagingProxy)
	mux.Handle("/v1/channels/{channelID}/members", messagingProxy)
	mux.Handle("/v1/messages/{messageID}", messagingProxy)
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/{id}", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/{id}", ctx.SpecificSessionHandler)
	mux.HandleFunc("/v1/ws", ctx.WebSocketConnectionHandler)
	wrappedMux := &handlers.CORS{Handler: mux}

	log.Printf("server listening at: %s", addr)
	http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux)
}
