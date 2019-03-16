package templ

type Archiver interface {
	Extract(src string, replacements map[string]interface{}) error
	LoadFile(src string) ([]byte, error)
}
