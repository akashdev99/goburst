package profiler

import (
	"goburst/pkg/cutefmt"
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
		cutefmt.Errorf("Stopping API request burst due to error : %v \n", err)
		return nil, err
	case <-profiler.apiBurster.Done:
		cutefmt.Successf("Api request burst Complete\n")
		profiler.statsProfiler.Stop()
	}

	if profiler.statsProfiler.IsSystemStatsProfiling() {
		select {
		case err := <-profiler.statsProfiler.Err:
			cutefmt.Errorf("Stopping System stats profiling due to error : %v \n", err)
			return nil, err
		case profilerData = <-profiler.statsProfiler.ProcessStats:
			cutefmt.Successf("System stats profiling Complete\n")
			return profilerData, nil
		}
	}
	return nil, nil
}
