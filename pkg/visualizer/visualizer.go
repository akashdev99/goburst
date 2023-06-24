package visualizer

import (
	"fmt"
	"goburst/internal/profiler"
	"goburst/pkg/timeseries"
	"os"
	"strconv"

	"time"

	ps "github.com/mitchellh/go-ps"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func Visualize(title string, processStatMap map[int]*profiler.ProcessStats) {
	f, _ := os.Create(title + ".html")
	for pid, processStats := range processStatMap {
		MakeLineChart(title+"-CPU", "CPU", makeDescription(pid, processStats), processStats.Cpu, f)
		MakeLineChart(title+"-MEM", "MEM(bytes)", makeDescription(pid, processStats), processStats.Mem, f)
	}
}

func makeDescription(pid int, processStats *profiler.ProcessStats) string {
	timeTaken := processStats.Cpu.EndTime() - processStats.Cpu.StartTime()

	var process string
	process, err := getProcessName(pid)
	if err != nil {
		process = strconv.Itoa(pid)
	}
	return fmt.Sprintf("Process \"%v\" took %d seconds to complete", process, timeTaken)
}

func getProcessName(pid int) (string, error) {
	p, err := ps.FindProcess(pid)
	if err != nil {
		fmt.Printf("Failed to find process name for pid %v : error - %v \n", pid, err)
		return "", err
	}

	return p.Executable(), nil
}

func MakeLineChart(title string, XLabel string, description string, stats timeseries.TimeSeries, f *os.File) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: description,
		}))

	xAxisTime := []string{}
	for _, unixTime := range stats.GetTimeList() {
		t := time.Unix(int64(unixTime), 0)
		xAxisTime = append(xAxisTime, t.Format("15:04:05.00000"))
	}

	// Put data into instance
	line.SetXAxis(xAxisTime).
		AddSeries(XLabel, generateLineItems(stats)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	line.Render(f)
	fmt.Println("Succesfully Created Visualization for " + title)
}

func generateLineItems(timeSeries timeseries.TimeSeries) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < timeSeries.Size(); i++ {
		items = append(items, opts.LineData{Value: timeSeries.GetValue(i)})
	}
	return items
}
