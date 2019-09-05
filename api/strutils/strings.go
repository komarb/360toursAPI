package strutils

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeString(a []string, i int) []string {
	emptySlice := make([]string, 0)
	if len(a) == 1 {
		return emptySlice
	}
	a[i] = a[len(a)-1]
	a[len(a)-1] = ""
	return a[:len(a)-1]
}
