package http

import (
	"github.com/ceo-suite/internal/product"
)

// formatProduct formats the given product
// into the respective HTTP-format object.
func formatProduct(p product.Product) (productHTTP, error) {
	status := p.Status.String()
	date := p.Date.Format(dateFormat)
	startTime := p.StartTime.Format(timeFormat)
	endTime := p.EndTime.Format(timeFormat)

	return productHTTP{
		ID:          &p.ID,
		Name:        &p.Name,
		Images:      &p.Images,
		Location:    &p.Location,
		Date:        &date,
		StartTime:   &startTime,
		EndTime:     &endTime,
		Status:      &status,
		Capacity:    &p.Capacity,
		Price:       &p.Price,
		MinCharge:   &p.MinCharge,
		DailyRate:   &p.DailyRate,
		Promo:       &p.Promo,
		PromoPrice:  &p.PromoPrice,
		Address:     &p.Address,
		Distance:    &p.Distance,
		Description: &p.Description,
		Latitude:    &p.Latitude,
		Longitude:   &p.Longitude,
		Rating:      &p.Rating,
	}, nil
}
