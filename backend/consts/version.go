package consts

const (
	Version = "1755969359"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
	Platform   string `json:"platform"`
}
