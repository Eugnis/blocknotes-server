package server

import (
	root "blocknotes_server/pkg"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type noteRouter struct {
	noteService root.NoteService
}

func NewNoteRouter(ns root.NoteService, router *mux.Router) *mux.Router {
	noteRouter := noteRouter{ns}

	// router.HandleFunc("/", noteRouter.createNoteHandler).Methods("PUT")
	router.HandleFunc("/view/{address}", noteRouter.getNoteHandler).Methods("GET")
	router.HandleFunc("/list", noteRouter.listNotesHandler).Methods("POST")
	// router.HandleFunc("/update_admin", ValidateAdminMiddleware(attributeRouter.updateAttributeHandlerAdmin, a)).Methods("POST")
	// router.HandleFunc("/remove/{id}", ValidateAdminMiddleware(attributeRouter.removeAttributeHandler, a)).Methods("DELETE")
	return router
}

func (sr *noteRouter) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	note, err := decodeNote(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = sr.noteService.Create(&note)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, map[string]string{"result": "true"})
}

func (sr *noteRouter) getNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// log.Println(vars)
	address := vars["address"]

	note, err := sr.noteService.GetByNoteAddress(address)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, note)
}

func (sr *noteRouter) listNotesHandler(w http.ResponseWriter, r *http.Request) {
	nsr, err := decodeNoteSearchRequest(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// if nsr.Count == 0 || nsr.From == 0 {
	// 	Error(w, http.StatusBadGateway, "Please provide both /from/count")
	// 	return
	// }

	if nsr.Count > 5000 {
		Error(w, http.StatusBadGateway, "You can get max 5000 records")
		return
	}

	attrs, err := sr.noteService.ListNotes(nsr)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, attrs)
}

// func (sr *noteRouter) removeNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	log.Println(vars)
// 	id := vars["id"]

// 	err := sr.noteService.Remove(id)
// 	if err != nil {
// 		Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	newAttr := root.Attribute{ID: bson.ObjectIdHex(id)}
// 	err = sr.userService.UpdateAttributes(&newAttr, true)
// 	if err != nil {
// 		Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	Json(w, http.StatusOK, "true")
// }

// func (sr *noteRouter) updateAttributeHandlerAdmin(w http.ResponseWriter, r *http.Request) {
// 	attr, err := decodeAttribute(r)
// 	if err != nil {
// 		log.Printf(err.Error())
// 		Error(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	// log.Printf("username " + user.ID)
// 	// log.Printf("workplaceid " + user.WorkplaceID)
// 	err = sr.noteService.Update(&attr)
// 	if err != nil {
// 		log.Println(err.Error())
// 		Error(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	// ur.getMeHandler(w, r)

// 	Json(w, http.StatusOK, true)
// }

func decodeNote(r *http.Request) (root.Note, error) {
	var s root.Note
	if r.Body == nil {
		return s, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&s)
	return s, err
}

func decodeNoteSearchRequest(r *http.Request) (root.NoteSearch, error) {
	var s root.NoteSearch
	if r.Body == nil {
		return s, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&s)
	return s, err
}
