package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type StationInJSON struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type LineInJSON struct {
	Name       string   `json:"name"`
	StationIDs []string `json:"stations"`
}

type JSONData struct {
	Stations []StationInJSON `json:"stations"`
	Lines    []LineInJSON    `json:"lines"`
}

func loadFromJSON(path string) (*JSONData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open data file: %w", err)
	}

	var rv *JSONData
	if err := json.Unmarshal(data, &rv); err != nil {
		return nil, fmt.Errorf("unmarhsal json: %w", err)
	}

	return rv, nil
}

func writeToDB(db *sqlx.DB, data *JSONData) error {
	writeToStations := func() error {
		q := "INSERT INTO stations (id, name, longitude, latitude) VALUES ($1, $2, $3, $4)"

		for _, station := range data.Stations {
			_, err := db.Exec(q, station.ID, station.Name, station.Longitude, station.Latitude)
			if err != nil {
				log.Printf("Writing %v to station: %s", station, err)
				return err
			}
		}
		return nil
	}
	writeToLineHasStations := func() error {
		q := "INSERT INTO line_has_stations (line_name, station_id) VALUES ($1, $2)"
		for _, line := range data.Lines {
			for _, stationID := range line.StationIDs {
				_, err := db.Exec(q, line.Name, stationID)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := writeToStations(); err != nil {
		return fmt.Errorf("write to stations: %w", err)
	}

	if err := writeToLineHasStations(); err != nil {
		return fmt.Errorf("write to line_has_stations: %w", err)
	}

	return nil
}

func PrintLineNamesByStationName(db *sqlx.DB, name string) error {
	q := "SELECT lhs.line_name FROM stations LEFT JOIN line_has_stations lhs on stations.id = lhs.station_id WHERE name = $1"
	row, err := db.Query(q, name)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer row.Close()
	fmt.Printf("Lines:\n")
	for row.Next() {
		var name string
		if err := row.Scan(&name); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func PrintStationNamesByLineName(db *sqlx.DB, name string) error {
	q := "SELECT s.name\nFROM line_has_stations\nLEFT JOIN stations s on line_has_stations.station_id = s.id\nWHERE line_name = $1"
	row, err := db.Query(q, name)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer row.Close()
	fmt.Printf("Stations:\n")
	for row.Next() {
		var name string
		if err := row.Scan(&name); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func main() {
	var (
		mode        string
		stationName string
		lineName    string
		dsn         string
		dataPath    string
	)

	flag.StringVar(&mode, "mode", "", "Running mode")
	flag.StringVar(&stationName, "station", "", "The name of the station")
	flag.StringVar(&lineName, "line", "", "The name of the line")
	flag.StringVar(&dsn, "dsn", "", "The DSN")
	flag.StringVar(&dataPath, "data-path", "", "The path of data json file")
	flag.Parse()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	switch mode {
	case "listLines":
		if err := PrintLineNamesByStationName(db, stationName); err != nil {
			log.Fatalf("PrintLineNamesByStationName: %s", err)
		}
	case "listStations":
		if err := PrintStationNamesByLineName(db, lineName); err != nil {
			log.Fatalf("PrintStationNamesByLineName: %s", err)
		}
	default:

		data, err := loadFromJSON(dataPath)
		if err != nil {
			log.Fatalf("Load from JSON: %s", err)
		}

		if err := writeToDB(db, data); err != nil {
			log.Fatalf("Write to database: %s", err)
		}
	}

}
