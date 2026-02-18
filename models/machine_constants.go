package models

// SLMDescriptions provides optional human-readable descriptions for SLM types
// This is NOT an enforcement list - the system accepts any SLM value from the database
// New SLM values will automatically appear in metadata until documented here
var SLMDescriptions = map[string]string{
	"KGP - WINCOR DW":   "KGP - WINCOR DW",
	"GPS - NCR":         "GPS - NCR",
	"NCR":               "NCR",
	"NCR/ASABA":         "NCR/ASABA",
	"WINCOR - KGP":      "WINCOR - KGP",
	"WINCOR":            "WINCOR",
	"ARGENTA":           "ARGENTA",
	"NCR - KGP":         "NCR - KGP",
	"KGP - WINCOR DU":   "KGP - WINCOR DU",
	"KGP - NCR":         "KGP - NCR",
	"GPS":               "GPS",
	"KGP-WINCOR DW":     "KGP-WINCOR DW",
	"ASABA-NCR":         "ASABA-NCR",
	"KGP - WINCOR":      "KGP - WINCOR",
	"WINCOR DU - KGP":   "WINCOR DU - KGP",
	"":                  "No SLM data available",
	"WINCOR DW - KGP":   "WINCOR DW - KGP",
	"NCR - ASABA":       "NCR - ASABA",
	"KGP-NCR":           "KGP-NCR",
	"KGP-WINCOR":        "KGP-WINCOR",
	"KGP-WINCOR HG":     "KGP-WINCOR HG",
	"GPS - WINCOR DU":   "GPS - WINCOR DU",
	"KGP":               "KGP",
	"KGP-WINCOR DU":     "KGP-WINCOR DU",
	"NCR-KGP":           "NCR-KGP",
	"KGP - WINCOR HG":   "KGP - WINCOR HG",
}

// NETDescriptions provides optional descriptions for network providers
var NETDescriptions = map[string]string{
	"NOSAIRIS": "NOSAIRIS",
	"SMS":      "SMS",
	"TANGARA":  "TANGARA",
	"":         "No NET data available",
	"SM":       "SM",
	"IFORTE":   "IFORTE",
}

// FLMNameDescriptions provides descriptions for FLM provider names
var FLMNameDescriptions = map[string]string{
	"ABS": "ABS",
	"BRS": "BRS",
	"TAG": "TAG",
	"":    "No FLM Name data available",
	"AVT": "AVT",
}

// FLMAreaMap maps FLM codes to their service areas
// This is reference data that helps associate FLM providers with geographic areas
var FLMAreaMap = map[string]string{
	"ABS - BALI":         "BALI",
	"ABS - BANDUNG":      "BANDUNG",
	"ABS - BOGOR":        "BOGOR",
	"ABS - CIAMIS":       "CIAMIS",
	"ABS - JAKARTA":      "JAKARTA",
	"ABS - KEDIRI":       "KEDIRI",
	"ABS - MAKASSAR":     "MAKASSAR",
	"ABS - MALANG":       "MALANG",
	"ABS - MEDAN":        "MEDAN",
	"ABS - SIDOARJO":     "SIDOARJO",
	"ABS - SUBANG":       "SUBANG",
	"ABS - SUKABUMI":     "SUKABUMI",
	"ABS - SURABAYA":     "SURABAYA",
	"AVT - BALI":         "BALI",
	"AVT - BALIKPAPAN":   "BALIKPAPAN",
	"AVT - BANDUNG":      "BANDUNG",
	"AVT - BANJARMASIN":  "BANJARMASIN",
	"AVT - BATAM":        "BATAM",
	"AVT - BENGKULU":     "BENGKULU",
	"AVT - CIDENG":       "CIDENG",
	"AVT - CIREBON":      "CIREBON",
	"AVT - JAMBI":        "JAMBI",
	"AVT - JEMBER":       "JEMBER",
	"AVT - KARAWANG":     "KARAWANG",
	"AVT - KEDIRI":       "KEDIRI",
	"AVT - KUDUS":        "KUDUS",
	"AVT - KUPANG":       "KUPANG",
	"AVT - LAMPUNG":      "LAMPUNG",
	"AVT - MALANG":       "MALANG",
	"AVT - MANADO":       "MANADO",
	"AVT - MATARAM":      "MATARAM",
	"AVT - MEDAN":        "MEDAN",
	"AVT - MERUYA":       "MERUYA",
	"AVT - PALEMBANG":    "PALEMBANG",
	"AVT - PEKANBARU":    "PEKANBARU",
	"AVT - PONTIANAK":    "PONTIANAK",
	"AVT - PURWOKERTO":   "PURWOKERTO",
	"AVT - RAWAMANGUN":   "RAWAMANGUN",
	"AVT - SAMARINDA":    "SAMARINDA",
	"AVT - SEMARANG":     "SEMARANG",
	"AVT - SERANG":       "SERANG",
	"AVT - SINGKAWANG":   "SINGKAWANG",
	"AVT - SOLO":         "SOLO",
	"AVT - SURABAYA":     "SURABAYA",
	"AVT - TASIKMALAYA":  "TASIKMALAYA",
	"AVT - TEGAL":        "TEGAL",
	"AVT - YOGYAKARTA":   "YOGYAKARTA",
	"BRS - BANDUNG":      "BANDUNG",
	"BRS - BOGOR":        "BOGOR",
	"BRS - CIREBON":      "CIREBON",
	"BRS - LEBAK BULUS":  "LEBAK BULUS",
	"BRS - MANADO":       "MANADO",
	"BRS - MEDAN":        "MEDAN",
	"BRS - SURABAYA":     "SURABAYA",
	"TAG - CIMONE":       "CIMONE",
	"-":                  "",
}

