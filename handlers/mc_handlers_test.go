package handlers_test

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		//fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		log.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

/*
func TestMCHandler_ParseHTML(t *testing.T) {
	url := "https://blog.xuite.net/lyricbox/2601/585471371"
	resp, err := http.Get(url)
	ok(t, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	mh := &handlers.MCHandler{}
	sections, sections2, err := mh.ParseHTML(string(body))
	ok(t, err)

	for i, section := range sections {
		for j, line := range section.Lines {
			log.Println(line)
			log.Println(sections2[i].Lines[j])
		}
		log.Println("")
	}
}
*/

/*
func TestMCHandler_ParseAlbumJson(t *testing.T) {
	mh := &handlers.MCHandler{}
	name := "重力と呼吸"
	album, err := mh.ParseAlbumJson(name)
	ok(t, err)
	for _, song := range album.Songs {
		for _, section := range song.SectionJP {
			for _, line := range section.Lines {
				log.Println(line)
			}
			log.Println("--- end section -------")
		}
	}
}
*/

/*
func TestMCHandler_ParseAlbumJson_2(t *testing.T) {
	mh := &handlers.MCHandler{}
	name := "重力と呼吸"
	album, err := mh.ParseAlbumJson(name)
	ok(t, err)

	songURLs := []string{}
	for key, _ := range album.Songs {
		songURLs = append(songURLs, key)
	}
	songURL := songURLs[2]
	log.Println(songURL)

	song := album.Songs[songURL]
	log.Println(song.Name)
	for i, section := range song.SectionJP {
		for j, line := range section.Lines {
			log.Println(line)
			log.Println(song.SectionCHT[i].Lines[j])
		}
		log.Println("--- end section -------")
	}
}
*/
