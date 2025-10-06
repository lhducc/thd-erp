package mail

import (
	"crypto/tls"
	"erp/backend/config"
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendEmailWithAccountInfo(receiver, receiverName, accountPassword string) error {
	c := config.Mail

	// Build email content
	emailContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
				max-width: 600px;
				margin: 20px auto;
				padding: 20px;
				background-color: #f0f2f5;
			}
			.header {
				color: #C22621;
				font-size: 20px;
				margin-bottom: 24px;
				text-align: center;
				font-weight: bold;
			}
			.content {
				background-color: #ffffff;
				padding: 24px;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0,0,0,0.05);
			}
			.account-info {
				background-color: #f8f9fa;
				padding: 16px;
				border-radius: 6px;
				margin: 16px 0;
				border: 1px solid #e2e8f0;
			}
			.notice {
				background-color: #fff3cd;
				border-left: 4px solid #ffc107;
				padding: 12px 16px;
				margin: 16px 0;
				border-radius: 4px;
			}
			.footer {
				margin-top: 24px;
				padding-top: 16px;
				border-top: 1px solid #e2e8f0;
				font-size: 14px;
				color: #666;
				text-align: center;
			}
			ul, ol {
				margin: 16px 0;
				padding-left: 24px;
			}
			li {
				margin: 8px 0;
			}
		</style>
	</head>
	<body>
		<div class="content">
			<div class="header">Chào mừng bạn đến với ERP!</div>
			<p>Xin chào <strong>%s</strong>,</p>
			<p>Chúc mừng! Tài khoản ERP của bạn đã được tạo thành công. Dưới đây là thông tin đăng nhập của bạn:</p>
			<div class="account-info">
				<p><strong>Email đăng nhập:</strong> %s</p>
				<p><strong>Mật khẩu:</strong> %s</p>
				<p><strong>Đường dẫn truy cập hệ thống:</strong> 
					<a href="https://erp.thdcybersecurity.com/" target="_blank">https://erp.thdcybersecurity.com/</a>
				</p>
			</div>
			<div class="notice">
				<strong>Lưu ý quan trọng:</strong>
				<ol>
					<li>Để đảm bảo an toàn, vui lòng đổi mật khẩu ngay sau khi đăng nhập lần đầu tiên.</li>
					<li>Không chia sẻ thông tin tài khoản với bất kỳ ai để bảo vệ dữ liệu của bạn.</li>
					<li>Nếu bạn gặp bất kỳ khó khăn nào trong quá trình đăng nhập, vui lòng liên hệ với bộ phận hỗ trợ của chúng tôi.</li>
				</ol>
			</div>
			<p>Trân trọng,</p>
			<div class="footer">
				<p><em>Đây là email tự động. Vui lòng không trả lời email này.</em></p>
				<p>© 2025 ERP. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, receiverName, receiver, accountPassword)

	// CreateElementOfTimesheetList message
	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "[ERP] Tài khoản đăng nhập của bạn")
	m.SetBody("text/html", emailContent)

	// SMTP dialer with STARTTLS
	dialer := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	dialer.SSL = false
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.Host,
	}

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("gửi email thất bại: %w", err)
	}

	return nil
}
