package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

//Conf 서버 설정 파일
type Conf struct {
	ClientChannelSize int

	ServerName  string
	ServiceType string

	TCPBind string
	UDPBind string

	LogLevel string //로그 출력 레벨

	ReadTimeOutSeconds  int //read 대기 시간
	WriteTimeOutSeconds int //write 대기 시간
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
