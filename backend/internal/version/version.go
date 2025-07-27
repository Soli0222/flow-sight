package version

var (
	// Application version - update this manually
	Version = "1.0.1"
)

type BuildInfo struct {
	Version string `json:"version"`
}

func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version: Version,
	}
}
