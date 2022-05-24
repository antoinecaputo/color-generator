package gohtml

import (
	"color-generator/constants"
	"github.com/bamboutech/golog"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
)

func FctOutputHTML(w http.ResponseWriter, _filepath string, _templateVar map[string]interface{}, _funcMap template.FuncMap) {

	// ■■■■■■■■■■ Parse template ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Parsing template from %s", _filepath)
	t, err := template.New(filepath.Base(_filepath)).Funcs(_funcMap).ParseFiles(path.Join(constants.WorkingDir, _filepath))
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
