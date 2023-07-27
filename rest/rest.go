package rest

import (
	"GoBlockChain/blockchain"
	"GoBlockChain/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string

// 타입을 설정하고 리시버 함수를 아래에서 구현 인터페이스 implement와 비슷
type URL string

// implement 대신 receiver 함수로 override 느낌
func (u URL) String() string {
	return string("hello" + u)
}

// impelment대신 interface method를 구현하기만 하면 됨
func (receiver URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:4000%s", receiver)
	return []byte(url), nil
}

// struct field tag -> json return 시 보여질 모습
type URLDesc struct {
	URL     URL    `json:"url"`
	Method  string `json:"method"` //method 소문자로 하면 json export 불가능
	Desc    string `json:"desc"`
	Payload string `json:"payload,omitempty"` //omitempty 데이터 없을 시 해당 프로퍼티 안내림  '-' 사용시 필드 무시
}

type addBlockBody struct {
	Message string //todo 소문자로 쓰면 안되네 unmarshal 할 때 소문자면 private이라 못하는 듯
}

func Start(portNo int) {
	port = fmt.Sprintf(":%d", portNo)
	router := mux.NewRouter()

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := []URLDesc{
			{
				URL:     URL("/"), // go는 string이지만 원하는 타입을 만들어 낼수 있음
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

		//json이라는 것을 알리고 byte로 데이터 바꾸고 Fprintf로 writer에 작성함
		writer.Header().Add("Content-Type", "application/json")
		//bytes, err := json.Marshal(data)
		//utils.HandleError(err)
		//fmt.Fprintf(writer, "%s", bytes)

		//Marshal로 바이트 변환할 필요없이 json encoder 사용
		json.NewEncoder(writer).Encode(data)

	}).Methods("GET")

	router.HandleFunc("/blocks", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			writer.Header().Add("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(blockchain.GetBlockchain().ListBlocks())
		case "POST":
			var addBlockBody addBlockBody
			err := json.NewDecoder(request.Body).Decode(&addBlockBody)
			utils.HandleError(err)
			log.Println(addBlockBody)
			//blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
			writer.WriteHeader(http.StatusCreated)
		}

	}).Methods("GET", "POST")

	router.HandleFunc("/blocks/{id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		log.Println(vars)
		id := vars["id"]
		log.Println(id)

	}).Methods("GET")

	log.Println("Running at http://localhost:4000")
	log.Fatal(http.ListenAndServe(port, router))
}
