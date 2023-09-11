package main

import (
	"bytes"
	"endiatx/lib/golog"
	"fmt"
	"log"
	"os"
)

func main() {
	logger, err := golog.NewLogger(golog.LogLvlDebug)
	if err != nil {
		log.Fatalln(err)
	}

	err = logger.FctAddFile("./log.log")
	if err != nil {
		log.Fatalln(err)
	}

	buffer := new(bytes.Buffer)
	err = logger.FctAddOutput(buffer)
	if err != nil {
		log.Fatalln(err)
	}

	err = logger.FctAddOutput(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Start()

	fmt.Printf("BUFFER : %s\n", buffer.String())
}
