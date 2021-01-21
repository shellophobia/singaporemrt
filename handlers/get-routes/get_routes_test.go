package getroutes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.myteksi.net/goscripts/zendesk/common"
)

func TestGetRouteEstimate(t *testing.T) {
	t.Run("gets the route estimate", func(t *testing.T) {
		testCases := []struct {
			sourceStation      string
			destinationStation string
			queryTime          string
			stationCount       int64
			estimatedTime      int64
			isNotOperational   bool
			err                error
		}{
			{
				// peak hours
				sourceStation:      "NS1",
				destinationStation: "NS2",
				queryTime:          "2019-01-31T08:00",
				stationCount:       1,
				estimatedTime:      12,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// peak hours
				sourceStation:      "NS1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T08:00",
				stationCount:       0,
				estimatedTime:      15,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// peak hours
				sourceStation:      "CC1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T08:00",
				stationCount:       0,
				estimatedTime:      15,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// peak hours
				sourceStation:      "DT1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T08:00",
				stationCount:       1,
				estimatedTime:      10,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// night hours
				sourceStation:      "DT1",
				destinationStation: "NS2",
				queryTime:          "2019-01-31T01:00",
				stationCount:       0,
				estimatedTime:      10,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// night hours
				sourceStation:      "DT1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T01:00",
				stationCount:       1,
				estimatedTime:      0,
				isNotOperational:   true,
			},
			{
				// night hours
				sourceStation:      "CC1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T01:00",
				stationCount:       0,
				estimatedTime:      10,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// night hours
				sourceStation:      "TE1",
				destinationStation: "TE2",
				queryTime:          "2019-01-31T01:00",
				stationCount:       1,
				estimatedTime:      8,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// non-peak hours
				sourceStation:      "DT1",
				destinationStation: "NS2",
				queryTime:          "2019-01-31T17:00",
				stationCount:       0,
				estimatedTime:      10,
				isNotOperational:   false,
				err:                nil,
			},
			{
				// non-peak hours
				sourceStation:      "DT1",
				destinationStation: "DT2",
				queryTime:          "2019-01-31T17:00",
				stationCount:       1,
				estimatedTime:      8,
				isNotOperational:   false,
			},
			{
				// non-peak hours
				sourceStation:      "CC1",
				destinationStation: "CG2",
				queryTime:          "2019-01-31T17:00",
				stationCount:       0,
				estimatedTime:      10,
				isNotOperational:   false,
				err:                nil,
			},
			// TODO: Add a test case for error check
		}
		for _, testCase := range testCases {
			stationCount, estimatedTime, isNotOperational, err := getRouteEstimate(testCase.sourceStation, testCase.destinationStation, testCase.queryTime)
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.stationCount, stationCount)
			assert.Equal(t, testCase.estimatedTime, estimatedTime)
			assert.Equal(t, testCase.isNotOperational, isNotOperational)
		}
	})
}

func TestFetchRoutes(t *testing.T) {
	// Test cases can be more enhanced to check the route changes
	t.Run("fetches the routes with estimated time", func(t *testing.T) {
		req := &common.GetRoutesRequest{
			Source:      "Marsiling",
			Destination: "Yio Chu Kang",
			StartTime:   "2019-01-31T17:00",
		}
		routes, err := fetchRoutes(req)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(routes))
		for code, route := range routes {
			if code == "NS12" {
				stationPath := generateStationList(route)
				assert.Equal(t, 8, len(stationPath))
				assert.Equal(t, "NS8", stationPath[0].Code)
				assert.Equal(t, "NS9", stationPath[1].Code)
			} else if code == "TE8" {
				stationPath := generateStationList(route)
				assert.Equal(t, 16, len(stationPath))
			} else if code == "EW21" {
				stationPath := generateStationList(route)
				assert.Equal(t, 21, len(stationPath))
			}
		}
	})

	t.Run("returns no routes on a non-operational path", func(t *testing.T) {
		req := &common.GetRoutesRequest{
			Source:      "Bencoolen",
			Destination: "Ubi",
			StartTime:   "2019-01-31T01:00",
		}
		routes, err := fetchRoutes(req)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(routes))
	})
}
