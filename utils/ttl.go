package utils

func BuildTTL(ttl int) int {
	var min = 60
	var max = 86400

	if ttl < min {
		// Auto TTL
		return 1
	} else if ttl >= max {
		// TTL Upperbound
		return max
	} else {
		return ttl
	}
}
