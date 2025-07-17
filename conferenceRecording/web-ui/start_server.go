package webui

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

const cssDir = "./web-ui/static"

func StartWebUi() {

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir(cssDir))))
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/", homePageHandler)
	logrus.Info("Server start. Go to  http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logrus.Fatalf("info: can't start server, err_text: %v", err)
	}
}
