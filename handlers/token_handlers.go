package handlers

import (
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/width"
	"log"
	"strconv"
	"strings"
	//"errors"
)

type Token struct {
	Surface string
	Reading string
}

type Line struct {
	Tokens         []*Token
	TranslationCHT string
}

type LyricsSection struct {
	Lines []*Line
}

type Lyrics struct {
	Name     string
	JPAuthor string
	CHAuthor string
	Sections []LyricsSection
}

type TokenHandlers struct {
}

func hira2kata(hira rune) rune {
	if (hira >= 'ぁ' && hira <= 'ゖ') || (hira >= 'ゝ' && hira <= 'ゞ') {
		return hira + 0x60
	}
	return hira
}

func Hira2kata(hira string) string {
	return strings.Map(hira2kata, hira)
}

func kata2hira(kata rune) rune {
	if (kata >= 'ァ' && kata <= 'ヶ') || (kata >= 'ヽ' && kata <= 'ヾ') {
		return kata - 0x60
	}
	return kata
}

func Kata2hira(kata string) string {
	return strings.Map(kata2hira, kata)
}

func (h *TokenHandlers) ToRuby(s string) string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		log.Println(err)
		return ""
	}

	//tokens := t.Tokenize("寿司が食べたい。") // t.Analyze("寿司が食べたい。", tokenizer.Normal)
	tokens := t.Tokenize(strings.TrimSpace(s))
	output := strings.Builder{}

	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			continue
		}
		//log.Println(token.Surface)
		//output.WriteString("<ruby>")

		var hira string

		// process half width characters
		narrowWidth := width.Narrow.String(strings.TrimSpace(token.Surface))
		_, err := strconv.Atoi(narrowWidth)
		if err == nil {
			output.WriteString(narrowWidth)
			continue
		}

		features := token.Features()
		if len(features) >= 9 {
			hira = Kata2hira(strings.TrimSpace(features[7]))
			//log.Println(hira)
			if token.Surface != hira {
				output.WriteString("<ruby>")
				output.WriteString(token.Surface)
				if hira != "" {
					output.WriteString("<rt>")
					output.WriteString(hira)
					output.WriteString("</rt>")
				}
				output.WriteString("</ruby>")
			} else {
				output.WriteString(token.Surface)
			}
		}
	}
	log.Println(output.String())
	return output.String()
}

func (h *TokenHandlers) ToLine(s string) Line {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		log.Println(err)
		return Line{}
	}

	line := Line{}
	line.Tokens = []*Token{}

	tokens := t.Tokenize(strings.TrimSpace(s))
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			continue
		}

		// process half width characters
		narrowWidth := width.Narrow.String(strings.TrimSpace(token.Surface))
		_, err := strconv.Atoi(narrowWidth)
		if err == nil {
			outToken := Token{}
			outToken.Surface = narrowWidth
			line.Tokens = append(line.Tokens, &outToken)
			continue
		}

		var hira string
		features := token.Features()
		if len(features) >= 9 {
			hira = Kata2hira(strings.TrimSpace(features[7]))
			if hira == "" {
				outToken := Token{}
				outToken.Surface = ""
				outToken.Reading = ""
				line.Tokens = append(line.Tokens, &outToken)
			}
			if token.Surface != hira {
				outToken := Token{}
				outToken.Surface = strings.TrimSpace(token.Surface)
				if hira != "" {
					outToken.Reading = hira
				}
				line.Tokens = append(line.Tokens, &outToken)
			} else {
				outToken := Token{}
				outToken.Surface = hira
				line.Tokens = append(line.Tokens, &outToken)
			}
		} else {
			outToken := Token{}
			outToken.Surface = token.Surface
			line.Tokens = append(line.Tokens, &outToken)
		}
	}
	return line
}

func (h *TokenHandlers) ToLyrics(songName string, jpSections []Section, chSections []Section) (*Lyrics, error) {

	lyrics := Lyrics{}
	lyrics.Name = songName

	for i, section := range jpSections {
		outSection := LyricsSection{}
		for j, s := range section.Lines {
			line := h.ToLine(s)
			line.TranslationCHT = chSections[i].Lines[j]
			outSection.Lines = append(outSection.Lines, &line)
		}
		lyrics.Sections = append(lyrics.Sections, outSection)
	}

	/*
		if len(lines) != len(translations) {
			return nil, errors.New("lines count do not match")
		}


		for i:=0; i<len(lines); i++ {
			s := lines[i]
			line := h.ToLine(s)
			line.TranslationCHT = translations[i]
			out.Lines = append(out.Lines, &line)
		}

	*/

	//return &out, nil
	return &lyrics, nil
}
