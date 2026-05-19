package utils

func TicketPaidEmailTemplate(name string, movieTitle string, seatNumber string, startTime string) string {
	return `
		<h2>Pembayaran Tiket Berhasil</h2>
		<p>Halo ` + name + `,</p>
		<p>Pembayaran tiket kamu sudah berhasil.</p>
		<p><strong>Film:</strong> ` + movieTitle + `</p>
		<p><strong>Kursi:</strong> ` + seatNumber + `</p>
		<p><strong>Jadwal:</strong> ` + startTime + `</p>
		<p>Silakan datang tepat waktu.</p>
	`
}
