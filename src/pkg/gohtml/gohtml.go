package gohtml

import (
	"color-generator/constants"
	"html/template"
	"lib/golog"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func FctOutputHTML(w http.ResponseWriter, _staticFilepath string, _templateVar map[string]interface{}, _funcMap template.FuncMap) {

	// ■■■■■■■■■■ Parse template ■■■■■■■■■■

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}

	executableDir := filepath.Dir(executablePath)

	templatePath := path.Join(executableDir, "static", _staticFilepath)

	constants.Log.FctLog(golog.LogLvlDebug, "Parsing template from %s", templatePath)
	t, err := template.New(filepath.Base(templatePath)).Funcs(_funcMap).ParseFiles(templatePath)
	if err != nil {
		constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

	// ■■■■■■■■■■ Execute template ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvlDebug, "Executing template with params %+v", _templateVar)
	if err = t.Execute(w, _templateVar); err != nil {
		constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
