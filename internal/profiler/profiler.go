package profiler

import (
	"fmt"
)

type Profiler struct {
	TaskName      string
	statsProfiler *StatsProfiler
	apiBurster    *ApiBurster
}

func NewProfiler(task string, method string, url string, headers []string,
	pidList []int, Iteration int) *Profiler {
	return &Profiler{
		TaskName:      task,
		statsProfiler: NewStatsProfiler(pidList),
		apiBurster:    NewApiBurster(method, url, headers, Iteration),
	}
}

func (profiler *Profiler) Profile(method string, url string, headers []string, iteration int, interval int) (profilerData map[int]*ProcessStats, err error) {
	go profiler.apiBurster.BurstRequests(method, url, headers, iteration)

	if profiler.statsProfiler.IsSystemStatsProfiling() {
		go profiler.statsProfiler.CpuMemProfiler(interval)
	}

	select {
	case err := <-profiler.apiBurster.Err:
		profiler.statsProfiler.Stop()
		fmt.Printf("Stopping API request bombarding due to error : %v \n", err)
		return nil, err
	case <-profiler.apiBurster.Done:
		profiler.statsProfiler.Stop()
		fmt.Printf("\n%v number of requests sent to URL \n", iteration)
	}

	if profiler.statsProfiler.IsSystemStatsProfiling() {
		select {
		case err := <-profiler.statsProfiler.Err:
			fmt.Printf("Stopping System stats profiling due to error : %v \n", err)
			return nil, err
		case profilerData = <-profiler.statsProfiler.ProcessStats:
			fmt.Printf("System stats profiling successfull")
			return profilerData, nil
		}
	}
	return nil, nil
}
