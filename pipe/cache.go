package pipe

import "sync"

const (
	CACHE_SIZE = 500
)

var LinesChan = make(chan Line, CACHE_SIZE)
var lineCache = struct {
	sync.RWMutex
	lines []Line
}{sync.RWMutex{}, []Line{}}

func addCacheLine(line Line) {
	LinesChan <- line
	lineCache.Lock()
	defer lineCache.Unlock()
	if len(lineCache.lines) >= CACHE_SIZE {
		lineCache.lines = lineCache.lines[:CACHE_SIZE]
	}
	lineCache.lines = append([]Line{line}, lineCache.lines...)
}

func GetLineCache() []Line {
	lineCache.RLock()
	defer lineCache.RUnlock()
	return lineCache.lines
}
