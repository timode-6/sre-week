package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Span struct {
	host     string
	id       int
	parentID int
	start    int
	end      int
	cpu      int
}

type Event struct {
	time  int
	delta int
	span  *Span
}

type Interval struct {
	host       string
	start      int
	end        int
	cpu        int
	topSpanID  int
	topSpanCPU int
}

func parseInput(filename string) ([]Span, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var spans []Span
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 6 {
			return nil, fmt.Errorf("invalid input format")
		}
		id, _ := strconv.Atoi(fields[1])
		parentID, _ := strconv.Atoi(fields[2])
		start, _ := strconv.Atoi(fields[3])
		end, _ := strconv.Atoi(fields[4])
		cpu, _ := strconv.Atoi(fields[5])
		spans = append(spans, Span{fields[0], id, parentID, start, end, cpu})
	}
	return spans, scanner.Err()
}

func findRootSpan(span Span, spans []Span) int {
	if span.parentID == -1 {
		return span.id
	}
	for _, s := range spans {
		if s.id == span.parentID {
			return findRootSpan(s, spans)
		}
	}
	return span.id
}

func findIntervals(spans []Span, threshold int, isHost bool) []Interval {
	var events []Event
	for _, span := range spans {
		events = append(events, Event{span.start, span.cpu, &span})
		events = append(events, Event{span.end + 1, -span.cpu, &span})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].time < events[j].time
	})

	var intervals []Interval
	var currentInterval *Interval
	currentCPU := 0
	spanCPU := make(map[int]int)

	for i, event := range events {
		currentCPU += event.delta

		if currentCPU >= threshold && currentInterval == nil {
			currentInterval = &Interval{start: event.time, end: event.time, cpu: currentCPU}
			if isHost {
				currentInterval.host = spans[0].host
			}
		}

		if currentCPU < threshold && currentInterval != nil {
			currentInterval.end = event.time - 1
			intervals = append(intervals, *currentInterval)
			currentInterval = nil
		}

		if currentInterval != nil {
			currentInterval.cpu = max(currentInterval.cpu, currentCPU)
		}

		for _, span := range spans {
			if span.start <= event.time && event.time <= span.end {
				rootSpanID := findRootSpan(span, spans)
				spanCPU[rootSpanID] += span.cpu
			}
		}

		if currentInterval != nil {
			topSpanID, topSpanCPU := 0, 0
			for id, cpu := range spanCPU {
				if cpu > topSpanCPU || (cpu == topSpanCPU && id < topSpanID) {
					topSpanID, topSpanCPU = id, cpu
				}
			}
			currentInterval.topSpanID = topSpanID
			currentInterval.topSpanCPU = topSpanCPU
		}

		if i < len(events)-1 && events[i+1].time > event.time {
			for k := range spanCPU {
				delete(spanCPU, k)
			}
		}
	}

	return intervals
}

func mergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	merged := []Interval{intervals[0]}

	for _, interval := range intervals[1:] {
		last := &merged[len(merged)-1]
		if interval.start <= last.end+1 &&
			interval.cpu == last.cpu &&
			interval.topSpanID == last.topSpanID &&
			interval.topSpanCPU == last.topSpanCPU {
			last.end = max(last.end, interval.end)
		} else {
			merged = append(merged, interval)
		}
	}

	return merged
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide input file name")
		return
	}

	spans, err := parseInput(os.Args[1])
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	globalIntervals := findIntervals(spans, 80, false)
	globalIntervals = mergeIntervals(globalIntervals)

	for _, interval := range globalIntervals {
		fmt.Printf("%d %d %d %d %d\n", interval.start, interval.end, interval.cpu, interval.topSpanID, interval.topSpanCPU)
	}

	hosts := make(map[string]bool)
	for _, span := range spans {
		hosts[span.host] = true
	}

	var sortedHosts []string
	for host := range hosts {
		sortedHosts = append(sortedHosts, host)
	}
	sort.Strings(sortedHosts)

	for _, host := range sortedHosts {
		hostSpans := make([]Span, 0)
		for _, span := range spans {
			if span.host == host {
				hostSpans = append(hostSpans, span)
			}
		}
		intervals := findIntervals(hostSpans, 90, true)
		intervals = mergeIntervals(intervals)
		for _, interval := range intervals {
			fmt.Printf("%s %d %d %d %d %d\n", host, interval.start, interval.end, interval.cpu, interval.topSpanID, interval.topSpanCPU)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}