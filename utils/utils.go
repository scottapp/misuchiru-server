package utils

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Join(filepath.Dir(b), "../")
)

func ParseTimestamp(str string) (string, error) {
	parts := strings.Split(str, "/")
	if len(parts) == 0 {
		return "", errors.New("error timestamp string")
	}
	last := parts[len(parts)-1]
	if strings.HasSuffix(last, ".html") {
		tmp2 := strings.Split(last, ".")
		if len(tmp2) != 5 {
			return "", errors.New("error timestamp string")
		}
		return tmp2[1], nil
	} else {
		return "", errors.New("error timestamp string")
	}
}

func GetTime(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

func getPathByTime(timestamp string) string {
	tm, err := GetTime(timestamp)
	if err != nil {
		panic(err)
	}
	year, month, day := tm.Date()
	monthStr := fmt.Sprintf("%02d", int(month))
	return strconv.Itoa(year) + "/" + monthStr + "/" + strconv.Itoa(day)
}

func MakePath(baseDir string, tm time.Time) string {
	dateDir := "./" + baseDir + "/" + strconv.Itoa(tm.Year()) + "/" + tm.Format("2006-01-02")
	return dateDir
}

func CreateDateDir(baseDir string, tm time.Time) (string, error) {
	dateDir := MakePath(baseDir, tm)
	//fmt.Println(dateDir)
	err := os.MkdirAll(dateDir, 0744)
	if err != nil {
		return "", err
	}
	return dateDir, nil
}

func LocalFileExists(baseDir string, url string) (bool, error) {
	items := strings.Split(url, "/")
	last := items[len(items)-1]

	timestamp, err := ParseTimestamp(url)
	if err != nil {
		return false, err
	}
	tm, err := GetTime(timestamp)
	if err != nil {
		return false, err
	}

	dirPath := MakePath(baseDir, tm)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}

	fullPath := dirPath + "/" + last
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func ReadLines(dataDir string, fileName string) ([]string, error) {
	f, err := os.Open(Basepath + dataDir + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return lines, nil
}

func ReadFile(filePath string) []string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(string(dat[:])))
	urls := []string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, strings.Trim(record[2], " "))
	}
	return urls
}

func GetChunks(items []string, chunkSize int) [][]string {
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

func StripString(str string, replacement string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(str, replacement)
}

func GetSortedValueInts(aMap map[string]int) []int {
	vals := []int{}
	for _, val := range aMap {
		vals = append(vals, val)
	}
	sort.Ints(vals)
	return vals
}
