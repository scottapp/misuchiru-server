package server_test

import (
	"bytes"
	"encoding/json"
	//"github.com/gin-gonic/gin"
	//"github.com/gomodule/redigo/redis"
	//"github.com/rafaeljusto/redigomock"
	//"github.com/stretchr/testify/assert"
	"github.com/scottapp/misuchiru-server/server"
	"log"
	//"net/http"
	//"net/http/httptest"
	"testing"
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func Test_NewServer(t *testing.T) {
	server := server.NewServer()
	log.Println(server)
}
