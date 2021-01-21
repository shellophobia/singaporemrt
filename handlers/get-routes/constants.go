package getroutes

// TrainLineTimeExceptionRules would have the list of rules that are configurable to assist in determining the best route based on the period of day and time taken
/*
The structure is
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
var TrainLineTimeExceptionRules = timeExceptionRule{
	"NS": {
		"6:00AM - 9:00AM": {
			NextStationTimeInMinutes: 12,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"6:00PM - 9:00PM": {
			NextStationTimeInMinutes: 12,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
	},
	"NE": {
		"6:00AM - 9:00AM": {
			NextStationTimeInMinutes: 12,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"6:00PM - 9:00PM": {
			NextStationTimeInMinutes: 12,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
	},
	"DT": {
		"10:00PM - 11:59PM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			LineChangeTimeInMinutes: 10,
		},
		"12:00AM - 6:00AM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			LineChangeTimeInMinutes: 10,
		},
		"6:00AM - 9:00AM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"6:00PM - 9:00PM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"default": {
			NextStationTimeInMinutes: 8,
			LineChangeTimeInMinutes:  10,
		},
	},
	"CG": {
		"10:00PM - 11:59PM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			LineChangeTimeInMinutes: 10,
		},
		"12:00AM - 6:00AM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			LineChangeTimeInMinutes: 10,
		},
	},
	"CE": {
		"10:00PM - 11:59PM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			LineChangeTimeInMinutes: 10,
		},
		"12:00AM - 6:00AM": {
			IsNotOperational:        true,
			DaysOfWeek:              []string{"Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			LineChangeTimeInMinutes: 10,
		},
	},
	"TE": {
		"10:00PM - 11:59PM": {
			NextStationTimeInMinutes: 8,
			LineChangeTimeInMinutes:  10,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		},
		"12:00AM - 6:00AM": {
			NextStationTimeInMinutes: 8,
			LineChangeTimeInMinutes:  10,
			DaysOfWeek:               []string{"Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		},
		"6:00AM - 9:00AM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"6:00PM - 9:00PM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"default": {
			NextStationTimeInMinutes: 8,
			LineChangeTimeInMinutes:  10,
		},
	},
	"default": {
		"6:00AM - 9:00AM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"6:00PM - 9:00PM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  15,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		"10:00PM - 11:59PM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  10,
			DaysOfWeek:               []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		},
		"12:00AM - 6:00AM": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  10,
			DaysOfWeek:               []string{"Tuesday", "Wednesday", "Thursday", "Fri", "Friday", "Sunday"},
		},
		"default": {
			NextStationTimeInMinutes: 10,
			LineChangeTimeInMinutes:  10,
		},
	},
}
