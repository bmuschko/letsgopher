package templ

type Archiver interface {
	Extract(src string, replacements map[string]string) error
	LoadFile(src string) ([]byte, error)
}
