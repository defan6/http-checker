package checker

import "time"

type Result struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}
