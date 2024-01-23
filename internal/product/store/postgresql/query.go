package postgresql

const queryGetProduct = `
	SELECT
		p.id,
		p.name,
		p.images,
		p.location,
		p.date,
		p.start_time,
		p.end_time,
		p.status,
		p.capacity,
		p.price,
		p.min_charge,
		p.daily_rate,
		p.promo,
		p.promo_price,
		p.address,
		p.distance,
		p.description,
		p.latitude,
		p.longitude,
		p.rating,
		p.create_time,
		p.update_time
	FROM
		product p
	%s
`
