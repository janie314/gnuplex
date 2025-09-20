package consts

const (
	Version = "1758409293"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
	Platform   string `json:"platform"`
}
