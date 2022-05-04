package main

import (
	"color-generator/constants"
	"color-generator/pkg/colors"
	"fmt"
	"github.com/bamboutech/golog"
	"log"
	"net/http"
)

func main() {

	// ■■■■■■■■■■ Log ■■■■■■■■■■

	logger, err := golog.FctCreateLogger(golog.TrcMth_Dual, golog.LogLvl_Debug, "bamboutech", "color-generator.log")
	if err != nil {
		log.Fatalln(err)
		return
	}
	constants.Log = logger

	// ■■■■■■■■■■ Multiplexer ■■■■■■■■■■

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", colors.GenerationHandler)

	mux.HandleFunc("/post", colors.EvaluationHandler)

	// ■■■■■■■■■■ Server ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Info, "---------- Application starting on port %d ----------", constants.PORT)
	// url := "http://127.0.0.1:" + fmt.Sprintf("%d", constants.PORT)
	//	 _ = exec.Command("explorer", url).Run()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", constants.PORT), mux))
}
