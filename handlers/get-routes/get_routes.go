package getroutes

import (
	"fmt"
	"strings"
	"time"

	"gitlab.myteksi.net/goscripts/zendesk/common"

	"github.com/thoas/go-funk"
	"gitlab.myteksi.net/goscripts/zendesk/utils"
)

func fetchRoutes(req *common.GetRoutesRequest) (map[string]*common.RouteNode, error) {
	sourceStationNodes := stationNameCodeMap[req.Source]
	destinationStationNodes := stationNameCodeMap[req.Destination]
	var routeNodeListForwardTraversal, routeNodeListBackwardTraversal []*common.RouteNode
	pathNodes := map[string]*common.RouteNode{}
	visitedRouteNodesForwardTraversal := map[string]*common.RouteNode{}
	visitedRouteNodesBackwardTraversal := map[string]*common.RouteNode{}
	// Store all the stations on different lines for the source station in the first level of breadth first traversal
	for _, sourceStationNode := range sourceStationNodes {
		lineName, stNumber, err := utils.GetStationMetadataFromCode(sourceStationNode)
		if err != nil {
			return nil, err
		}
		routeNode := &common.RouteNode{Station: trainLine[lineName][stNumber]}
		routeNodeListForwardTraversal = append(routeNodeListForwardTraversal, routeNode)
		visitedRouteNodesForwardTraversal[trainLine[lineName][stNumber].Code] = routeNode
	}
	// Store all the stations on different lines for the destination station in the first level of breadth first traversal
	for _, destinationStationNode := range destinationStationNodes {
		lineName, stNumber, err := utils.GetStationMetadataFromCode(destinationStationNode)
		if err != nil {
			return nil, err
		}
		routeNode := &common.RouteNode{Station: trainLine[lineName][stNumber]}
		routeNodeListBackwardTraversal = append(routeNodeListBackwardTraversal, routeNode)
		visitedRouteNodesBackwardTraversal[trainLine[lineName][stNumber].Code] = routeNode
	}
	for {
		if len(routeNodeListBackwardTraversal) == 0 && len(routeNodeListForwardTraversal) == 0 {
			// If all the nodes have been visited stop searching
			break
		}

		// Breadth first search in forward and backward direction to find the routes faster

		// forward traversal of stations
		var tempStationListForwardTraversal, tempStationListBackwardTraversal []*common.RouteNode
		for _, routeNode := range routeNodeListForwardTraversal {
			if routeNode.IsNotOperational {
				// If node was marked as not operational in last iteration, terminate search for this path
				continue
			}
			// If the node has been visited before it means that it is a potential valid route
			if _, ok := visitedRouteNodesBackwardTraversal[routeNode.Station.Code]; ok {
				if _, ok := pathNodes[routeNode.Station.Code]; ok {
					// If the node has already been added in the path route, ignore
					continue
				}
				if visitedRouteNodesBackwardTraversal[routeNode.Station.Code].IsNotOperational {
					continue // Do not merge with a non-operational path
				}
				pathNodes[routeNode.Station.Code] = routeNode
				// connect the nodes
				routeNode.NextNode = visitedRouteNodesBackwardTraversal[routeNode.Station.Code].NextNode
				routeNode.StationCount += visitedRouteNodesBackwardTraversal[routeNode.Station.Code].StationCount + 1 // add the station count, +1 is to count the first station
				routeNode.EstimatedTime += visitedRouteNodesBackwardTraversal[routeNode.Station.Code].EstimatedTime   // add the time taken
				if visitedRouteNodesBackwardTraversal[routeNode.Station.Code].NextNode != nil {
					visitedRouteNodesBackwardTraversal[routeNode.Station.Code].NextNode.PrevNode = routeNode
				}
				continue
			}
			// Populate next station
			err := populateTempTraversalList(visitedRouteNodesForwardTraversal, routeNode, routeNode.Station.NextStation, &tempStationListForwardTraversal, req.StartTime, true)
			if err != nil {
				return nil, err
			}

			// Populate previous station
			err = populateTempTraversalList(visitedRouteNodesForwardTraversal, routeNode, routeNode.Station.PrevStation, &tempStationListForwardTraversal, req.StartTime, true)
			if err != nil {
				return nil, err
			}

			// iterate over linked stations
			for _, linkedStation := range routeNode.Station.LinkedStations {
				err = populateTempTraversalList(visitedRouteNodesForwardTraversal, routeNode, linkedStation, &tempStationListForwardTraversal, req.StartTime, true)
				if err != nil {
					return nil, err
				}
			}
		}

		// start back traversal
		for _, routeNode := range routeNodeListBackwardTraversal {
			if routeNode.IsNotOperational {
				// If node was marked as not operational in last iteration, terminate search for this path
				continue
			}
			// If the node has been visited before it means that it is a potential valid route
			if _, ok := visitedRouteNodesForwardTraversal[routeNode.Station.Code]; ok {
				if _, ok := pathNodes[routeNode.Station.Code]; ok {
					// If node was marked as not operational in last iteration, terminate search for this path
					continue
				}
				if visitedRouteNodesForwardTraversal[routeNode.Station.Code].IsNotOperational {
					continue // Do not merge with a non-operational path
				}
				pathNodes[routeNode.Station.Code] = routeNode
				// connect the nodes
				routeNode.PrevNode = visitedRouteNodesForwardTraversal[routeNode.Station.Code].PrevNode
				routeNode.StationCount += visitedRouteNodesForwardTraversal[routeNode.Station.Code].StationCount   // add the station count
				routeNode.EstimatedTime += visitedRouteNodesForwardTraversal[routeNode.Station.Code].EstimatedTime // add the time taken
				if visitedRouteNodesForwardTraversal[routeNode.Station.Code].PrevNode != nil {
					visitedRouteNodesForwardTraversal[routeNode.Station.Code].PrevNode.NextNode = routeNode
				}
				continue
			}
			// Populate next station
			err := populateTempTraversalList(visitedRouteNodesBackwardTraversal, routeNode, routeNode.Station.NextStation, &tempStationListBackwardTraversal, req.StartTime, false)
			if err != nil {
				return nil, err
			}

			// Populate previous station
			err = populateTempTraversalList(visitedRouteNodesBackwardTraversal, routeNode, routeNode.Station.PrevStation, &tempStationListBackwardTraversal, req.StartTime, false)
			if err != nil {
				return nil, err
			}

			// iterate over linked stations
			for _, linkedStation := range routeNode.Station.LinkedStations {
				err = populateTempTraversalList(visitedRouteNodesBackwardTraversal, routeNode, linkedStation, &tempStationListBackwardTraversal, req.StartTime, false)
				if err != nil {
					return nil, err
				}
			}
		}
		routeNodeListForwardTraversal = tempStationListForwardTraversal
		routeNodeListBackwardTraversal = tempStationListBackwardTraversal
	}
	return pathNodes, nil
}

