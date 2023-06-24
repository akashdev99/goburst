package profiler

import (
	"context"
	"fmt"
	"goburst/pkg/timeseries"
	"time"

	"github.com/struCoder/pidusage"
)

// //////////////////////////////////////////////////////////////
// Process Stats
// //////////////////////////////////////////////////////////////

type ProcessStats struct {
	Cpu          timeseries.TimeSeries
	Mem          timeseries.TimeSeries
	ReqProccesed timeseries.TimeSeries
	Pid          int
}

func NewProcessStats(pid int) *ProcessStats {
	return &ProcessStats{
		Cpu:          *timeseries.NewTimeSeries(),
		Mem:          *timeseries.NewTimeSeries(),
		ReqProccesed: *timeseries.NewTimeSeries(),
		Pid:          pid,
	}
}

// //////////////////////////////////////////////////////////////
// Stats Profiler
// //////////////////////////////////////////////////////////////

type StatsProfiler struct {
	MonitoredProcessList []int
	ProcessStats         chan map[int]*ProcessStats
	Err                  chan error
	ctx                  context.Context
	cancelFn             context.CancelFunc
}

func NewStatsProfiler(pidList []int) *StatsProfiler {
	ctx, cancelFn := context.WithCancel(context.Background())

	return &StatsProfiler{
		MonitoredProcessList: pidList,
		ProcessStats:         make(chan map[int]*ProcessStats, 1),
		Err:                  make(chan error, 1),
		ctx:                  ctx,
		cancelFn:             cancelFn,
	}
}

func (profiler *StatsProfiler) CpuMemProfiler(intervalSec int) {
	processStats := make(map[int]*ProcessStats)
	for _, pid := range profiler.MonitoredProcessList {
		processStats[pid] = NewProcessStats(pid)
	}

	var cpuUsage, rssBytes float64

	for {
		select {
		case <-profiler.ctx.Done():
			fmt.Println("Profiling Done")
			profiler.ProcessStats <- processStats
			return
		default:
			for _, pid := range profiler.MonitoredProcessList {
				sysInfo, err := pidusage.GetStat(pid)
				if err != nil {
					profiler.Err <- err
					return
				}
				cpuUsage, rssBytes = sysInfo.CPU, sysInfo.Memory

				processStats[pid].Cpu.Add(int(cpuUsage))
				processStats[pid].Mem.Add(int(rssBytes))
			}
			time.Sleep(time.Millisecond * time.Duration(intervalSec))

		}
	}
}

func (profiler *StatsProfiler) IsSystemStatsProfiling() bool {
	return len(profiler.MonitoredProcessList) != 0
}

func (profiler *StatsProfiler) Stop() {
	profiler.cancelFn()
}
