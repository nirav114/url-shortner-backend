package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"
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
