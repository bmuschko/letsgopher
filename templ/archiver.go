package templ

type Archiver interface {
	Extract(src string) error
	LoadFile(src string) (string, error)
}
