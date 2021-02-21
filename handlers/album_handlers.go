package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

type Song struct {
	ID         int
	URL        string
	Name       string
	Album      string
	SectionJP  []Section
	SectionCHT []Section
	JPAuthor   string
	CHAuthor   string
}

type Album struct {
	ID         int
	URL        string
	Name       string
	SongByID   map[int]*Song
	URLIDTable map[string]int
}

type AlbumDB struct {
	Albums     map[string]*Album
	AlbumByID  map[int]*Album
	SongByName map[string]*Song
	SongByURL  map[string]*Song
}

func ReadAlbum(albumID int) (*Album, error) {
	fileName := "./data/album_" + strconv.Itoa(albumID) + ".json"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	album := Album{}
	err = json.Unmarshal(data, &album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func ReadAllAlbums() *AlbumDB {
	fileName := "./data/all_albums.json"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	albumDB := AlbumDB{}
	err = json.Unmarshal(data, &albumDB)
	if err != nil {
		log.Fatal(err)
	}

	if albumDB.AlbumByID == nil {
		albumDB.AlbumByID = map[int]*Album{}
	}

	for _, album := range albumDB.Albums {
		albumFromFile, err := ReadAlbum(album.ID)
		if err != nil {
			//log.Println("error reading album, " + album.Name + ", ID=" + strconv.Itoa(album.ID))
			continue
		}
		albumFromFile.ID = album.ID
		albumDB.Albums[albumFromFile.Name] = albumFromFile
		albumDB.AlbumByID[album.ID] = albumFromFile
		log.Println("loaded " + strconv.Itoa(album.ID) + ", " + album.Name)
	}
	return &albumDB
}
