package artifactory

func Bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
