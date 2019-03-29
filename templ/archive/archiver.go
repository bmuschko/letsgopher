package archive

type Archiver interface {
	Extract(archiveFile string, targetDir string, replacements map[string]interface{}) error
	LoadManifestFile(src string) ([]byte, error)
}
