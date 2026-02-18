package models

// StatusDescriptions provides optional human-readable descriptions for ticket statuses
// This is NOT an enforcement list - the system accepts any status value from the database
// New statuses will automatically appear in metadata, using their code as description until documented here
var StatusDescriptions = map[string]string{
	"0.NEW":                   "New ticket",
	"1.Req FD ke HD":          "Request FD to HD",
	"2.Kirim FLM":             "Send FLM",
	"21.Req Replenish":        "Request Replenish",
	"3.SLM":                   "SLM Machine",
	"4.SLM-Net":               "SLM Network",
	"5.Menunggu Update":       "Waiting for Update",
	"6.Follow-up Sales team":  "Follow-up Sales team",
	"8.Wait transaction":      "Wait transaction",
}

// ModeDescriptions provides optional descriptions for terminal modes
var ModeDescriptions = map[string]string{
	"Closed":     "Terminal is closed",
	"In Service": "Terminal is in service",
	"nan":        "No mode data available",
	"Off-line":   "Terminal is offline",
	"Supervisor": "Supervisor mode",
}

// PriorityDescriptions provides optional descriptions for ticket priorities
var PriorityDescriptions = map[string]string{
	"1.High":    "High priority",
	"2.Middle":  "Middle priority",
	"3.Low":     "Low priority",
	"4.Minimum": "Minimum priority",
}

// StatusInfo provides detailed information about a ticket status
type StatusInfo struct {
	Code        string `json:"code" example:"0.NEW"`                    // Status code from database
	Description string `json:"description" example:"New ticket"`        // Human-readable description
	IsDocumented bool  `json:"is_documented" example:"true"`            // Whether this value has documentation
}

// ModeInfo provides detailed information about a terminal mode
type ModeInfo struct {
	Code        string `json:"code" example:"Off-line"`                        // Mode code from database
	Description string `json:"description" example:"Terminal is offline"`      // Human-readable description
	IsDocumented bool  `json:"is_documented" example:"true"`                   // Whether this value has documentation
}

// PriorityInfo provides detailed information about a priority level
type PriorityInfo struct {
	Code        string `json:"code" example:"1.High"`                   // Priority code from database
	Description string `json:"description" example:"High priority"`     // Human-readable description
	IsDocumented bool  `json:"is_documented" example:"true"`            // Whether this value has documentation
}

// BuildStatusInfo creates StatusInfo with optional description
func BuildStatusInfo(code string) StatusInfo {
	desc, documented := StatusDescriptions[code]
	if !documented {
		desc = code // Use code as description for undocumented values
	}
	return StatusInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// BuildModeInfo creates ModeInfo with optional description
func BuildModeInfo(code string) ModeInfo {
	desc, documented := ModeDescriptions[code]
	if !documented {
		desc = code
	}
	return ModeInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// BuildPriorityInfo creates PriorityInfo with optional description
func BuildPriorityInfo(code string) PriorityInfo {
	desc, documented := PriorityDescriptions[code]
	if !documented {
		desc = code
	}
	return PriorityInfo{
		Code:         code,
		Description:  desc,
		IsDocumented: documented,
	}
}

// MetadataResponse provides information about valid ticket field values
// All values are queried from the database in real-time for accuracy
type MetadataResponse struct {
	Success    bool           `json:"success" example:"true"`                           // Operation success status
	Message    string         `json:"message" example:"Metadata retrieved successfully"` // Response message
	Statuses   []StatusInfo   `json:"statuses"`                                         // Available status values from database
	Modes      []ModeInfo     `json:"modes"`                                            // Available mode values from database
	Priorities []PriorityInfo `json:"priorities"`                                       // Available priority values from database
	LastUpdated string        `json:"last_updated" example:"2024-01-15T10:30:00Z"`     // When metadata was last refreshed
}

