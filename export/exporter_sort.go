package export

import (
	"sort"
	"strconv"
)

// Implementation of the sort interface of exporter.
// Since there is no way to identify number and string
// just simplify use the isNumeric func to identify
// number and string.

type detailRowSlice struct {
	detailRows  [][]string
	sortColumns []int
}

var _ sort.Interface = (*detailRowSlice)(nil)

func (s *detailRowSlice) Len() int {
	return len(s.detailRows)
}

func (s *detailRowSlice) Swap(i, j int) {
	s.detailRows[i], s.detailRows[j] = s.detailRows[j], s.detailRows[i]
}

func (s *detailRowSlice) Less(i, j int) bool {
	for _, sortColumn := range s.sortColumns {
		iValue, iErr := strconv.ParseFloat(s.detailRows[i][sortColumn], 64)
		jValue, jErr := strconv.ParseFloat(s.detailRows[j][sortColumn], 64)

		if iErr == nil && jErr == nil && iValue != jValue {
			return iValue < jValue
		}

		if s.detailRows[i][sortColumn] != s.detailRows[j][sortColumn] {
			return s.detailRows[i][sortColumn] < s.detailRows[j][sortColumn]
		}
	}

	return true
}
