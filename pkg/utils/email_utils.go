package utils

import (
    "gopkg.in/gomail.v2"
    "os"
)

func SendEmail(to string, subject string, body string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", os.Getenv("EMAIL_SENDER"))
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/plain", body)

    dialer := gomail.NewDialer(
        os.Getenv("EMAIL_SMTP_HOST"),
        587, 
        os.Getenv("EMAIL_SENDER"),
        os.Getenv("EMAIL_PASSWORD"),
    )
    if err := dialer.DialAndSend(msg); err != nil {
        return err
    }

    return nil
}


func SendEmailWithAttachment(to string, subject string, body string, attachmentPath string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", os.Getenv("EMAIL_SENDER"))
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/plain", body)

    // Attach the file
    msg.Attach(attachmentPath)

    dialer := gomail.NewDialer(
        os.Getenv("EMAIL_SMTP_HOST"),
        587, 
        os.Getenv("EMAIL_SENDER"),
        os.Getenv("EMAIL_PASSWORD"),
    )

    if err := dialer.DialAndSend(msg); err != nil {
        return err
    }

    return nil
}