// SLMInfo provides detailed information about SLM types
type SLMInfo struct {
	Code         string `json:"code" example:"KGP - WINCOR DW"`
	Description  string `json:"description" example:"KGP - WINCOR DW"`
	IsDocumented bool   `json:"is_documented" example:"true"` // Whether this value has documentation
}

// FLMInfo provides detailed information about FLM providers
type FLMInfo struct {
	Code         string `json:"code" example:"AVT - BANDUNG"`
	Description  string `json:"description" example:"AVT - BANDUNG"`
	Area         string `json:"area,omitempty" example:"BANDUNG"`
	IsDocumented bool   `json:"is_documented" example:"true"` // Whether this value has documentation
}

// NETInfo provides detailed information about network providers
type NETInfo struct {
	Code         string `json:"code" example:"NOSAIRIS"`
	Description  string `json:"description" example:"NOSAIRIS"`
	IsDocumented bool   `json:"is_documented" example:"true"` // Whether this value has documentation
}

// FLMNameInfo provides information about FLM provider names
type FLMNameInfo struct {
	Code         string `json:"code" example:"AVT"`
	Description  string `json:"description" example:"AVT"`
	IsDocumented bool   `json:"is_documented" example:"true"` // Whether this value has documentation
}

// BuildSLMInfo creates SLMInfo with optional description
func BuildSLMInfo(code string) SLMInfo {
	desc, documented := SLMDescriptions[code]
	if !documented {
		desc = code
	}
	return SLMInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// BuildFLMInfo creates FLMInfo with optional description and area lookup
func BuildFLMInfo(code string) FLMInfo {
	// Check if we have area mapping
	area := FLMAreaMap[code]

	// Use code as description (FLM codes are self-descriptive)
	desc := code
	if code == "-" || code == "" {
		desc = "No FLM data available"
	}

	// Check if documented (exists in our area map)
	_, documented := FLMAreaMap[code]

	return FLMInfo{
		Code:         code,
		Description:  desc,
		Area:         area,
		IsDocumented: documented,
	}
}

// BuildNETInfo creates NETInfo with optional description
func BuildNETInfo(code string) NETInfo {
	desc, documented := NETDescriptions[code]
	if !documented {
		desc = code
	}
	return NETInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// BuildFLMNameInfo creates FLMNameInfo with optional description
func BuildFLMNameInfo(code string) FLMNameInfo {
	desc, documented := FLMNameDescriptions[code]
	if !documented {
		desc = code
	}
	return FLMNameInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// GetFLMArea returns the service area for a given FLM code
func GetFLMArea(flmCode string) string {
	if area, ok := FLMAreaMap[flmCode]; ok {
		return area
	}
	return ""
}

// MachineMetadataResponse provides information about valid machine field values
// All values are queried from the database in real-time for accuracy
type MachineMetadataResponse struct {
	Success     bool          `json:"success" example:"true"`
	Message     string        `json:"message" example:"Machine metadata retrieved successfully"`
	SLMs        []SLMInfo     `json:"slms"`      // Available SLM types from database
	FLMs        []FLMInfo     `json:"flms"`      // Available FLM providers from database
	NETs        []NETInfo     `json:"nets"`      // Available network providers from database
	FLMNames    []FLMNameInfo `json:"flm_names"` // Available FLM provider names from database
	LastUpdated string        `json:"last_updated" example:"2024-01-15T10:30:00Z"` // When metadata was last refreshed
}

