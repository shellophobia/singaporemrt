package getroutes

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"gitlab.myteksi.net/goscripts/zendesk/utils"

	"github.com/thoas/go-funk"

	"github.com/gorilla/schema"
	"gitlab.myteksi.net/goscripts/zendesk/common"
)

const (
	INVALID_PREV_STATION_NUMBER = int64(-1)            // This is used in finding the closest station while inserting a new station
	INVALID_NEXT_STATION_NUMBER = int64(math.MaxInt64) // This is used in finding the closest station while inserting a new station
	QUERY_TIME_FORMAT           = "2006-01-02T15:04"   // This is the expected format in which startTime parameter in getQueryRoutes is expected
	DEFAULT_KEY                 = "default"            // For the TrainLineTimeExceptionRules map, for default values this will be the key
)

var decoder *schema.Decoder

// init function is automatically executed on package load
func init() {
	buildTrainLineMap()
	decoder = schema.NewDecoder()
}

type IHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type handler struct{}

func NewHandlerImpl() IHandler {
	return &handler{}
}

// Handle method would return the response to be returned to the API
func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Retrieve params from query params
	routeRequest := &common.GetRoutesRequest{}
	err := decoder.Decode(routeRequest, r.URL.Query())
	if err != nil {
		writeErrorResponse(err, w, 500)
	}
	err = validateRequest(routeRequest)
	if err != nil {
		writeErrorResponse(err, w, 400)
		return
	}
	routes, err := fetchRoutes(routeRequest)
	if err != nil {
		writeErrorResponse(err, w, 500)
		return
	}
	routeResponse, err := generateRouteResponse(routes, routeRequest)
	if err != nil {
		writeErrorResponse(err, w, 500)
		return
	}
	writeSuccessResponse(w, 200, routeResponse)
}

// validates the request such that only startTime is optional
func validateRequest(req *common.GetRoutesRequest) error {
	// validate start time if present
	if req.StartTime != "" {
		_, err := time.Parse(QUERY_TIME_FORMAT, req.StartTime)
		if err != nil {
			return fmt.Errorf("invalid start time")
		}
	}
	if _, ok := stationNameCodeMap[req.Source]; !ok {
		return fmt.Errorf("invalid source station")
	}
	if _, ok := stationNameCodeMap[req.Destination]; !ok {
		return fmt.Errorf("invalid destination station")
	}
	return nil
}

// Note: the write error and success response can be moved into a helper function to build a framework like setup for http server
func writeErrorResponse(err error, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResp := &common.ErrorResponse{
		Code:    statusCode,
		Message: err.Error(),
	}
	_ = json.NewEncoder(w).Encode(errResp)
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// Method to generate the route response
func generateRouteResponse(routes map[string]*common.RouteNode, req *common.GetRoutesRequest) (*common.GetRoutesResponse, error) {
	var suggestedRoutes []*common.SuggestedRoute
	shortestPathStations := int64(math.MaxInt64)
	shortestPathTimeTaken := int64(math.MaxInt64)
	for _, routeNode := range routes {

		stationPath := generateStationList(routeNode)

		var verboseRoutes, routeStations []string
		for idx, station := range stationPath {
			routeStations = append(routeStations, station.Code)
			if idx+1 != len(stationPath) { // Skip the route generate for last node as it would be covered with previous node's verboseRoute
				verboseRoute, err := generateVerboseRoute(station, stationPath[idx+1])
				if err != nil {
					return nil, err
				}
				verboseRoutes = append(verboseRoutes, verboseRoute)
			}
		}

		// Update the shortest path values
		if routeNode.StationCount < shortestPathStations {
			shortestPathStations = routeNode.StationCount
		}
		if routeNode.EstimatedTime < shortestPathTimeTaken {
			shortestPathTimeTaken = routeNode.EstimatedTime
		}

		suggestedRoute := &common.SuggestedRoute{
			StationsTravelled:      routeNode.StationCount,
			Route:                  routeStations,
			VerboseRoute:           verboseRoutes,
			EstimatedTimeInMinutes: routeNode.EstimatedTime,
			ShortestRoute:          false,
		}
		suggestedRoutes = append(suggestedRoutes, suggestedRoute)
	}
	// find the shortest path based on either time or number of stations and update the value in suggested routes
	for _, suggestedRoute := range suggestedRoutes {
		if req.StartTime == "" {
			if suggestedRoute.StationsTravelled == shortestPathStations {
				suggestedRoute.ShortestRoute = true
			}
		} else {
			if suggestedRoute.EstimatedTimeInMinutes == shortestPathTimeTaken {
				suggestedRoute.ShortestRoute = true
			}
		}
	}
	return &common.GetRoutesResponse{Source: req.Source, Destination: req.Destination, SuggestedRoutes: suggestedRoutes}, nil
}

func generateStationList(routeNode *common.RouteNode) []*common.Station {
	// traverse route as we have the a node in the middle so first we traverse backwards to get the
	// first node and then traverse forward from the middle node to reach the end node and create an ordered list to create the path
	stationPath := []*common.Station{routeNode.Station}
	// traverse backwards
	startNode := &common.RouteNode{}
	*startNode = *routeNode
	for {
		startNode = startNode.PrevNode
		if startNode == nil {
			break
		}
		stationPath = append(stationPath, startNode.Station)
	}
	stationPath = funk.Reverse(stationPath).([]*common.Station)
	// traverse forwards
	startNode = &common.RouteNode{}
	*startNode = *routeNode
	for {
		startNode = startNode.NextNode
		if startNode == nil {
			break
		}
		stationPath = append(stationPath, startNode.Station)
	}
	return stationPath
}

func generateVerboseRoute(startStation *common.Station, endStation *common.Station) (string, error) {
	startTrainLine, _, err := utils.GetStationMetadataFromCode(startStation.Code)
	if err != nil {
		return "", nil
	}
	endTrainLine, _, err := utils.GetStationMetadataFromCode(endStation.Code)
	if err != nil {
		return "", nil
	}
	// Not storing the format in a constant as this is the only place where it is used
	if startTrainLine == endTrainLine {
		return fmt.Sprintf("Take %s line from %s to %s", startTrainLine, stationCodeNameMap[startStation.Code], stationCodeNameMap[endStation.Code]), nil
	} else {
		return fmt.Sprintf("Change from %s line to %s line", startTrainLine, endTrainLine), nil
	}
}
