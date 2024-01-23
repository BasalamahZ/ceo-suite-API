package postgresql

const queryGetBooking = `
	SELECT
		b.id,
		b.user_id,
		b.product_id,
		p.date,
		p.start_time,
		p.end_time,
		p.price,
		b.status,
		b.create_time,
		b.update_time
	FROM
		booking b
	LEFT JOIN
		product p
	ON
		b.product_id = p.id
	%s
`
