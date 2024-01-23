package postgresql

import (
	"errors"
	"time"

	"github.com/ceo-suite/internal/product"
	"github.com/ceo-suite/internal/product/service"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	errInvalidCommit   = errors.New("cannot do commit on non-transactional querier")
	errInvalidRollback = errors.New("cannot do rollback on non-transactional querier")
)

// store implements user/service.PGStore
type store struct {
	db *sqlx.DB
}

// storeClient implements user/service.PGStoreClient
type storeClient struct {
	q sqlx.Ext
}

// New creates a new store.
func New(db *sqlx.DB) (*store, error) {
	s := &store{
		db: db,
	}

	return s, nil
}

func (s *store) NewClient(useTx bool) (service.PGStoreClient, error) {
	var q sqlx.Ext

	// determine what object should be use as querier
	q = s.db
	if useTx {
		var err error
		q, err = s.db.Beginx()
		if err != nil {
			return nil, err
		}
	}

	return &storeClient{
		q: q,
	}, nil
}

func (sc *storeClient) Commit() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Commit()
	}
	return errInvalidCommit
}

func (sc *storeClient) Rollback() error {
	if tx, ok := sc.q.(*sqlx.Tx); ok {
		return tx.Rollback()
	}
	return errInvalidRollback
}

// productDB denotes a data in the store.
type productDB struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Images      pq.StringArray `db:"images"`
	Location    string         `db:"location"`
	Date        time.Time      `db:"date"`
	StartTime   time.Time      `db:"start_time"`
	EndTime     time.Time      `db:"end_time"`
	Status      product.Status `db:"status"`
	Capacity    int            `db:"capacity"`
	Price       int64          `db:"price"`
	MinCharge   int64          `db:"min_charge"`
	DailyRate   int64          `db:"daily_rate"`
	Promo       bool           `db:"promo"`
	PromoPrice  *int64         `db:"promo_price"`
	Address     string         `db:"address"`
	Distance    float32        `db:"distance"`
	Description string         `db:"description"`
	Latitude    string         `db:"latitude"`
	Longitude   string         `db:"longitude"`
	Rating      float32        `db:"rating"`
	CreateTime  time.Time      `db:"create_time"`
	UpdateTime  *time.Time     `db:"update_time"`
}

// format formats database struct into domain struct.
func (pdb *productDB) format() product.Product {
	p := product.Product{
		ID:          pdb.ID,
		Name:        pdb.Name,
		Location:    pdb.Location,
		Date:        pdb.Date,
		StartTime:   pdb.StartTime,
		EndTime:     pdb.EndTime,
		Status:      pdb.Status,
		Capacity:    pdb.Capacity,
		Price:       pdb.Price,
		MinCharge:   pdb.MinCharge,
		DailyRate:   pdb.DailyRate,
		Promo:       pdb.Promo,
		Address:     pdb.Address,
		Distance:    pdb.Distance,
		Description: pdb.Description,
		Latitude:    pdb.Latitude,
		Longitude:   pdb.Longitude,
		Rating:      pdb.Rating,
		CreateTime:  pdb.CreateTime,
	}

	if len(pdb.Images) > 0 {
		images := make([]string, 0)
		for _, image := range pdb.Images {
			images = append(images, image)
		}
		p.Images = images
	}

	if pdb.PromoPrice != nil {
		p.PromoPrice = *pdb.PromoPrice
	}

	if pdb.UpdateTime != nil {
		p.UpdateTime = *pdb.UpdateTime
	}

	return p
}
