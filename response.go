package scrap

import (
	"net/http"
	"time"
)

// Represents a response from the server, containing the original
// http.Response object, extra contextual/statistic data, and various
// convenience functions.
type ServerResponse struct {
	http.Response
	Stats ResponseStats
}

// Some statistics on the completed request, for example, the time
// required to retrieve the file from the network.
type ResponseStats struct {
	Start    time.Time
	Duration time.Duration
}
