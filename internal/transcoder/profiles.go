package transcoder

type QualityProfile struct {
	ID         string
	Resolution string
	MaxRate    string
	BufSize    string
	CRF        string
}

var QualityProfiles = map[string]QualityProfile{
	"source": {Resolution: "source", MaxRate: "", BufSize: "", CRF: ""},
	"high":   {Resolution: "native", MaxRate: "35M", BufSize: "40M", CRF: "20"},
	"medium": {Resolution: "native", MaxRate: "12M", BufSize: "16M", CRF: "23"},
	"low":    {Resolution: "1080p", MaxRate: "8M", BufSize: "10M", CRF: "26"},
	"mobile": {Resolution: "720p", MaxRate: "4M", BufSize: "6M", CRF: "28"},
}

func GetQualityProfile(id string) (QualityProfile, bool) {
	profile, exists := QualityProfiles[id]
	return profile, exists
}

func GetAvailableQualities() []string {
	return []string{"source", "high", "medium", "low", "mobile"}
}
