package gohtml

import (
	"color-generator/constants"
	"html/template"
	"lib/golog"
	"net/http"
	"path"
	"path/filepath"
)

func FctOutputHTML(w http.ResponseWriter, _filepath string, _templateVar map[string]interface{}, _funcMap template.FuncMap) {

	// ■■■■■■■■■■ Parse template ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvlDebug, "Parsing template from %s", _filepath)
	t, err := template.New(filepath.Base(_filepath)).Funcs(_funcMap).ParseFiles(path.Join(constants.WorkingDir, _filepath))
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
