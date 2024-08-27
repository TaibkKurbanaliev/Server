package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	dataBase *DataBase
	Router   *mux.Router
}

func NewServer(dbConnectionString string, currencyToken string) *Server {
	var server *Server = new(Server)
	server.dataBase = NewDataBase("postgres", "", "", "disable")
	server.Router = mux.NewRouter()
	server.Router.StrictSlash(true)
	server.initServerHandelFunctions()
	// if err != nil {
	// 	panic(err)
	// }
	return server
}

func (server Server) initServerHandelFunctions() {
	server.Router.HandleFunc("/user/", server.createHandler).Methods("POST")
	server.Router.HandleFunc("/user/", server.getAllHandler).Methods("GET")
	server.Router.HandleFunc("/user/{id:[0-9]+}/", server.getHandler).Methods("GET")
	server.Router.HandleFunc("/user/{id:[0-9]+}/", server.deleteHandler).Methods("DELETE")
	server.Router.HandleFunc("/user/", server.updateHandler).Methods("UPDATE")
	server.Router.HandleFunc("/user/wallPaper/{id:[0-9]+}", server.createWallPaper).Methods("POST")
	server.Router.HandleFunc("/user/wallPaper/{id:[0-9]+}", server.getWallPaperById).Methods("GET")
}

func (server *Server) createHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server.dataBase.Add(user, "user")
}

func (server *Server) getAllHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Print(req.URL.Path)
}

func (server *Server) getHandler(w http.ResponseWriter, req *http.Request) {
	// implement logic to handle request
}

func (server *Server) deleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server.dataBase.Delete(id, "user")
}

func (server *Server) updateHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server.dataBase.Update(user.ID, "user", user)
}

func (server *Server) createWallPaper(w http.ResponseWriter, req *http.Request) {
	var jsonImage JsonImage
	err := json.NewDecoder(req.Body).Decode(&jsonImage)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var WallPaper WallPaper
	err = WallPaper.Init(jsonImage)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = server.dataBase.Add(WallPaper, "wallPaper")

	if err != nil {
		log.Println(err.Error())
	}
}

func (server *Server) getWallPaperById(w http.ResponseWriter, req *http.Request) {
	//  implement logic to handle request
}
