package main

import (
	"GoBlockChain/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

type URL string

// impelment대신 interface method를 구현하기만 하면 됨
func (receiver URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:4000%s", receiver)
	return []byte(url), nil
}

// struct field tag -> json return 시 보여질 모습
type URLDesc struct {
	URL     URL    `json:"url"`
	Method  string `json:"method"` //소문자로 하면 json export 불가능
	Desc    string `json:"desc"`
	Payload string `json:"payload,omitempty"` //omitempty 데이터 없을 시 해당 프로퍼티 안내림  '-' 사용시 필드 무시
}

func (receiver URLDesc) String() string {
	return "hi"
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := []URLDesc{
			{
				URL:     "/",
				Method:  "GET",
				Desc:    "See Doc",
				Payload: "",
			},
			{
				URL:     "/blocks",
				Method:  "POST",
				Desc:    "Add a block",
				Payload: "data:string",
			},
		}

		bytes, err := json.Marshal(data)
		utils.HandleError(err)

		writer.Header().Add("Content-Type", "application/json")
		fmt.Printf("%s", bytes) //byte -> string

		//Marshal로 바이트 변환할 필요없이 json encoder 사용
		json.NewEncoder(writer).Encode(data)

	})

	log.Println("Running at http://localhost:4000")
	log.Fatal(http.ListenAndServe(port, nil))

}
