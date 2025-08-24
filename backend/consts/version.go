package consts

const (
	Version = "1756048384"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
	Platform   string `json:"platform"`
}
