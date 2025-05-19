package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendEmailWithAccountInfo(receiver, receiverName, password string) error {
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
				}
				.header {
					color: #C22621;
					font-size: 20px;
					margin-bottom: 24px;
					display: flex; 
					justify-content: center; 
					align-items: center;
				}
				.content {
					background-color: #f8f9fa;
					padding: 24px;
					border-radius: 8px;
				}
				.account-info {
					background-color: #fff;
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
				}
				.footer {
					margin-top: 24px;
					padding-top: 16px;
					border-top: 1px solid #e2e8f0;
					font-size: 14px;
					color: #666;
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
				<div class="header">
					Chào mừng bạn đến với ERP!
				</div>
				
				<p>Xin chào <b>%s</b>,</p>
				
				<p>Chúc mừng! Tài khoản ERP của bạn đã được tạo thành công. Dưới đây là thông tin đăng nhập của bạn:</p>
				
				<div class="account-info">
					<p><strong>Email đăng nhập:</strong> %s</p>
					<p><strong>Mật khẩu:</strong> %s</p>
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
					<p><i>Đây là email tự động. Vui lòng không trả lời email này.</i></p>
					<p>© 2025 ERP. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, receiverName, receiver, password)

	m := gomail.NewMessage()
	m.SetHeader("From", "tranvandu3802@gmail.com")
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "[HỆ THỐNG APP CHẤM CÔNG] TÀI KHOẢN MAIL SỬ DỤNG ĐỂ ĐĂNG NHẬP")
	m.SetBody("text/html", emailContent)

	// d := gomail.NewDialer(mailer.Host, mailer.Port, mailer.UserName, mailer.Password)
	d := gomail.NewDialer("smtp.gmail.com", 587, "tranvandu3802@gmail.com", "hzim xqcv ebff hfot")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
