package download

// Downloader retrieves a template from an URL.
type Downloader interface {
	Download(url string) (string, error)
}
