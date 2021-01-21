package getroutes

import (
	"encoding/csv"
	"gitlab.myteksi.net/goscripts/zendesk/common"
	"gitlab.myteksi.net/goscripts/zendesk/utils"
	"io"
	"log"
	"os"
)

// Builds the cache for querying the path between stations
func buildTrainLineMap() {
	stationMapPath := os.Getenv("STATION_MAP_PATH")
	if stationMapPath == "" {
		log.Fatalln("STATION_MAP_PATH Env variable not defined")
	}
	csvfile, err := os.Open(os.Getenv("STATION_MAP_PATH"))
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	reader := csv.NewReader(csvfile)
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error in reading the file row : %v\n", err.Error())
		}
		rowCount++
		if rowCount == 1 {
			continue // Skip processing the header of csv
		}
		station := &common.Station{
			Code:        record[0],
			Name:        record[1],
			OpeningDate: record[2],
		}

		// store the map of station name to a list of station codes
		if _, ok := stationNameCodeMap[station.Name]; ok {
			// Link found
			// Update the links of all mapped stations
			for _, stationCode := range stationNameCodeMap[station.Name] {
				if stationCode == station.Code {
					continue // avoid duplicate insertions
				}
				// Train line should have the mapped station
				lineCode, stNumber, err := utils.GetStationMetadataFromCode(stationCode)
				if err != nil {
					log.Fatalln("error in retrieving station metadata from code")
				}
				linkedStation := trainLine[lineCode][stNumber]
				linkedStation.LinkedStations = append(linkedStation.LinkedStations, station)
				station.LinkedStations = append(station.LinkedStations, linkedStation) // Link to new station
			}
			stationNameCodeMap[station.Name] = append(stationNameCodeMap[station.Name], station.Code)
		} else {
			stationNameCodeMap[station.Name] = []string{station.Code}
		}

		// store the code to station name mapping
		stationCodeNameMap[station.Code] = station.Name

		// Build train station graph which will used for calculating routes which is using linked list data structure
		lineCode, stNumber, err := utils.GetStationMetadataFromCode(station.Code)
		if err != nil {
			log.Fatalln("error in retrieving station metadata from code")
		}
		if _, ok := trainLine[lineCode]; ok {
			// find the closest linked nodes to insert new station
			prevStationNumber := INVALID_PREV_STATION_NUMBER
			nextStationNumber := INVALID_NEXT_STATION_NUMBER
			for number, _ := range trainLine[lineCode] {
				if number > prevStationNumber && number < stNumber {
					prevStationNumber = number
				}
				if number < nextStationNumber && number > stNumber {
					nextStationNumber = number
				}
			}
			if prevStationNumber != INVALID_PREV_STATION_NUMBER {
				// Insert new station
				nextStation := trainLine[lineCode][prevStationNumber].NextStation
				trainLine[lineCode][prevStationNumber].NextStation = station
				station.PrevStation = trainLine[lineCode][prevStationNumber]
				station.NextStation = nextStation
				if nextStation != nil {
					nextStation.PrevStation = station
				}
			}
			if station.PrevStation == nil && nextStationNumber != INVALID_NEXT_STATION_NUMBER {
				// new station is the 1st node
				trainLine[lineCode][nextStationNumber].PrevStation = station
				station.NextStation = trainLine[lineCode][nextStationNumber]
			}
			trainLine[lineCode][stNumber] = station
		} else {
			trainLine[lineCode] = map[int64]*common.Station{
				stNumber: station,
			}
		}
	}
}
