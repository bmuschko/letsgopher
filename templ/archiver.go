package templ

type Archiver interface {
	Extract(src, dest string) (string, error)
}
