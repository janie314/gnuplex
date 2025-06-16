package consts

const (
	Version = "1750036030"
)

type VersionInfo struct {
	Version    string `json:"version"`
	SourceHash string `json:"source_hash"`
}
