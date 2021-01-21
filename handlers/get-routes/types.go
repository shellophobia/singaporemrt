package getroutes

import "gitlab.myteksi.net/goscripts/zendesk/common"

var stationNameCodeMap = map[string][]string{} // Key is station name and value is a list of station codes mapped to it
var stationCodeNameMap = map[string]string{} // Reverse map of stationNameCodeMap. Key is station code and value is station name

// trainLine type would be structured as
/*
{
	<TrainLineCode>: {
		<TimeRange>: {
			NextStationTimeInMinutes: // This will give the estimated time to get to the next station
			LineChangeTimeInMinutes: // This will give the estimated time to change the line on the same station
			DaysOfWeek: // This is an array of the days of the week to which the time range config applies
			IsNotOperational: // Boolean to denote whether the line is operational in the time range. Since golang's default value for boolean is false
							  // The value is only specified if line is not operational
		}
	}
}
 */
var trainLine = map[string]map[int64]*common.Station{}

// Metadata for train line
type trainLineMeta struct {
	NextStationTimeInMinutes int64
	LineChangeTimeInMinutes  int64
	IsNotOperational         bool
	DaysOfWeek               []string
}

// TimeExceptionRule would have the rule that will be configurable to assist in determining the best route based on the period of day and time taken
type timeExceptionRule map[string]map[string]*trainLineMeta
