package consts

const (
	Version = "1754854321"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
	Platform   string `json:"platform"`
}
