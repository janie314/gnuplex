package consts

const (
	Version = "1759079502"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
	Platform   string `json:"platform"`
	GoVersion  string `json:"go_version"`
}
