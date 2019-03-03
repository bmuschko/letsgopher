package templ

type Archiver interface {
	Extract(src string) error
}
