package server

import (
	root "blocknotes_server/pkg"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(
	ns root.NoteService,

) *Server {
	s := Server{router: mux.NewRouter()}
	go ns.BatchNotesFetcher()
	// ns.NotesFixer()

	// go gts.TickerCounter(u, at, pk, wp, sch, act, attrstats)

	NewNoteRouter(ns, s.newSubrouter("/note"))

	return &s
}

func (s *Server) Start() {
	log.Println("Listening on port 8086")
	if err := http.ListenAndServe(":8086", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
