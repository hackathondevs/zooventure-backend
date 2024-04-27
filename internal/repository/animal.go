package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type AnimalRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (AnimalQueryerItf, error)
}

type AnimalRepository struct {
	db *sqlx.DB
}

func NewAnimalRepository(db *sqlx.DB) AnimalRepositoryItf {
	return &AnimalRepository{db}
}

type AnimalQueryerItf interface {
	txCompat
	Save(ctx context.Context, animal *model.Animal) (int64, error)
	GetAll(ctx context.Context) ([]model.Animal, error)
	FetchTopRelated(ctx context.Context, name string, lat, long float64) (model.Animal, error)
}

type animalQueryer struct {
	ext sqlx.ExtContext
}

func (r *AnimalRepository) NewClient(withTx bool, tx sqlx.ExtContext) (AnimalQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &animalQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &animalQueryer{tx}, nil
	}
	return &animalQueryer{r.db}, nil
}

func (q *animalQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *animalQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *animalQueryer) Ext() sqlx.ExtContext {
	return q.ext
}
func (q *animalQueryer) Save(ctx context.Context, animal *model.Animal) (int64, error) {
	return 0, nil
}

func (q *animalQueryer) GetAll(ctx context.Context) ([]model.Animal, error) {
	var animals []model.Animal
	if err := sqlx.SelectContext(ctx, q.ext, &animals, qGetAllAnimals); err != nil {
		return nil, err
	}
	return animals, nil
}

func (q *animalQueryer) FetchTopRelated(ctx context.Context, name string, lat, long float64) (model.Animal, error) {
	query, args, err := sqlx.Named(qFetchTopRelated, fiber.Map{"Name": name, "Latitude": lat, "Longitude": long})
	if err != nil {
		return model.Animal{}, err
	}
	var animal model.Animal
	if err := sqlx.GetContext(ctx, q.ext, &animal, query, args...); err != nil {
		return model.Animal{}, err
	}
	return animal, nil
}
