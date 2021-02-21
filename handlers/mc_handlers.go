package handlers

import (
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/scottapp/misuchiru-server/utils"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
)

/*
type Song struct {
	URL string
	Name string
	SectionJP []Section
	SectionCHT []Section
}
*/

/*
type Album struct {
	//URL []string
	Name string
	Songs map[string]*Song
}
*/

type Section struct {
	Lines []string
}

type MCHandler struct {
}

func (h *MCHandler) ParseAlbumJson(albumName string) (*Album, error) {
	data, err := ioutil.ReadFile(utils.Basepath + "/data/" + albumName + ".json")
	if err != nil {
		log.Fatal(err)
	}

	album := Album{}
	err = json.Unmarshal(data, &album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

//func (h *MCHandler) ParseHTML(content string) ([]Section, []Section, error) {
func (h *MCHandler) ParseHTML(doc *html.Node) ([]Section, []Section, error) {
	/*
		doc, err := htmlquery.Parse(strings.NewReader(content))
		if err != nil {
			return nil, nil, err
		}
	*/

	jp := htmlquery.FindOne(doc, `//td[@id="lyric_jp"]`)
	if jp == nil {
		return nil, nil, errors.New("jp lyrics not found!")
	}

	sections := []Section{}
	nodes := htmlquery.Find(jp, `//p`)
	for _, node := range nodes {
		s := Section{}
		s.Lines = []string{}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				s.Lines = append(s.Lines, c.Data)
			}
		}
		sections = append(sections, s)
	}

	ch := htmlquery.FindOne(doc, `//td[@id="lyric_ch"]`)
	if ch == nil {
		return nil, nil, errors.New("ch lyrics not found!")
	}

	sections2 := []Section{}
	nodes2 := htmlquery.Find(ch, `//p`)
	for _, node := range nodes2 {
		s := Section{}
		s.Lines = []string{}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				s.Lines = append(s.Lines, c.Data)
			}
		}
		sections2 = append(sections2, s)
	}

	return sections, sections2, nil
}
