package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"project-cahillaw/servers/gateway/handlers"
	"project-cahillaw/servers/gateway/models/users"
	"project-cahillaw/servers/gateway/sessions"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlsKeyPath, keyok := os.LookupEnv("TLSKEY")
	tlsCertPath, certok := os.LookupEnv("TLSCERT")

	//will switch for widget/other ms we end up having
	widgetsaddr, wdgok := os.LookupEnv("WIDGETSADDR")
	if wdgok == false {
		os.Stdout.WriteString("widgetsaddr not set")
		os.Exit(1)
	}
    
	wdgAddrArray := strings.Split(widgetsaddr, ",")
	wdgAddrArrayLength := len(wdgAddrArray)
	wdgAddrToUseNum := 0
	wdgAddrToUse := wdgAddrArray[wdgAddrToUseNum]

	if keyok == false || certok == false {
		os.Stdout.WriteString("key/certificate not set")
		os.Exit(1)
	}

	//session key is int for right now
	//claim the address for redis store like rServe:6379, and store it as an Env Varibale
	sessionKey, sessok := os.LookupEnv("SESSIONKEY")
	redisaddr, redisok := os.LookupEnv("REDISADDR")
	dsn, dsnok := os.LookupEnv("DSN")

	if sessok == false {
		os.Stdout.WriteString("signing key not set")
		os.Exit(1)
	}

	if redisok == false {
		os.Stdout.WriteString("redisaddr not set")
		os.Exit(1)
	}

	if dsnok == false {
		os.Stdout.WriteString("dsn not set")
		os.Exit(1)
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		os.Stdout.WriteString(err.Error())
		os.Exit(1)
	}
	fmt.Println(pong, err)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		os.Stdout.WriteString("error opening database")
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	} else {
		fmt.Printf("successfully connected!\n")
	}

	ctx := &handlers.HandlerContext{}
	ctx.Key = sessionKey
	ctx.SessionsStore = sessions.NewRedisStore(client, time.Hour)
	stdb := &users.SQLStore{}
	stdb.DB = db
	ctx.UserStore = stdb

	// 验证并返回用户信息,存储在request的header里 
	widgetsDirector := func(r *http.Request) {
		serverName := wdgAddrToUse
		if wdgAddrToUseNum < wdgAddrArrayLength-1 {
			wdgAddrToUseNum = wdgAddrToUseNum + 1
			wdgAddrToUse = wdgAddrArray[wdgAddrToUseNum]
		} else {
			wdgAddrToUseNum := 0
			wdgAddrToUse = wdgAddrArray[wdgAddrToUseNum]
		}

		r.Header.Del("X-User")
		ss := &handlers.SessionState{}
		sid, _ := sessions.GetSessionID(r, sessionKey)
		errGetSession := ctx.SessionsStore.Get(sid, &ss)
		if errGetSession == nil {
			encoded, _ := json.Marshal(ss.User)
			r.Header.Set("X-User", string(encoded))
		}
		r.Host = serverName
		r.URL.Host = serverName
		r.URL.Scheme = "http"
	}
	

	mux := mux.NewRouter()
	// so waht is httputil and what is ReverseProxy{}
	widgetsProxy := &httputil.ReverseProxy{Director: widgetsDirector}
	mux.Handle("/v1/widgets/textbox", widgetsProxy)
	mux.Handle("/v1/widgets/textbox/{textboxID}", widgetsProxy)
	mux.Handle("/v1/widgets/recenttracks", widgetsProxy)
	mux.Handle("/v1/widgets/recenttracks/{recenttracksID}", widgetsProxy)
	mux.Handle("/v1/widgets/textbox", widgetsProxy)
	mux.Handle("/v1/widgets/textbox/{textboxID}", widgetsProxy)
	mux.Handle("/v1/widgets/spotifyplaylist", widgetsProxy)
	mux.Handle("/v1/widgets/spotifyplaylist/{spotifyplaylistID}", widgetsProxy)
	mux.Handle("/v1/widgets/featuredmusic", widgetsProxy)
	mux.Handle("/v1/widgets/featuredmusic/{featuredmusicID}", widgetsProxy)

	mux.Handle("/v1/widgets/{widgetID}/like", widgetsProxy)
	mux.Handle("/v1/widgets/{widgetID}/comment", widgetsProxy)
	mux.Handle("/v1/comments/{commentID}", widgetsProxy)
	mux.Handle("/v1/comments/{commentID}/like", widgetsProxy)

	mux.Handle("/v1/profile/{userID}", widgetsProxy)

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/{userID}", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/{sessionID}", ctx.SpecificSessionHandler)

	wrappedMux := handlers.AddFiveResponseHeaders(mux, "Access-Control-Allow-Origin", "*", "Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE", "Access-Control-Allow-Headers", "Content-Type, Authorization", "Access-Control-Expose-Headers", "Authorization", "Access-Control-Max-Age", "600")

	log.Printf("server is lsitening at %s", addr)

	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
