package handlers_test

import (
	"github.com/scottapp/misuchiru-server/handlers"
	"log"
	"testing"
)

func Test_ToRuby(t *testing.T) {
	s := "寿司が食べたい。"
	h := &handlers.TokenHandlers{}
	h.ToRuby(s)
}

func Test_ToRuby_02(t *testing.T) {
	s := "君は嬉しそうに　しばらく空を見ていた"
	h := &handlers.TokenHandlers{}
	h.ToRuby(s)
}

func Test_ToLine_1(t *testing.T) {
	s := "寿司が食べたい。"
	h := &handlers.TokenHandlers{}
	line := h.ToLine(s)
	for _, token := range line.Tokens {
		log.Println("word=" + token.Surface + ", reading=" + token.Reading)
	}
	equals(t, 5, len(line.Tokens))
}

func Test_ToLine_2(t *testing.T) {
	s := "8８, 狙う, 狙っている"
	h := &handlers.TokenHandlers{}
	line := h.ToLine(s)
	for _, token := range line.Tokens {
		log.Println("word=" + token.Surface + ", reading=" + token.Reading)
	}
	equals(t, "88", line.Tokens[0].Surface)
	equals(t, 5, len(line.Tokens))
}

func Test_ToLine_english_string(t *testing.T) {
	s := "this is a test string"
	h := &handlers.TokenHandlers{}
	line := h.ToLine(s)
	for _, token := range line.Tokens {
		log.Println("word=" + token.Surface + ", reading=" + token.Reading)
	}
	//equals(t, "88", line.Tokens[0].Surface)
	//equals(t, 5, len(line.Tokens))
}

/*
func Test_ToLyrics_1(t *testing.T) {
	str := []string{"寿司が食べたい。", "君は嬉しそうに　 しばらく空を見ていた"}
	trans := []string{"test translation 1", "test translation 2"}
	h := &handlers.TokenHandlers{}
	s1 := []handlers.Section{handlers.Section{Lines:str}}
	s2 := []handlers.Section{handlers.Section{Lines:trans}}
	lyrics, err := h.ToLyrics(s1, s2)
	ok(t, err)
	for _, section := range lyrics.Sections {
		for _, line := range section.Lines {
			log.Println("new line:")
			for _, token := range line.Tokens {
				log.Println("word=" + token.Surface + ", reading=" + token.Reading)
			}
		}
	}
	equals(t, 5, len(lyrics.Sections[0].Lines[0].Tokens))
	equals(t, 14, len(lyrics.Sections[0].Lines[1].Tokens))
}
*/
