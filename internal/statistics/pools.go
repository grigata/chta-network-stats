package statistics

import (
	"sort"

	"github.com/grigata/chta-network-stats/internal/models"
)

func AnalyzePools(blocks []models.NetworkBlock) []models.PoolStats {

	pools := make(map[string]*models.PoolStats)

	for _, block := range blocks {

		ps, ok := pools[block.Pool]
		if !ok {
			ps = &models.PoolStats{
				Name: block.Pool,
			}
			pools[block.Pool] = ps
		}

		ps.TotalBlocks++

		switch block.Type {
		case "NORMAL":
			ps.NormalBlocks++

		case "LOW-DIFF":
			ps.LowDiffBlocks++
		}
	}

	result := make([]models.PoolStats, 0, len(pools))

	for _, pool := range pools {
		pool.Percent =
			float64(pool.TotalBlocks) * 100 /
				float64(len(blocks))

		result = append(result, *pool)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalBlocks > result[j].TotalBlocks
	})

	return result
}
