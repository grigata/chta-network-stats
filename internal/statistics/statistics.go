package statistics

import (
	"sort"

	"github.com/grigata/chta-network-stats/internal/models"
)

type Statistics struct {
	TotalBlocks int

	NormalBlocks  int
	CheetahBlocks int

	AverageGap float64
	MinGap     int64
	MaxGap     int64

	Pools map[string]int
}

func Calculate(blocks []models.NetworkBlock) Statistics {

	stats := Statistics{
		TotalBlocks: len(blocks),
		Pools:       make(map[string]int),
	}

	if len(blocks) == 0 {
		return stats
	}

	stats.MinGap = blocks[0].Gap

	var totalGap int64

	for _, block := range blocks {

		stats.Pools[block.Pool]++

		if block.Type == "CHEETAH" {
			stats.CheetahBlocks++
		} else {
			stats.NormalBlocks++
		}

		totalGap += block.Gap

		if block.Gap < stats.MinGap {
			stats.MinGap = block.Gap
		}

		if block.Gap > stats.MaxGap {
			stats.MaxGap = block.Gap
		}
	}

	stats.AverageGap = float64(totalGap) / float64(len(blocks))

	return stats

}
func SortedPools(stats Statistics) []PoolStat {

	pools := make([]PoolStat, 0, len(stats.Pools))

	for name, blocks := range stats.Pools {

		pools = append(pools, PoolStat{
			Name:    name,
			Blocks:  blocks,
			Percent: float64(blocks) * 100 / float64(stats.TotalBlocks),
		})
	}

	sort.Slice(pools, func(i, j int) bool {

		if pools[i].Blocks == pools[j].Blocks {
			return pools[i].Name < pools[j].Name
		}

		return pools[i].Blocks > pools[j].Blocks
	})

	return pools
}
