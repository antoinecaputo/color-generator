package gohtml

import (
	"color-generator/constants"
	"color-generator/pkg/colors"
	"fmt"
	"github.com/bamboutech/golog"
	"html/template"
	"math/rand"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	tmplColors := make([]colors.ColorTyp, 5)
	for i := 0; i < 5; i++ {
		tmplColors[i] = colors.GetColor(rand.Intn(colors.ColorsLength))
	}

	templateVar := map[string]interface{}{
		"colors": tmplColors,
	}

	// ■■■■■■■■■■ HTML Output ■■■■■■■■■■

	fctOutputHTML(w, "./static/index.gohtml", templateVar)

	buff, err := constants.Log.FctGetBufferContent()
	fmt.Println(err)
	fmt.Println(buff)
}

func fctOutputHTML(w http.ResponseWriter, _filepath string, _templateVar map[string]interface{}) {
	// ■■■■■■■■■■ Parse template ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Parsing template from %s", _filepath)
	t, err := template.ParseFiles(_filepath)
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

	// ■■■■■■■■■■ Execute template ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Executing template with params %+v", _templateVar)
	if err = t.Execute(w, _templateVar); err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
