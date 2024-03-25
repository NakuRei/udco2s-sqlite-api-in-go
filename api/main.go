package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ApiConfig struct {
	DBPath    string
	TableName string
	Port      string
}

type SensorData struct {
	ID                 int     `json:"id"`
	CO2PPM             int     `json:"co2_ppm"`
	HumidityPercentage float64 `json:"humidity_percentage"`
	TemperatureCelsius float64 `json:"temperature_celsius"`
	Timestamp          string  `json:"timestamp"`
}

func loadConfig() ApiConfig {
	config := ApiConfig{
		DBPath:    os.Getenv("DB_PATH"),
		TableName: os.Getenv("TABLE_NAME"),
		Port:      os.Getenv("PORT"),
	}

	if config.DBPath == "" {
		config.DBPath = "./db/udco2s_data.db"
	}
	if config.TableName == "" {
		config.TableName = "sensor_data"
	}
	if config.Port == "" {
		config.Port = "8080"
	}
	return config
}

func initDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
}

// parseAndFormatDateTimeはISO 8601形式の日付と時刻の文字列を受け取り、
// UTCタイムゾーンに基づいてフォーマットされた新しい文字列を返す。
func parseAndFormatDateTime(dateTimeStr string) (string, error) {
	const layout = time.RFC3339 // ISO 8601フォーマット
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return "", err
	}
	// UTCに変換し、データベースが期待するフォーマットにする
	return t.UTC().Format("2006-01-02 15:04:05"), nil
}

func sensorDataHandler(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// GETリクエスト以外は許可しない
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// クエリパラメータの取得
		startParam := r.URL.Query().Get("start")
		endParam := r.URL.Query().Get("end")
		if startParam == "" || endParam == "" {
			http.Error(w, "Start and end query parameters are required", http.StatusBadRequest)
			return
		}

		// 日付と時刻の文字列をUTCフォーマットに変換
		start, err := parseAndFormatDateTime(startParam)
		if err != nil {
			http.Error(w, "Invalid start time format", http.StatusBadRequest)
			return
		}
		end, err := parseAndFormatDateTime(endParam)
		if err != nil {
			http.Error(w, "Invalid end time format", http.StatusBadRequest)
			return
		}

		// データの取得
		data, err := getSensorDataByDateTimeRange(db, tableName, start, end)
		if err != nil {
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
			return
		}

		// 結果が空の場合は、404を返す
		if len(data) == 0 {
			http.Error(w, "No data found for the specified range", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func getSensorDataByDateTimeRange(db *sql.DB, tableName string, startDateTime, endDateTime string) ([]SensorData, error) {
	// NOTE: クエリのパラメータを直接埋め込むのはセキュリティ上のリスクがあるため
	// NOTE: このやり方はtableNameによるSQLインジェクションに対して脆弱である
	query := fmt.Sprintf("SELECT id, co2_ppm, humidity_percentage, temperature_celsius, timestamp FROM %s WHERE timestamp BETWEEN ? AND ?", tableName)
	rows, err := db.Query(query, startDateTime, endDateTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []SensorData
	for rows.Next() {
		var d SensorData
		if err := rows.Scan(&d.ID, &d.CO2PPM, &d.HumidityPercentage, &d.TemperatureCelsius, &d.Timestamp); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func main() {
	config := loadConfig()
	db := initDB(config.DBPath)
	defer db.Close()

	http.HandleFunc("/sensor_data", sensorDataHandler(db, config.TableName))

	log.Printf("Server started on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
