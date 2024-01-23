package postgresql

import (
	"errors"
	"time"

	"github.com/ceo-suite/internal/booking"
	"github.com/ceo-suite/internal/booking/service"
	"github.com/jmoiron/sqlx"
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

// bookingDB denotes a data in the store.
type bookingDB struct {
	ID         int64          `db:"id"`
	UserID     int64          `db:"user_id"`
	ProductID  int64          `db:"product_id"`
	Date       time.Time      `db:"date"`
	StartTime  time.Time      `db:"start_time"`
	EndTime    time.Time      `db:"end_time"`
	Status     booking.Status `db:"status"`
	Price      int64          `db:"price"`
	CreateTime time.Time      `db:"create_time"`
	UpdateTime *time.Time     `db:"update_time"`
}

// format formats database struct into domain struct.
func (bdb *bookingDB) format() booking.Booking {
	b := booking.Booking{
		ID:         bdb.ID,
		UserID:     bdb.UserID,
		ProductID:  bdb.ProductID,
		Date:       bdb.Date,
		StartTime:  bdb.StartTime,
		EndTime:    bdb.EndTime,
		Status:     bdb.Status,
		Price:      bdb.Price,
		CreateTime: bdb.CreateTime,
	}

	if bdb.UpdateTime != nil {
		b.UpdateTime = *bdb.UpdateTime
	}

	return b
}
