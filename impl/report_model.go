package impl

import (
	"encoding/json"
	"log"
)

// type ReportItem struct {
// 	Source SourceItem `json:"_source"`
// }

func (p ReportItem) FromJson(data []byte) (*ReportItem, error) {
	var result ReportItem
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
	}
	return &result, err
}

type ReportItem struct {
	Host          string     `json:"host"`
	Environment   string     `json:"environment"`
	PuppetVersion string     `json:"puppet_version"`
	Status        string     `json:"status"`
	Metrics       MetricItem `json:"metrics"`
	Logs          []LogItem  `json:"logs"`
}

type MetricItem struct {
	Resources MetricResourcesItem `json:"resources"`
	Events    MetricEventsItem    `json:"events"`
}

type MetricResourcesItem struct {
	Failed          int `json:"failed"`
	FailedToRestart int `json:"failed_to_restart"`
}

type MetricEventsItem struct {
	Failed int `json:"failure"`
}

type LogItem struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}
