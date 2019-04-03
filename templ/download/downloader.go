package download

type Downloader interface {
	Download(url string) (string, error)
}
