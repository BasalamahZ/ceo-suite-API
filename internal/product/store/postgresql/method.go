package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ceo-suite/internal/product"
	"github.com/jmoiron/sqlx"
)

// GetProductByID returns a product with the given product ID.
func (sc *storeClient) GetProductByID(ctx context.Context, id int64) (product.Product, error) {
	query := fmt.Sprintf(queryGetProduct, "WHERE p.id = $1")
	// query single row
	var pdb productDB
	err := sc.q.QueryRowx(query, id).StructScan(&pdb)
	if err != nil {
		if err == sql.ErrNoRows {
			return product.Product{}, product.ErrDataNotFound
		}
		return product.Product{}, err
	}

	return pdb.format(), nil
}

// GetProducts returns list of products that satisfy the given
// filter.
func (sc *storeClient) GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if filter.Location != "" {
		addConditions = append(addConditions, "p.location LIKE :location OR p.address LIKE :location")
		argsKV["location"] = "%" + filter.Location + "%"
	}

	if !filter.StartTime.IsZero() {
		addConditions = append(addConditions, "p.start_time >= :start_time")
		argsKV["start_time"] = filter.StartTime
	}

	if !filter.EndTime.IsZero() {
		addConditions = append(addConditions, "p.end_time <= :end_time")
		argsKV["end_time"] = filter.EndTime
	}

	if !filter.Date.IsZero() {
		addConditions = append(addConditions, "p.date = :date")
		argsKV["date"] = filter.Date
	}

	if filter.Capacity > 0 {
		addConditions = append(addConditions, "p.capacity <= :capacity")
		argsKV["capacity"] = filter.Capacity
	}

	if filter.Promo == true {
		addConditions = append(addConditions, "p.promo = :promo")
		argsKV["promo"] = filter.Promo
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}

	// construct query
	query := fmt.Sprintf(queryGetProduct, addCondition)

	// prepare query
	query, args, err := sqlx.Named(query, argsKV)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sc.q.Rebind(query)

	// query to database
	rows, err := sc.q.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// read products
	products := make([]product.Product, 0)
	for rows.Next() {
		var row productDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		products = append(products, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
