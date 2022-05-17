package handlers

import (
	"color-generator/constants"
	"color-generator/pkg/colors"
	"color-generator/pkg/gohtml"
	"encoding/json"
	"fmt"
	"github.com/bamboutech/golog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GenerationHandler(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "No user id provided", http.StatusBadRequest)
		return
	}

	// ■■■■■■■■■■ Colors generation ■■■■■■■■■■

	rand.Seed(time.Now().UnixNano())

	tmplColors := make([]interface{}, 5)
	tmplColorsIndexes := make(map[int]bool, 5)
	var i int
	for i < 5 {
		randIndex := rand.Intn(colors.FctGetColorsLength())
		if _, exist := tmplColorsIndexes[randIndex]; exist {
			continue
		}
		tmplColorsIndexes[randIndex] = true

		tmplColors[i] = struct {
			Index int
			Color colors.ColorTyp
		}{
			Index: randIndex,
			Color: colors.FctGetColor(randIndex),
		}
		i++
	}

	// ■■■■■■■■■■ HTML Output ■■■■■■■■■■

	templateVar := map[string]interface{}{
		"userId": userId,
		"colors": tmplColors,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	gohtml.FctOutputHTML(w, "./static/index.gohtml", templateVar)
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

	// ■■■■■■■■■■ Get form data ■■■■■■■■■■

	var formData struct {
		UserId      string            `json:"user_id"`
		Evaluations map[string]string `json:"evaluations"`
	}

	constants.Log.FctLog(golog.LogLvl_Debug, "Decoding form data")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&formData)
	if err != nil {
		constants.Log.FctLog(golog.LogLvl_Err, "   = %s", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if formData.UserId == "" {
		constants.Log.FctLog(golog.LogLvl_Err, "   = No user id provided")
		http.Error(w, "No user id provided", http.StatusBadRequest)
	}

	if len(formData.Evaluations) == 0 {
		constants.Log.FctLog(golog.LogLvl_Err, "   = No evaluations provided")
		http.Error(w, "No evaluations provided", http.StatusBadRequest)
	}

	constants.Log.FctLog(golog.LogLvl_Debug, "   = OK")

	// ■■■■■■■■■■ Opening output file ■■■■■■■■■■

	constants.Log.FctLog(golog.LogLvl_Debug, "Opening output file")

	csvFile, err := os.OpenFile("output.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	user := formData.UserId
	date := time.Now()
	i := 1

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
		var colorName string
		switch i {
		case 1:
			colorName = constants.COLOR_1
		case 2:
			colorName = constants.COLOR_2
		case 3:
			colorName = constants.COLOR_3
		case 4:
			colorName = constants.COLOR_4
		case 5:
			colorName = constants.COLOR_5
		}

		csvRow := fmt.Sprintf("%s,%s,%s,%s,%d\n", user, date.Format("2006/02/01"), colorName, color.Value, evaluation)

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
