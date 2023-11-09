package initializers

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(email string) {
	msg := []byte("From: Sora <1371934864@qq.com>\r\n" + "To: recipient@example.net\r\n" + "Subject:Reset your password!\r\n" + "\r\n" + "Please click this link and reset password.\r\n")
	destEmail := email
	auth := smtp.PlainAuth("", "1371934864@qq.com", os.Getenv("SMTPPWD"), "smtp.qq.com")

	err := smtp.SendMail("smtp.qq.com:587", auth, "1371934864@qq.com", []string{destEmail}, msg)
	if err != nil {
		fmt.Println(err)
	}

}
