package api

func init() {
	// http.DefaultClient.Transport
}

func GetUrl(path string) string {
	return "http://" + Server + path
}