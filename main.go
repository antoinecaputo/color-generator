package main

import (
	"color-generator/constants"
	"color-generator/pkg/handlers"
	"fmt"
	"github.com/bamboutech/golog"
	"log"
	"net/http"
	"path"
)

func main() {

	// ■■■■■■■■■■ Log ■■■■■■■■■■

	logger, err := golog.FctCreateLogger(golog.TrcMth_File, constants.LogLevel, "bamboutech", "color-generator.log")
	if err != nil {
		log.Fatalln(err)
		return
	}
	constants.Log = logger

	// ■■■■■■■■■■ Multiplexer ■■■■■■■■■■

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(path.Join(constants.WorkingDir, "/static")))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.GenerationHandler)

	mux.HandleFunc("/post", handlers.EvaluationHandler)

	// ■■■■■■■■■■ Server ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Info, "---------- Application starting on port %d ----------", constants.PORT)
	// url := fmt.Sprintf("%s:%d", constants.IP, constants.PORT)
	// fmt.Printf("Application starts on %s with default user 0\n", url)
	// fmt.Printf("http://%s", url)
	//	 _ = exec.Command("explorer", url).Run()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", constants.PORT), mux))
}