// This method will populate the nodes which will be used in the next traversal
func populateTempTraversalList(visitedRouteNodes map[string]*common.RouteNode, currentRouteNode *common.RouteNode, nextStation *common.Station, tempTraversalList *[]*common.RouteNode,
	queryTimeString string, forwardTraversal bool) error {
	if nextStation == nil {
		return nil
	}
	nextRouteNode := &common.RouteNode{Station: nextStation}
	if _, ok := visitedRouteNodes[nextStation.Code]; !ok {
		// mark as visited
		visitedRouteNodes[nextRouteNode.Station.Code] = nextRouteNode
		stationCount, estimatedTime, isNotOperational, err := getRouteEstimate(currentRouteNode.Station.Code, nextRouteNode.Station.Code, queryTimeString)
		if err != nil {
			return err
		}
		nextRouteNode.StationCount = currentRouteNode.StationCount + stationCount
		nextRouteNode.EstimatedTime = currentRouteNode.EstimatedTime + estimatedTime
		if forwardTraversal {
			nextRouteNode.PrevNode = currentRouteNode
		} else {
			nextRouteNode.NextNode = currentRouteNode
		}
		nextRouteNode.IsNotOperational = isNotOperational
		*tempTraversalList = append(*tempTraversalList, nextRouteNode)
	}
	return nil
}

