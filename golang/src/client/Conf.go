package client

import (
	"encoding/json"
	"fmt"
	"os"
)

//Conf 서버 설정 파일
type Conf struct {
	TCPConnect       string
	UDPConnect       string
	ConcurrentAccess int
	Total            int
}

func (c *Conf) Load(filePath string) bool {

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
