package main

import (
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	batch_crawler "github.com/pytorchtw/go-batch-crawler"
	"github.com/scottapp/misuchiru-server/handlers"
	"github.com/scottapp/misuchiru-server/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ResponseHandler struct {
	Crawler        *batch_crawler.Crawler
	Result         *handlers.AlbumDB
	Errors         []error
	ContentHandler *handlers.MCHandler
}

func NewResponseHandler() (*ResponseHandler, error) {
	h := ResponseHandler{}
	h.Result = &handlers.AlbumDB{}
	h.Result.Albums = map[string]*handlers.Album{}
	h.Result.SongByName = map[string]*handlers.Song{}
	h.Result.SongByURL = map[string]*handlers.Song{}
	h.ContentHandler = &handlers.MCHandler{}
	return &h, nil
}

func (h *ResponseHandler) findSongByURL(url string) (*handlers.Song, error) {
	song, ok := h.Result.SongByURL[url]
	if !ok {
		return nil, errors.New(url + " not found")
	}
	return song, nil
}

func (h *ResponseHandler) HandleResponse(r *colly.Response) {
	songURL := strings.TrimLeft(r.Request.URL.String(), "https://")
	songURL = strings.TrimLeft(songURL, "http://")
	log.Println("OK:", songURL)

	doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
	if err != nil {
		log.Println(err)
		h.Errors = append(h.Errors, err)
		return
	}

	//s1, s2, err := h.ContentHandler.ParseHTML(string(r.Body))
	s1, s2, err := h.ContentHandler.ParseHTML(doc)
	if err == nil {
		song, err := h.findSongByURL(songURL)
		if err != nil {
			log.Fatal(err)
		}

		jpAuthor := htmlquery.FindOne(doc, `//td[@id="jp_author"]`)
		if jpAuthor == nil {
			h.Errors = append(h.Errors, errors.New("jp author not found"))
		}

		chAuthor := htmlquery.FindOne(doc, `//td[@id="ch_author"]`)
		if chAuthor == nil {
			h.Errors = append(h.Errors, errors.New("ch author not found"))
		}

		song.SectionJP = s1
		song.SectionCHT = s2
		song.JPAuthor = htmlquery.InnerText(jpAuthor)
		song.CHAuthor = htmlquery.InnerText(chAuthor)
		return
	}

	albumTitle := htmlquery.FindOne(doc, `//span[@class="titlename"]`)
	if albumTitle == nil {
		h.Errors = append(h.Errors, errors.New("album title not found"))
		return
	}
	album := handlers.Album{}
	album.URL = r.Request.URL.String()
	album.URLIDTable = map[string]int{}
	album.SongByID = map[int]*handlers.Song{}
	album.Name = htmlquery.InnerText(albumTitle)
	h.Result.Albums[album.Name] = &album

	jpAuthor := htmlquery.FindOne(doc, `//td[@id="jp_author"]`)
	if jpAuthor == nil {
		h.Errors = append(h.Errors, errors.New("jp author not found"))
	}

	chAuthor := htmlquery.FindOne(doc, `//td[@id="ch_author"]`)
	if chAuthor == nil {
		h.Errors = append(h.Errors, errors.New("ch author not found"))
	}

	table := htmlquery.FindOne(doc, `//table[@id="lyric_album"]`)
	if table == nil {
		h.Errors = append(h.Errors, errors.New("album not found"))
		return
	}

	links := htmlquery.Find(table, `//a`)
	count := 1
	for _, link := range links {
		parentC := htmlquery.SelectAttr(link.Parent, "class")
		if parentC == "catalog_track" || parentC == "lang-jp" {
			originURL := htmlquery.SelectAttr(link, "href")
			curSongURL := htmlquery.SelectAttr(link, "href")
			curSongURL = strings.TrimLeft(curSongURL, "https://")
			curSongURL = strings.TrimLeft(curSongURL, "http://")

			songName := htmlquery.InnerText(link)
			song, _ := h.findSongByURL(curSongURL)
			if song != nil {
				log.Println("found, skip saving song " + songName + ", " + curSongURL)
				continue
			}

			log.Println("creating new song data " + strconv.Itoa(count) + ", " + songName)
			song = &handlers.Song{}
			song.ID = count
			song.Name = songName
			song.Album = album.Name
			song.URL = curSongURL
			album.URLIDTable[curSongURL] = count
			album.SongByID[song.ID] = song
			h.Result.SongByName[songName] = song
			h.Result.SongByURL[curSongURL] = song
			h.Crawler.Queue.AddURL(originURL)
			count++
			log.Println(strconv.Itoa(count)+", crawling "+song.Name, song.URL)
		}
	}
}

