package pipe

import "sync"

const (
	CACHE_SIZE = 500
)

var LinesChan = make(chan Line)
var lineCache = struct {
	sync.RWMutex
	lines []Line
}{sync.RWMutex{}, []Line{}}

func addCacheLine(line Line) {
	go func() { LinesChan <- line }()
	lineCache.Lock()
	defer lineCache.Unlock()
	if len(lineCache.lines) >= CACHE_SIZE {
		lineCache.lines = lineCache.lines[1:]
	}
	lineCache.lines = append(lineCache.lines, line)
}

func GetLineCache() []Line {
	lineCache.RLock()
	defer lineCache.RUnlock()
	return lineCache.lines
}
