package main

import (
	"color-generator/constants"
	"color-generator/pkg/handlers"
	"fmt"
	"lib/golog"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {

	// ■■■■■■■■■■ Log ■■■■■■■■■■

	logger, err := golog.NewLogger(constants.LogLevel)
	if err != nil {
		log.Fatalln(err)
		return
	}
	constants.Log = logger

	err = constants.Log.FctAddOutput(os.Stdout)
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = constants.Log.FctAddFile("color-generator.log")
	if err != nil {
		log.Fatalln(err)
		return
	}

	constants.Log.Start()

	// ■■■■■■■■■■ Multiplexer ■■■■■■■■■■

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(path.Join(constants.WorkingDir, "/static")))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.GenerationHandler)

	mux.HandleFunc("/post", handlers.EvaluationHandler)

	// ■■■■■■■■■■ Server ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvlInfo, "---------- Application starting on port %d ----------", constants.PORT)
	url := fmt.Sprintf("%s:%d", constants.IP, constants.PORT)
	fmt.Printf("http://%s", url)
	// _ = exec.Command("explorer", url).Run()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", constants.PORT), mux))
}
