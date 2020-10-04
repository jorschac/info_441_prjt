package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"project-cahillaw/servers/widgets/widgetsrc"

	"github.com/gorilla/mux"
)

func main() {

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	dsn, dsnok := os.LookupEnv("DSN")
	if dsnok == false {
		os.Stdout.WriteString("dsn not set")
		os.Exit(1)
	}

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

	ctx := widgetsrc.WidgetContext{}
	stdb := &widgetsrc.SQLStore{}
	stdb.DB = db
	ctx.WStore = stdb

	mux := mux.NewRouter()

	mux.HandleFunc("/hello", ctx.TopTracksHandler)

	//widgets
	mux.HandleFunc("/v1/widgets/textbox", ctx.TextBoxHandlerCreate)
	mux.HandleFunc("/v1/widgets/textbox/{textboxID}", ctx.TextBoxHandler)
	mux.HandleFunc("/v1/widgets/recenttracks", ctx.RecentTracksHandlerCreate)
	mux.HandleFunc("/v1/widgets/recenttracks/{recenttracksID}", ctx.RecentTracksHandler)
	mux.HandleFunc("/v1/widgets/topmusic", ctx.TopMusicHandlerCreate)
	mux.HandleFunc("/v1/widgets/topmusic/{topmusicID}", ctx.TopMusicHandler)
	mux.HandleFunc("/v1/widgets/spotifyplaylist", ctx.SpotifyPlaylistHandlerCreate)
	mux.HandleFunc("/v1/widgets/spotifyplaylist/{spotifyplaylistID}", ctx.SpotifyPlaylistHandler)
	mux.HandleFunc("/v1/widgets/featuredmusic", ctx.FeaturedMusicHandlerCreate)
	mux.HandleFunc("/v1/widgets/featuredmusic/{featuredmusicID}", ctx.FeaturedMusicHandler)

	//social interaction
	mux.HandleFunc("/v1/widgets/{widgetID}/like", ctx.WidgetLikeHandler)
	mux.HandleFunc("/v1/widgets/{widgetID}/comment", ctx.WidgetCommentHandlerCreate)
	mux.HandleFunc("/v1/comments/{commentID}", ctx.WidgetCommentHandler)
	mux.HandleFunc("/v1/comments/{commentID}/like", ctx.WidgetCommentLikeHandler)

	//user profile page
	mux.HandleFunc("/v1/profile/{userID}", ctx.UserProfileHandler)

	log.Printf("server is lsitening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))

}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Returning a request for %s", r.URL.Path)
	w.Write([]byte("Hello World\n"))
}