func getRouteEstimate(startStationCode, endStationCode string, queryTimeString string) (int64, int64, bool, error) {
	startLineName, _, err := utils.GetStationMetadataFromCode(startStationCode)
	if err != nil {
		return 0, 0, false, err
	}
	endLineName, _, err := utils.GetStationMetadataFromCode(endStationCode)
	if err != nil {
		return 0, 0, false, err
	}
	var stationCount int64
	if startLineName == endLineName {
		stationCount = 1 // If both stations are on same line count the station
	}
	if queryTimeString == "" {
		return stationCount, 0, false, nil // Skip processing if query time string wasn't provided
	}
	var lineTimeConfig map[string]*trainLineMeta
	var ok bool
	lineTimeConfig, ok = TrainLineTimeExceptionRules[startLineName]
	if !ok {
		// Assuming default is always there
		lineTimeConfig, ok = TrainLineTimeExceptionRules[DEFAULT_KEY]
	}
	if lineTimeConfig == nil || !ok {
		return 0, 0, false, fmt.Errorf("missing train line config")
	}
	var eligibleTrainLineMeta *trainLineMeta
	// Iterate through all time ranges
	for timeRange, trainLineMeta := range lineTimeConfig {
		if timeRange == DEFAULT_KEY {
			// Skip time range parsing if default values are present
			eligibleTrainLineMeta = trainLineMeta
			continue
		}
		isTimeConfigApplicable, err := isTimeConfigApplicable(timeRange, queryTimeString, trainLineMeta)
		if err != nil {
			return 0, 0, false, err
		}
		if isTimeConfigApplicable {
			eligibleTrainLineMeta = trainLineMeta
			break
		}
	}
	if eligibleTrainLineMeta == nil {
		// lookup for default station config in default time
		if lineTimeConfig, ok := TrainLineTimeExceptionRules[DEFAULT_KEY]; ok {
			if eligibleTrainLineMeta, ok = lineTimeConfig[DEFAULT_KEY]; !ok {
				return 0, 0, false, fmt.Errorf("missing train line config")
			}
		} else {
			return 0, 0, false, fmt.Errorf("missing train line config")
		}
	}
	estimedTimeInMinutes, isNotOperational := getEstimatedTimeFromTrainLineMeta(eligibleTrainLineMeta, startLineName == endLineName)
	return stationCount, estimedTimeInMinutes, isNotOperational, nil
}

func isTimeConfigApplicable(timeRange string, queryTimeString string, trainLineMeta *trainLineMeta) (bool, error) {
	// parse time range
	timeStrings := strings.Split(timeRange, " - ")
	if len(timeStrings) != 2 {
		return false, fmt.Errorf("invalid time range in config")
	}
	// parse query time string
	queryTime, err := time.Parse(QUERY_TIME_FORMAT, queryTimeString)
	if err != nil {
		return false, err
	}
	// Get the start time equivalent for query time
	startTime, err := time.Parse(time.Kitchen, timeStrings[0])
	if err != nil {
		return false, err
	}
	startTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	// Get the end time equivalent for query time
	endTime, err := time.Parse(time.Kitchen, timeStrings[1])
	if err != nil {
		return false, err
	}
	endTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
	if startTime.After(queryTime.UTC()) || endTime.Before(queryTime.UTC()) {
		return false, nil
	}

	// Proceed to validate further
	if funk.ContainsString(trainLineMeta.DaysOfWeek, queryTime.Weekday().String()) {
		return true, nil
	}
	return false, nil
}

func getEstimatedTimeFromTrainLineMeta(trainLineMeta *trainLineMeta, sameLine bool) (int64, bool) {
	if trainLineMeta.IsNotOperational && sameLine {
		return 0, true
	}
	if sameLine {
		return trainLineMeta.NextStationTimeInMinutes, false
	}
	return trainLineMeta.LineChangeTimeInMinutes, false
}
