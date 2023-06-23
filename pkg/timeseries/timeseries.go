package timeseries

import (
	"time"
)

type DataPoint []int

type TimeSeries struct {
	data []DataPoint
}

func NewDataPoint(time int, value int) DataPoint {
	return []int{time, value}
}

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{}
}

func (timeSeries *TimeSeries) AddAtTime(unixtime string, data DataPoint) {
	timeSeries.data = append(timeSeries.data, data)
}

func (timeSeries *TimeSeries) Add(value int) {
	timeSeries.data = append(timeSeries.data, NewDataPoint(int(time.Now().Unix()), value))
}

func (timeSeries *TimeSeries) GetTimeList() (timeList []int) {
	for _, datapoint := range timeSeries.data {
		timeList = append(timeList, datapoint[0])
	}
	return
}

func (timeSeries *TimeSeries) GetValues() (values []int) {
	for _, datapoint := range timeSeries.data {
		values = append(values, datapoint[0])
	}
	return
}

func (timeSeries *TimeSeries) Size() int {
	return len(timeSeries.data)
}

func (timeSeries *TimeSeries) GetValue(index int) int {
	return timeSeries.data[index][1]
}

func (timeSeries *TimeSeries) GetTime(index int) int {
	return timeSeries.data[index][0]
}

func (timeSeries *TimeSeries) StartTime() int {
	return timeSeries.data[0][0]
}

func (timeSeries *TimeSeries) EndTime() int {
	return timeSeries.data[timeSeries.Size()-1][0]
}
