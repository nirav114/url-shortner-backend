package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/nirav114/url-shortner-backend.git/config"
	"gopkg.in/gomail.v2"
)

const OTPExpiry = 10 * time.Minute

func GenerateOTP(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	otp := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	return otp[:length], nil
}

func ValidateOTP(email string, submittedOTP string) (bool, error) {
	userData, err := RetrieveUserOTPData(email)
	if err != nil {
		return false, fmt.Errorf("error retrieving user data: %v", err)
	}

	if userData == nil {
		return false, fmt.Errorf("no OTP data found for email")
	}

	if userData.OTP != submittedOTP {
		return false, fmt.Errorf("invalid OTP")
	}

	ttl, err := RedisClient.TTL(ctx, email+":user").Result()
	if err != nil {
		return false, fmt.Errorf("error checking TTL: %v", err)
	}

	if ttl <= 0 {
		return false, fmt.Errorf("OTP has expired")
	}

	return true, nil
}

func SendOTPEmail(email string, otp string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EnvConfig.MAIL_ID)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Your OTP Code")
	mailer.SetBody("text/plain", fmt.Sprintf("Your OTP is: %s. It will expire in 10 minutes.", otp))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, config.EnvConfig.MAIL_ID, config.EnvConfig.APP_KEY)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}
