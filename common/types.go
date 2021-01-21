package common

// GetRoutesRequest has the expected parameters for GetRoutes request
type GetRoutesRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	StartTime   string `json:"startTime"` // Optional
}

// Route has the suggested route with the metadata about route
type SuggestedRoute struct {
	StationsTravelled      int64    `json:"stationsTravelled"`
	Route                  []string `json:"route"`
	VerboseRoute           []string `json:"verboseRoute"`
	EstimatedTimeInMinutes int64    `json:"estimatedTimeInMinutes"`
	ShortestRoute          bool     `json:"shortestRoute"` // This will denote whether it's the shortest route
}

// GetRoutesResponse has the response for get route request
type GetRoutesResponse struct {
	Source          string            `json:"source"`
	Destination     string            `json:"destination"`
	SuggestedRoutes []*SuggestedRoute `json:"suggestedRoutes"`
}

// ErrorResponse
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RouteNode is used for internal operation to fetch the routes from source to destination
type RouteNode struct {
	Station          *Station
	PrevNode         *RouteNode
	NextNode         *RouteNode
	StationCount     int64
	EstimatedTime    int64
	IsNotOperational bool
}

// Station is the
type Station struct {
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	OpeningDate    string     `json:"openingDate"`
	LinkedStations []*Station `json:"linkedStations"`
	NextStation    *Station   `json:"nextStation"`
	PrevStation    *Station   `json:"prevStation"`
}
