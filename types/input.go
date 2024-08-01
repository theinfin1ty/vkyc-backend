package types

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	OTP      string `json:"otp" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type AddUserInput struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Role        string `json:"role" binding:"required"`
	CountryCode string `json:"countryCode"`
	Phone       string `json:"phone"`
}

type UpdateUserInput struct {
	FirstName   string  `json:"firstName" binding:"required"`
	LastName    string  `json:"lastName" binding:"required"`
	Email       string  `json:"email" binding:"required"`
	CountryCode string  `json:"countryCode"`
	Phone       string  `json:"phone"`
	DOB         string  `json:"dob"`
	Department  *string `json:"department"`
}

type MeetingInput struct {
	Title            string  `json:"title"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Email            string  `json:"email"`
	CountryCode      string  `json:"countryCode"`
	Phone            string  `json:"phone"`
	Remark           string  `json:"remark"`
	ScheduleDateTime *string `json:"scheduleDateTime"`
	AgentID          *int    `json:"agentID"`
}

type JoinMeetingInput struct {
	MeetingCode string `json:"meetingCode"`
	Password    string `json:"password"`
}

type VerificationInput struct {
	MeetingID int    `json:"meetingID"`
	Document  string `json:"document"`
}

type QuestionInput struct {
	Question string `json:"question"`
	Status   string `json:"status"`
}

type ConfigurationInput struct {
	Name        *string `json:"name"`
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}
