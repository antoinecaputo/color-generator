package handlers

import (
	"color-generator/constants"
	"color-generator/pkg/colors"
	"color-generator/pkg/gohtml"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"lib/golog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const CookieUserId = "userId"

func GenerationHandler(w http.ResponseWriter, r *http.Request) {

	constants.Log.FctLog(golog.LogLvlDebug, "Getting user id from cookie")
	c, err := r.Cookie(CookieUserId)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			constants.Log.FctLog(golog.LogLvlErr, "   = No cookie found, generating it")

			randomUUID := uuid.New()

			c = &http.Cookie{
				Name:    CookieUserId,
				Value:   randomUUID.String(),
				Expires: time.Now().Add(365 * 24 * time.Hour),
			}
			http.SetCookie(w, c)
		} else {
			constants.Log.FctLog(golog.LogLvlErr, "   = Error getting cookie: %s", err)
			http.Error(w, "Error getting cookie", http.StatusInternalServerError)
		}
	}
	constants.Log.FctLog(golog.LogLvlDebug, "   = OK, got %q", c.Value)

	userId := c.Value
	if userId == "" {
		http.Error(w, "No user id provided", http.StatusBadRequest)
		return
	}

	// ■■■■■■■■■■ Colors generation ■■■■■■■■■■

	paletteId, palette := colors.FctGeneratePalette()

	// ■■■■■■■■■■ HTML Output ■■■■■■■■■■

	templateVar := map[string]interface{}{
		"paletteId": paletteId.String(),
		"palette":   palette,
	}

	funcMap := template.FuncMap{
		"fctColorPosition": colors.FctColorPosition,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	gohtml.FctOutputHTML(w, "/static/index.gohtml", templateVar, funcMap)
}

func EvaluationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		CORSHandler(&w, r)
		return
	}

	if r.Method != http.MethodPost {
		constants.Log.FctLog(golog.LogLvlErr, "%s %s : Method not allowed", r.Method, r.Host)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// ■■■■■■■■■■ Get userId from cookie ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvlDebug, "Getting user id from cookie")
	c, err := r.Cookie(CookieUserId)
	if err != nil {
		// do not throw http: named cookie not present
		if !errors.Is(err, http.ErrNoCookie) {
			constants.Log.FctLog(golog.LogLvlErr, "   = Error getting cookie: %s", err)
			http.Error(w, "Error getting cookie", http.StatusInternalServerError)
			return
		}
	}

	var userId string
	if c != nil {
		userId = c.Value
	}

	if userId == "" {
		constants.Log.FctLog(golog.LogLvlErr, "   = No user id provided")
		http.Error(w, "No user id provided", http.StatusBadRequest)
		return
	}

	constants.Log.FctLog(golog.LogLvlDebug, "   = OK, got %q", c.Value)

	// ■■■■■■■■■■ Get form data ■■■■■■■■■■

	var formData struct {
		PaletteId   string            `json:"paletteId"`
		Evaluations map[string]string `json:"evaluations"`
	}

	constants.Log.FctLog(golog.LogLvlDebug, "Decoding form data")

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&formData)
	if err != nil {
		constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if len(formData.Evaluations) == 0 {
		constants.Log.FctLog(golog.LogLvlErr, "   = No evaluations provided")
		http.Error(w, "No evaluations provided", http.StatusBadRequest)
	}

	constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

	// ■■■■■■■■■■ Opening output file ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvlDebug, "Opening output file")

	csvFile, err := os.OpenFile(filepath.Join(constants.WorkingDir, "output.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
		http.Error(w, "Error while opening file", http.StatusInternalServerError)
	}
	defer csvFile.Close()

	csvFileInfo, err := csvFile.Stat()
	if err != nil {
		constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
		http.Error(w, "Error while getting file stat", http.StatusInternalServerError)
		return
	}

	// ■■■■■■■■■■ Writing CSV header ■■■■■■■■■■

	if csvFileInfo.Size() == 0 {
		constants.Log.FctLog(golog.LogLvlDebug, "   = file don't exist yet, creating file and writing header")
		csvHeader := constants.CSV_HEADER + "\n"
		_, err = csvFile.WriteString(csvHeader)
		if err != nil {
			constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
			http.Error(w, "Error writing file", http.StatusInternalServerError)
		}
	}

	constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

	// ■■■■■■■■■■ Writing CSV row ■■■■■■■■■■

	user := userId
	date := time.Now()
	i := 0

	for colorIndexStr, evaluationStr := range formData.Evaluations {
		constants.Log.FctLog(golog.LogLvlDebug, "Getting color index")
		colorIndex, err := strconv.Atoi(colorIndexStr)
		if err != nil {
			constants.Log.FctLog(golog.LogLvlErr, "   = %s", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

		constants.Log.FctLog(golog.LogLvlDebug, "Getting evaluation")
		evaluation, err := strconv.Atoi(evaluationStr)
		if err != nil {
			constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

		color := colors.FctGetColor(colorIndex)

		colorName := colors.FctColorPosition(i)

		constants.Log.FctLog(golog.LogLvlDebug, "Convert hex to rgb")
		rgb, err := colors.FctHex2RGB(color.Value)
		if err != nil {
			constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

		csvRow := fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,%d\n", user, date.Format("2006/01/02"), formData.PaletteId, colorName, rgb.Red, rgb.Blue, rgb.Green, evaluation)

		constants.Log.FctLog(golog.LogLvlDebug, "Writing CSV row")
		_, err = csvFile.WriteString(csvRow)
		if err != nil {
			constants.Log.FctLog(golog.LogLvlErr, "   = %s", err.Error())
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
		}
		constants.Log.FctLog(golog.LogLvlDebug, "   = OK")

		i++
	}
}

func CORSHandler(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
