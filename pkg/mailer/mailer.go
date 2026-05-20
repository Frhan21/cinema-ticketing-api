package mailer

import (
	"cinema-ticketing-api/internal/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer    *gomail.Dialer
	fromEmail string
	fromName  string
}

func NewMailer(cfg config.SMTPConfig) *Mailer {
	return &Mailer{
		dialer:    gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password),
		fromEmail: cfg.Sender,
		fromName:  "Cinema Ticketing",
	}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.fromName, m.fromEmail))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	return m.dialer.DialAndSend(msg)
}

// SendTicketPaidNotification mengirim email konfirmasi pembayaran tiket.
func (m *Mailer) SendTicketPaidNotification(to, userName, movieTitle, scheduleTime string, totalPrice float64) error {
	subject := "✅ Pembayaran Tiket Berhasil - " + movieTitle
	body := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
		<h2 style="color: #2d8a4e;">Pembayaran Berhasil!</h2>
		<p>Halo <strong>%s</strong>,</p>
		<p>Pembayaran tiket Anda telah berhasil. Berikut detailnya:</p>
		<table style="width:100%%; border-collapse: collapse;">
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Film</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Jadwal</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Total</strong></td><td style="padding: 8px; border: 1px solid #ddd;">Rp %.0f</td></tr>
		</table>
		<p style="margin-top: 20px;">Tunjukkan tiket ini saat masuk bioskop. Selamat menikmati film! 🎬</p>
	</div>
	`, userName, movieTitle, scheduleTime, totalPrice)
	return m.SendEmail(to, subject, body)
}

// SendTicketCancelledNotification mengirim email notifikasi pembatalan tiket.
func (m *Mailer) SendTicketCancelledNotification(to, userName, movieTitle string) error {
	subject := "❌ Tiket Dibatalkan - " + movieTitle
	body := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
		<h2 style="color: #e53e3e;">Tiket Dibatalkan</h2>
		<p>Halo <strong>%s</strong>,</p>
		<p>Tiket Anda untuk film <strong>%s</strong> telah dibatalkan.</p>
		<p>Jika ini bukan atas permintaan Anda, silakan hubungi customer service kami.</p>
	</div>
	`, userName, movieTitle)
	return m.SendEmail(to, subject, body)
}

// SendUpcomingFilmReminder mengirim email pengingat film yang akan tayang.
func (m *Mailer) SendUpcomingFilmReminder(to, userName, movieTitle, scheduleTime, studioName string) error {
	subject := "⏰ Pengingat Film - " + movieTitle + " 30 Menit Lagi!"
	body := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
		<h2 style="color: #d97706;">Film Anda Akan Segera Mulai!</h2>
		<p>Halo <strong>%s</strong>,</p>
		<p>Ini adalah pengingat bahwa film Anda akan dimulai dalam <strong>30 menit</strong>.</p>
		<table style="width:100%%; border-collapse: collapse;">
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Film</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Jadwal</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
			<tr><td style="padding: 8px; border: 1px solid #ddd;"><strong>Studio</strong></td><td style="padding: 8px; border: 1px solid #ddd;">%s</td></tr>
		</table>
		<p style="margin-top: 20px;">Harap tiba lebih awal. Selamat menikmati film! 🎬</p>
	</div>
	`, userName, movieTitle, scheduleTime, studioName)
	return m.SendEmail(to, subject, body)
}
