package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	dataBase             *DataBase
	Router               *mux.Router
	WallPaperStoragePath string
}

func NewServer(configurationFilePath string) *Server {
	var server *Server = new(Server)
	configuration, err := readConfigurationFile(configurationFilePath)

	if err != nil {
		log.Panic(err)
		return nil
	}

	server.dataBase = NewDataBase(configuration["dbConnectionString"].(string))
	server.WallPaperStoragePath = configuration["storagePath"].(string)
	server.Router = mux.NewRouter()
	server.Router.StrictSlash(true)
	server.initServerHandelFunctions()

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

	wallPaper, err := NewWallPaper(jsonImage, server.WallPaperStoragePath)

	if err != nil {
		log.Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := server.WallPaperStoragePath + jsonImage.FileName
	err = writeFileToStorage(jsonImage.Image, path)

	if err != nil {
		log.Panic(err.Error())
		return
	}

	err = server.dataBase.Add(*wallPaper, "wallPaper")

	if err != nil {
		log.Panic(err.Error())
		return
	}
}

func (server *Server) getWallPaperById(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])

	if err != nil {
		log.Panic(err)
		return
	}

	var wallPaper WallPaper
	err = server.dataBase.Select(id, "wallPaper", &wallPaper)

	if err != nil {
		log.Panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := readImageFromStorage(wallPaper.ImagePath, wallPaper.Name)

	if err != nil {
		log.Panic(err)
		return
	}

	if err != nil {
		log.Panic(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData.Bytes())
}

func readConfigurationFile(filePath string) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var jsonData []byte
	jsonData, err = io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &jsonMap)

	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func writeFileToStorage(data []byte, path string) error {
	img, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return err
	}

	out, _ := os.Create("." + path)
	defer out.Close()

	err = png.Encode(out, img)

	if err != nil {
		return err
	}

	return nil
}

func readImageFromStorage(filePath string, fileName string) (*bytes.Buffer, error) {
	file, err := os.Open("." + filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	byt := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(byt)

	if err != nil {
		return nil, err
	}

	var jsonImage = JsonImage{
		FileName: fileName,
		Image:    byt,
	}
	jsonBytes, err := json.Marshal(jsonImage)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBytes), err
}