func (h *ResponseHandler) HandleError(r *colly.Response, err error) {
	log.Println("ERROR:", r.Request.URL.String(), r.StatusCode)
	log.Println(err.Error())
}

func getData(filename string) ([]string, error) {
	lines, err := utils.ReadLines("data", filename)
	if err != nil {
		return nil, err
	}

	boards := []string{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			return nil, errors.New("error parsing data file")
		}
		boards = append(boards, strings.Trim(strings.TrimSpace(parts[1]), ","))
	}
	return boards, nil
}

func getChunks(items []string, chunkSize int) [][]string {
	chunks := [][]string{}
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

func saveAlbumsJson(result *handlers.AlbumDB) {
	for _, album := range result.Albums {
		log.Println("saving " + album.Name)
		file, err := json.MarshalIndent(album, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		//err = ioutil.WriteFile(utils.Basepath+"/data/"+album.Name+".json", file, 0644)
		err = ioutil.WriteFile(utils.Basepath+"/data/"+album.Name+".json", file, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func tempCreateAlbums() {
	albumNames := []string{
		"EVERYTHING",
		"Kind of Love",
		"Versus",
		"Atomic Heart",
		"深海",
		"BOLERO",
		"DISCOVERY",
		"1/42",
		"Q",
		"It's a Wonderful World",
		"シフクノオト",
		"I ♥ U",
		"HOME",
		"B-SIDE",
		"Supermarket Fantasy",
		"SENSE",
		"[(an imitation) blood orange]",
		"REFLECTION",
		"重力と呼吸",
	}

	result := handlers.AlbumDB{}
	result.Albums = map[string]*handlers.Album{}
	count := 1
	for _, name := range albumNames {
		album := handlers.Album{}
		album.ID = count
		album.Name = name
		count++
		result.Albums[name] = &album
	}
	log.Println("saving all albums")
	file, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(utils.Basepath+"/data/all_albums.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	/*
		//tempCreateAlbums()
		result := handlers.ReadAllAlbums()
		for albumName, album := range result.Albums {
			vals := utils.GetSortedValueInts(album.URLIDTable)
			for _, songID := range vals {
				log.Println(albumName, songID, album.SongByID[songID].Name)
			}
		}
		os.Exit(0)
	*/

	start := time.Now()

	file, err := os.OpenFile("log/info_"+start.Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)

	respHandler, err := NewResponseHandler()
	if err != nil {
		panic(err)
	}

	crawler := batch_crawler.NewCrawler([]string{""}, respHandler)
	respHandler.Crawler = crawler

	err = crawler.C.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "https://",
		Parallelism: 1,
		// Set a delay between requests to these domains
		Delay: 0 * time.Second,
		// Add an additional random delay
		RandomDelay: 0 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	myTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   20 * time.Second,
			KeepAlive: 20 * time.Second,
		}).DialContext,
		MaxIdleConns:          25,
		IdleConnTimeout:       20 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
	crawler.C.WithTransport(myTransport)
	crawler.C.SetRequestTimeout(5 * time.Second)

	// 19
	//crawler.Queue.AddURL("https://blog.xuite.net/lyricbox/2601/585394761")

	// 18, no lyrics
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/306922140")

	// 17
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/63883386")

	// 16
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/40380041")

	// 15
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361347")

	// 14
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361614")

	// 13
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361604")

	// 12
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361588")

	// 11
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361582")

	// 10
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361573")

	// 9
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361530")

	// 8
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361521")

	// 7
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361509")

	// 6
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361501")

	// 5
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361491")

	// 4
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361479")

	// 3
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361467")

	// 2
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361459")

	// 1
	//crawler.Queue.AddURL("http://blog.xuite.net/lyricbox/2601/21361450")

	//crawler.Run()
	//saveAlbumsJson(respHandler.Result)

	log.Println("running time: " + time.Since(start).String())
}
