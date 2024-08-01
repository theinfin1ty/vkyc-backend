package utils

const (
	DateLayout = "2006-01-02T15:04:05Z"
	DateFormat = "02-01-2006 15:04:05"

	ActiveStatus   = "active"
	InactiveStatus = "inactive"
	BlockedStatus  = "blocked"
	PassStatus     = "pass"
	FailStatus     = "fail"
	ApprovedStatus = "approved"
	RejectedStatus = "rejected"

	QAFailedRemark              = "QA_FAILED"
	OCRFailedRemark             = "OCR_FAILED"
	FacematchFailedRemark       = "FACE_MATCH_FAILED"
	LivenessFailedRemark        = "LIVENESS_FAILED"
	SignFailedRemark            = "SIGNATURE_FAILED"
	FailedAfterConnectingRemark = "FAILED_AFTER_CONNECTING"

	AgentRole = "agent"
	AdminRole = "admin"

	MeetingExpireBeforeStart = "MEETING_EXPIRY_BEFORE_START"
	MeetingExpireAfterStart  = "MEETING_EXPIRY_AFTER_START"
	MeetingTimeGap           = "MEETING_TIME_GAP"
	FacematchAPI             = "FACEMATCH_API"
	FacematchAPIKey          = "FACEMATCH_API_KEY"
	LivenessAPI              = "LIVENESS_API"
	LivenessAPIKey           = "LIVENESS_API_KEY"
	OCRAPI                   = "OCR_API"
	OCRAPIKey                = "OCR_API_KEY"
	SMSAPI                   = "SMS_API"
	SMSAPIKey                = "SMS_API_KEY"
	SMSSender                = "SMS_SENDER"
	SMSTemplate              = "SMS_TEMPLATE"
	MinFacematchScore        = "MIN_FACEMATCH_SCORE"
	MinLivenessScore         = "MIN_LIVENESS_SCORE"
	CompanyName              = "COMPANY_NAME"
	CompanyLogo              = "COMPANY_LOGO"
	CompanyLogo2             = "COMPANY_LOGO_2"
)
