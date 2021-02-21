package server

import (
	//"encoding/json"
	//"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scottapp/misuchiru-server/handlers"
	"github.com/scottapp/misuchiru-server/utils"
	"strconv"

	//"github.com/kuso/japanese-word-extractor/extractor"
	//"github.com/lithammer/shortuuid/v3"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	"time"
)

type Server struct {
	Router       *gin.Engine
	HttpServer   *http.Server
	TokenHandler *handlers.TokenHandlers
	AlbumDB      *handlers.AlbumDB
}

type ResponseSong struct {
	ID   int
	Name string
}

type ResponseAlbum struct {
	ID    int
	Name  string
	Songs []*ResponseSong
}

func NewServer() *Server {
	server := Server{}
	server.TokenHandler = &handlers.TokenHandlers{}

	server.AlbumDB = handlers.ReadAllAlbums()

	server.SetupRouter()

	server.HttpServer = &http.Server{
		Addr:           ":8081",
		Handler:        server.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &server
}

func (server *Server) GracefulShutdown(timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nshutdown with timeout: %s\n", timeout)

	if err := server.HttpServer.Shutdown(ctx); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
}

func (server *Server) SetupRouter() {
	server.Router = gin.Default()

	server.Router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := server.Router.Group("/v1")
	{
		v1.GET("/hello", server.Hello)
		v1.GET("/all_albums", server.GetAllAlbums)
		v1.GET("/albums", server.GetAlbums)
		v1.GET("/album/:album/songs", server.GetSongs)
		v1.GET("/song/:album/:song", server.GetSong)
	}
}

func (server *Server) Hello(c *gin.Context) {
	result := gin.H{"hello": "world"}
	c.JSON(http.StatusOK, result)
}

func (server *Server) GetAllAlbums(c *gin.Context) {
	albumIds := []int{19, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	outAlbums := []*ResponseAlbum{}
	for _, albumId := range albumIds {
		album, ok := server.AlbumDB.AlbumByID[albumId]
		if !ok {
			//c.JSON(404, gin.H{"error": "not found"})
			continue
		}
		outAlbum := ResponseAlbum{}
		songIds := utils.GetSortedValueInts(album.URLIDTable)
		songs := []*ResponseSong{}
		for i := 0; i < len(songIds); i++ {
			songId := songIds[i]
			outSong := ResponseSong{}
			outSong.ID = songId
			outSong.Name = album.SongByID[songId].Name
			songs = append(songs, &outSong)
		}
		outAlbum.ID = album.ID
		outAlbum.Name = album.Name
		outAlbum.Songs = songs
		outAlbums = append(outAlbums, &outAlbum)
	}
	c.JSON(http.StatusOK, gin.H{"albums": outAlbums})
}

func (server *Server) GetSongs(c *gin.Context) {
	albumId, err := strconv.Atoi(c.Param("album"))
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "not found"})
	}
	result := []*ResponseSong{}
	album, ok := server.AlbumDB.AlbumByID[albumId]
	if !ok {
		log.Println("cannot find album id " + strconv.Itoa(albumId))
		c.JSON(404, gin.H{"error": "not found"})
	}
	for id, song := range album.SongByID {
		obj := ResponseSong{}
		obj.ID = id
		obj.Name = song.Name
		result = append(result, &obj)
	}
	c.JSON(http.StatusOK, gin.H{"album": album.Name, "songs": result})
}

func (server *Server) GetAlbums(c *gin.Context) {
	result := []*ResponseAlbum{}
	for _, album := range server.AlbumDB.Albums {
		obj := ResponseAlbum{}
		obj.ID = album.ID
		obj.Name = album.Name
		result = append(result, &obj)
	}
	c.JSON(http.StatusOK, result)
}

func (server *Server) GetSong(c *gin.Context) {
	albumId, err := strconv.Atoi(c.Param("album"))
	if err != nil {
		log.Fatal(err)
	}
	songId, err := strconv.Atoi(c.Param("song"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(albumId, songId)

	album, ok := server.AlbumDB.AlbumByID[albumId]
	if !ok {
		log.Println(strconv.Itoa(albumId) + ", album not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	song, ok := album.SongByID[songId]
	if !ok {
		log.Println(strconv.Itoa(songId) + ", song not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	lyrics, err := server.TokenHandler.ToLyrics(song.Name, song.SectionJP, song.SectionCHT)
	if err != nil {
		log.Fatal(err)
	}
	lyrics.JPAuthor = song.JPAuthor
	lyrics.CHAuthor = song.CHAuthor

	result := gin.H{"lyrics": lyrics}
	c.JSON(http.StatusOK, result)
}
