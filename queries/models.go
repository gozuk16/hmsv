// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package queries

import ()

type CpuDownsampled struct {
	ID          int64
	Timestamp   string
	AvgCpuUsage float64
	MaxCpuUsage float64
}

type CpuOriginal struct {
	ID        int64
	Timestamp string
	CpuUsage  float64
}
