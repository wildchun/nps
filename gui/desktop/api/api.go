package api

func GetUrl(path string) string {
	return "http://" + Server + path
}
