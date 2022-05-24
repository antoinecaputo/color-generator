package handlers

import (
	"color-generator/constants"
	"color-generator/pkg/colors"
	"color-generator/pkg/gohtml"
	"encoding/json"
	"fmt"
	"github.com/bamboutech/golog"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type TemplateColorTyp struct {
	Index int
	Color colors.ColorTyp
}

func GenerationHandler(w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().UnixNano())

	constants.Log.FctLog(golog.LogLvl_Debug, "Getting user id from cookie")
	c, err := r.Cookie("userId")
	if err != nil {
		if err == http.ErrNoCookie {
			constants.Log.FctLog(golog.LogLvl_Err, "   = No cookie found, generating it")

			c = &http.Cookie{
				Name:    "userId",
				Value:   strconv.Itoa(rand.Int()),
				Expires: time.Now().Add(365 * 24 * time.Hour),
			}
			http.SetCookie(w, c)
		} else {
			constants.Log.FctLog(golog.LogLvl_Err, "   = Error getting cookie: %s", err)
			http.Error(w, "Error getting cookie", http.StatusInternalServerError)
		}
	}
	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK, got %q", c.Value)

	userId := c.Value
	if userId == "" {
		http.Error(w, "No user id provided", http.StatusBadRequest)
		return
	}

	// ■■■■■■■■■■ Colors generation ■■■■■■■■■■

	rand.Seed(time.Now().UnixNano())
	paletteId := rand.Int()

	tmplColors := make([]TemplateColorTyp, 5)
	tmplColorsIndexes := make(map[int]bool, 5)
	var i int
	for i < 5 {
		randIndex := rand.Intn(colors.FctGetColorsLength())
		/*
			if _, exist := tmplColorsIndexes[randIndex]; exist {
				continue
			}
		*/

		color := colors.FctGetColor(randIndex)

		switch i {
		// Title
		case 0:
			// Get a color that is not too close to the primary color
			if 255-color.Luma() < 100 {
				continue
			}
		// Text secondary
		case 2:
			primaryColor := tmplColors[1].Color
			// Get a color that is not too close to the primary color
			if math.Abs(float64(color.Luma()-primaryColor.Luma())) < 125 {
				continue
			}
		}

		tmplColorsIndexes[randIndex] = true

		tmplColors[i] = TemplateColorTyp{
			Index: randIndex,
			Color: color,
		}

		i++
	}

	// ■■■■■■■■■■ HTML Output ■■■■■■■■■■

	templateVar := map[string]interface{}{
		"paletteId": paletteId,
		"colors":    tmplColors,
	}

	funcMap := template.FuncMap{
		"fctColorPosition": colors.FctColorPosition,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	gohtml.FctOutputHTML(w, "/static/index.gohtml", templateVar, funcMap)
}

func EvaluationHandler(w http.ResponseWriter, r *http.Request) {
	CORSHandler(&w, r)

	if r.Method == "OPTIONS" {
		http.Error(w, "", http.StatusUnauthorized)
	}

	if r.Method != "POST" {
		constants.Log.FctLog(golog.LogLvl_Err, "%s %s : Method not allowed", r.Method, r.Host)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// ■■■■■■■■■■ Get userId from cookie ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Getting user id from cookie")
	c, err := r.Cookie("userId")
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = Error getting cookie: %s", err)
		http.Error(w, "Error getting cookie", http.StatusInternalServerError)
	}
	userId := c.Value
	if userId == "" {
		constants.Log.FctLog(golog.LogLvl_Err, "   = No user id provided")
		http.Error(w, "No user id provided", http.StatusBadRequest)
	}
	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK, got %q", c.Value)

	// ■■■■■■■■■■ Get form data ■■■■■■■■■■

	var formData struct {
		PaletteId   string            `json:"paletteId"`
		Evaluations map[string]string `json:"evaluations"`
	}

	constants.Log.FctLog(golog.LogLvl_Debug, "Decoding form data")

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&formData)
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if len(formData.Evaluations) == 0 {
		constants.Log.FctLog(golog.LogLvl_Err, "   = No evaluations provided")
		http.Error(w, "No evaluations provided", http.StatusBadRequest)
	}

	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

	// ■■■■■■■■■■ Opening output file ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Opening output file")

	csvFile, err := os.OpenFile(filepath.Join(constants.WorkingDir, "output.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Error while opening file", http.StatusInternalServerError)
	}
	defer csvFile.Close()

	csvFileInfo, err := csvFile.Stat()
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Error while getting file stat", http.StatusInternalServerError)
		return
	}

	// ■■■■■■■■■■ Writing CSV header ■■■■■■■■■■

	if csvFileInfo.Size() == 0 {
		constants.Log.FctLog(golog.LogLvl_Debug, "   = file don't exist yet, creating file and writing header")
		csvHeader := constants.CSV_HEADER + "\n"
		_, err = csvFile.WriteString(csvHeader)
		if err != nil {
			constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
			http.Error(w, "Error writing file", http.StatusInternalServerError)
		}
	}

	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

	// ■■■■■■■■■■ Writing CSV row ■■■■■■■■■■

	user := userId
	date := time.Now()
	i := 0

	for colorIndexStr, evaluationStr := range formData.Evaluations {
		constants.Log.FctLog(golog.LogLvl_Debug, "Getting color index")
		colorIndex, err := strconv.Atoi(colorIndexStr)
		if err != nil {
			constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

		constants.Log.FctLog(golog.LogLvl_Debug, "Getting evaluation")
		evaluation, err := strconv.Atoi(evaluationStr)
		if err != nil {
			constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

		color := colors.FctGetColor(colorIndex)

		colorName := colors.FctColorPosition(i)

		constants.Log.FctLog(golog.LogLvl_Debug, "Convert hex to rgb")
		rgb, err := colors.FctHex2RGB(color.Value)
		if err != nil {
			constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

		csvRow := fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,%d\n", user, date.Format("2006/01/02"), formData.PaletteId, colorName, rgb.Red, rgb.Blue, rgb.Green, evaluation)

		constants.Log.FctLog(golog.LogLvl_Debug, "Writing CSV row")
		_, err = csvFile.WriteString(csvRow)
		if err != nil {
			constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
		}
		constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

		i++
	}
}

func CORSHandler(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
