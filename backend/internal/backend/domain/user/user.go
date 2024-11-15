package user

import "regexp"

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (creds *Credentials) IsValidLogin() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return creds.Login != "" && re.MatchString(creds.Login)
}

func (creds *Credentials) IsValidPassword() bool {
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(creds.Password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(creds.Password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(creds.Password)
	hasSpecialChar := regexp.MustCompile(`[@$!%*?&]`).MatchString(creds.Password)
	hasMinLength := len(creds.Password) >= 8
	hasMaxLength := len(creds.Password) <= 48

	return hasLowercase && hasUppercase && hasDigit && hasSpecialChar && hasMinLength && hasMaxLength
}

func (creds *Credentials) IsValid() bool {
	return creds.IsValidLogin() && creds.IsValidPassword()
}

type User struct {
	Credentials Credentials
}

func NewUser(credentials Credentials) *User {
	return &User{
		Credentials: credentials,
	}
}
