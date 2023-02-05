package explorer

import (
	"GoBlockChain/blockchain"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string // 소문자로 하면 private이라 렌더링 데이터로 못넘김
	Blocks    []*blockchain.Block
}

func StartExplorer() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := homeData{"Home", blockchain.GetBlockchain().GetAllBlocks()}
		templates.ExecuteTemplate(writer, "home", data)
	})

	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			templates.ExecuteTemplate(writer, "add", nil)
		case "POST":
			request.ParseForm()
			inputData := request.Form.Get("blockData")
			blockchain.GetBlockchain().AddBlock(inputData)
			http.Redirect(writer, request, "/", http.StatusPermanentRedirect)

		}
	})

	fmt.Printf("server learning at %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
