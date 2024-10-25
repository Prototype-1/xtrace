package utils

import (
    "math/rand"
    "strconv"
    "time"
	"errors"
)

const otpLength = 6

func GenerateOTP() string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    otp := ""
    for i := 0; i < otpLength; i++ {
        otp += strconv.Itoa(r.Intn(10)) 
    }
    return otp
}

func CompareOTP(expectedOTP, providedOTP string) error {
    if expectedOTP != providedOTP {
        return errors.New("invalid OTP")
    }
    return nil
}