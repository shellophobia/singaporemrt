package utils

import "strconv"

// GetStationMetadataFromCode returns the train line name and station code in integer format
// Assumes that the train line is a 2 character value which is a prefix of stationcode and stationcode is always a number
// The logic can be enhanced further by either having another field for trainline or add these constraints on stationcode naming
func GetStationMetadataFromCode(stationCode string) (string, int64, error) {
	code, err := strconv.ParseInt(stationCode[2:], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return stationCode[:2], code, nil
}
