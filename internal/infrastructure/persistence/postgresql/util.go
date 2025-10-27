package postgresql

import "strconv"

func joinNums(nums []int, sep string) string {
	if len(nums) == 0 {
		return ""
	}

	s := strconv.Itoa(nums[0])

	for _, n := range nums[1:] {
		s += sep + strconv.Itoa(n)
	}

	return s
}
