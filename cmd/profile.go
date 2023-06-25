/*
Copyright © 2023 NAME HERE akashnandan99@gmail.com
*/
package cmd

import (
	"goburst/internal/profiler"
	"goburst/pkg/cutefmt"
	"goburst/pkg/visualizer"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Profie your api",
	Long:  `profile command runs the api endpoint multiple times while taking CPU and Memory Sample at a specified interval.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		headers, err := cmd.Flags().GetStringSlice("header")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		iteration, err := cmd.Flags().GetInt("iteration")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		method, err := cmd.Flags().GetString("method")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		pidList, err := cmd.Flags().GetIntSlice("pidlist")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		visualize, err := cmd.Flags().GetBool("visualize")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		taskName, err := cmd.Flags().GetString("name")
		if err != nil {
			cutefmt.Errorf("Error parsing --url flag :", err)
			return
		}

		profiler := profiler.NewProfiler(taskName, method, url, headers, pidList, iteration)
		processStats, err := profiler.Profile(method, url, headers, iteration, interval)
		if err != nil {
			cutefmt.Errorf("Profiling not complete due to error ", err)
			return
		}

		cutefmt.Successf("Profiling Complete")

		if visualize {
			visualizer.Visualize(taskName, processStats)
			cutefmt.Successf("Generated Graphs")
		}
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().StringP("method", "M", "GET", "Http Method")
	profileCmd.Flags().StringP("url", "u", "", "Add API endpoint to be profiled")
	profileCmd.Flags().StringSliceP("header", "H", []string{}, "List of headers to be added")
	profileCmd.Flags().IntP("iteration", "I", 1, "Number of times the endpoint requests will be sent")
	profileCmd.Flags().BoolP("visualize", "v", true, "Save the data captured in a line graph")
	profileCmd.Flags().IntSliceP("pidlist", "p", []int{}, "List of processes to profile")
	profileCmd.Flags().IntP("interval", "i", 1000, "Interval at which the profiling is done (milliseconds)")
	profileCmd.Flags().StringP("name", "n", "Perf Graph", "Title for the graphs generated")
}
