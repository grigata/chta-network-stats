package models

type PoolStats struct {
	Name          string
	TotalBlocks   int
	NormalBlocks  int
	LowDiffBlocks int
	Percent       float64
}
