package archive

type Archiver interface {
	Extract(src string, replacements map[string]interface{}) error
	LoadManifestFile(src string) ([]byte, error)
}
