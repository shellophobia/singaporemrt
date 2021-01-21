# Overview
This assignment is built in Golang language and the http server is a binary file
that can be generated based on the platform it has to be run on.
<br />

The HTTP server will run on http://localhost:8080 which has 1 API GET /trainRoutes which will return a list of routes.
<br />
The documentation is available below for the API usage and contract

# Steps to run the server
* Set the ENV variable "STATION_MAP_PATH" in the system
    ```shell script
      export STATION_MAP_PATH=<the path to StationMap.csv file>
    ```
* Execute the file "server"
    ```shell script
      ./server
      This will start the http server on port 8080
    ```
  
# API
### GET /trainRoutes

#### Curl
```shell script
curl --location --request GET 'http://localhost:8080/trainRoutes?source=Boon%20Lay&destination=Little%20India&startTime=2019-01-31T08:00'
```

#### Request params
```json
{
    "source": "Boon Lay",
    "destination": "Little India",
    "startTime": "2019-01-31T08:00" # Optional. If not provided the routes returned won't have estimated time. The time format has to be YYYY-MM-DDTHH:mm 
}
```

#### Response
```json
{
    "source": "Boon Lay",
    "destination": "Little India",
    "suggestedRoutes": [
        {
            "stationsTravelled": 13,
            "route": ["EW27","EW26","EW25","EW24", "EW23","EW22","EW21","CC22","CC21","CC20","CC19","DT9","DT10","DT11","DT12"],
            "verboseRoute": [
                "Take EW line from Boon Lay to Lakeside",
                "Take EW line from Lakeside to Chinese Garden",
                "Take EW line from Chinese Garden to Jurong East",
                "Take EW line from Jurong East to Clementi",
                "Take EW line from Clementi to Dover",
                "Take EW line from Dover to Buona Vista",
                "Change from EW line to CC line",
                "Take CC line from Buona Vista to Holland Village",
                "Take CC line from Holland Village to Farrer Road",
                "Take CC line from Farrer Road to Botanic Gardens",
                "Change from CC line to DT line",
                "Take DT line from Botanic Gardens to Stevens",
                "Take DT line from Stevens to Newton",
                "Take DT line from Newton to Little India"
             ],
            "estimatedTimeInMinutes": 150,
            "shortestRoute": true // This is determine based on estimated time if startTime param is provided in api request else it will be based on number of stations
        },
        // .... other routes
    ]
}
```
<br />

### Code structure
#### Handlers
This package serves as a controller layer which can have validations on the API request. The logic if reusable by multiple handlers can be added into "logic" package

#### Utils
This package consists of the common utility helper functions

#### Common
This package has the common types shared across the project

### Overview of trainRoutes logic
On package initialisation the train line graph with the stations is generated using the linked list data structure
Where a station/node has following attributes
```text
{
  Code // this param holds the station code
  Name // this param holds the station name
  OpeningDate // The opening date that is present in csv file. Although this isn't used anywhere in the logic
  PrevStation // This is a pointer to the previous station in the same train line
  NextStation // This is a pointer to next station in the same train line
  LinkedStations // This is a pointer list of station nodes that are linked to the station to change to a different line
}
```
The traversal to find the routes happens using Breadth first search.
<br /> To ensure that the discovery of path is faster, the traversal is done in both directions.
<br /> When a visited node is found in any of the traversal from the other traversal. e.g. In a backward traversal a node is found which was already traversed in forward traversal then if the route is operational it is considered eligible for a potential route

<br /> There are helper functions in the file that are used to find the estimated travel time based on the rules defined for different operating hours for the week.
<br /> The structure is defined with train line as the key and if there's a new override to be made in future for a time period that affects multiple train stations, then the train lines would have to be updated

#### Structure

```text
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
```

#### Considerations for train operating hours rules
* There is a default key at trainline level which means that all the train lines that don't have a specific rule would fall under the default rules
* There is also a default key under the trainline object which is a default timerange, meaning for that trainline what is the default behaviour if it doesn't fall any under time range check
* To denote which days is the time range applicable for, the weekdays have to be listed out in the DaysOfWeek attribute as an array e.g. ["Sunday", "Monday", ...]
* To mark if the train line is not operational in the time duration, a boolean flag IsNotOperational has been kept

#### Potential Improvement
The logic can further be optimised by caching the paths between 2 stations and using them to reduce computation time.

### Running the server without binary
To build the go server code locally and run
* https://golang.org/doc/install follow this
* For the ease of use, preferred editor is Goland IDE https://www.jetbrains.com/go/download/#section=mac
* If you're using the Goland IDE you would just need to go the main.go file and run it by using the play button next to `func main` which will start the server.
* Otherwise you can execute
    ```shell script
      go run main.go
    ```
  from terminal
* For running test cases, if you're using Goland IDE, you would be able to execute the tests in get_routes_test.go by using the play button next to test case name.
* Else you can follow https://golangcode.com/run-one-test/ to execute the test file
<br />

#### Command used to generate binary file
```shell script
   env GOOS=linux GOARCH=amd64 go build -i -o server .
```
This needs to be run from the parent package that has main.go file
