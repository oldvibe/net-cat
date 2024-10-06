package ncat

func ValidPort(port string) bool {
	for _, ch := range port {
		if ch < 48 || ch > 57 {
			return false
		}
	}

	if len(port) > 5 {
		return false
	}

	if port < "1" || port > "65535" {
		return false
	}

	if port <= "1023" {
		return false
	}

	return true
}
