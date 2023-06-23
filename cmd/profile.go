/*
Copyright © 2023 NAME HERE akashnandan99@gmail.com
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"goburst/pkg/timeseries"

	"github.com/spf13/cobra"
	"github.com/struCoder/pidusage"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Profie your api",
	Long:  `profile command runs the api endpoint multiple times while taking CPU and Memory Sample at a specified interval.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		headers, err := cmd.Flags().GetStringSlice("header")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		iteration, err := cmd.Flags().GetInt("iteration")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		method, err := cmd.Flags().GetString("method")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		pidList, err := cmd.Flags().GetIntSlice("pidlist")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			fmt.Println("Error parsing --url flag :", err)
		}

		profiler := NewProfiler("test", pidList, iteration)
		if err := profiler.Profile(method, url, headers, iteration, interval); err != nil {
			fmt.Println("Profiling not complete due to error ", err)
			return
		}

		fmt.Println("Profiling Complete")
	},
}

type Profiler struct {
	TaskName             string
	MonitoredProcessList []int
	Iterations           int
	Done                 chan bool
	ProfilerData         chan map[int]*ProcessStats
	Err                  chan bool
}

func NewProfiler(task string, processList []int, Iteration int) *Profiler {
	return &Profiler{
		TaskName:             task,
		MonitoredProcessList: processList,
		Iterations:           Iteration,
		Done:                 make(chan bool, 1),
		ProfilerData:         make(chan map[int]*ProcessStats, 1),
		Err:                  make(chan bool, 1),
	}
}

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

func (profiler *Profiler) Profile(method string, url string, headers []string, iteration int, interval int) error {
	profiler.BurstRequests(method, url, headers, iteration)

	if len(profiler.MonitoredProcessList) != 0 {
		profiler.CpuMemProfiler(interval)
	}
	return nil
}

func (profiler *Profiler) BurstRequests(method string, url string, headers []string, iteration int) {
	httpClient, request, err := createHttpClient(method, url, headers)
	if err != nil {
		fmt.Printf("Failed create http client : %v \n", err)
		return
	}

	startTime := time.Now().Unix()
	for i := 0; i < iteration; i++ {
		err := makeRequest(httpClient, request)
		if err != nil {
			fmt.Printf("API Failed , stopping profiling at count %v!!! : %v \n", i, err)
			profiler.Err <- true
			return
		}
	}
	endTime := time.Now().Unix()
	fmt.Printf("%v:Total Time took to complete %v request = %v second \n", profiler.TaskName, iteration, endTime-startTime)
	profiler.Done <- true
}

func (profiler *Profiler) CpuMemProfiler(intervalSec int) {
	processStats := make(map[int]*ProcessStats)
	for _, pid := range profiler.MonitoredProcessList {
		processStats[pid] = NewProcessStats(pid)
	}

	var cpuUsage, rssBytes float64

	for {
		select {
		case <-profiler.Done:
			fmt.Println("Profiling Done")
			profiler.ProfilerData <- processStats
			profiler.Err <- false
			return
		default:
			for _, pid := range profiler.MonitoredProcessList {
				sysInfo, err := pidusage.GetStat(pid)
				if err != nil {
					fmt.Println("Could not get system Cpu Info: ", err)
					profiler.Err <- true
					profiler.ProfilerData <- processStats
				}
				cpuUsage, rssBytes = sysInfo.CPU, sysInfo.Memory

				processStats[pid].Cpu.Add(int(cpuUsage))
				processStats[pid].Mem.Add(int(rssBytes))
				processStats[pid].Mem.Add(profiler.Iterations)
			}
			time.Sleep(time.Second * time.Duration(intervalSec))
		}
	}
}

func createHttpClient(method string, url string, headers []string) (*http.Client, *http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	var key, value string
	for _, header := range headers {
		key, value = getHeaderKeyValue(header)
		req.Header.Add(key, value)
	}

	return client, req, nil
}

func getHeaderKeyValue(header string) (string, string) {
	headerPair := strings.Split(header, ":")
	return strings.TrimSpace(headerPair[0]), strings.TrimSpace(headerPair[1])
}

func makeRequest(client *http.Client, req *http.Request) error {
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.StatusCode != 200 {
		fmt.Printf("Response Body \n: %v\n", string(body))
		return fmt.Errorf("response status code %d", res.StatusCode)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().StringP("method", "M", "GET", "Http Method")
	profileCmd.Flags().StringP("url", "u", "", "Add API endpoint to be profiled")
	profileCmd.Flags().StringSliceP("header", "H", []string{}, "List of headers to be added")
	profileCmd.Flags().IntP("iteration", "I", 1, "Number of times the endpoint requests will be sent")
	profileCmd.Flags().BoolP("visualize", "v", false, "Save the data captured in a line graph")
	profileCmd.Flags().IntSliceP("pidlist", "p", []int{}, "List of processes to profile")
	profileCmd.Flags().IntP("interval", "i", 1000, "Interval at which the profiling is done (milliseconds)")
}
