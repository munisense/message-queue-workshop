package message

import "time"

type Message struct {
	Results []Result
	Meta    Metadata
}

type Result struct {
	Timestamp time.Time
	Value     string
}
type Metadata struct {
	Attribute_id int
	Cluster_id   int
	Eui64        string
	Gateway_id   int
	Unit         string
}
