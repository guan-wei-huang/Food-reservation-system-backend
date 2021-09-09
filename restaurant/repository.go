package restaurant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var (
	ErrRestaurantIdWrong = errors.New("restaurant not found")
)

type Repository interface {
	Close()
	CreateRestaurant(ctx context.Context, rest *Restaurant) (int, error)
	CreateFood(ctx context.Context, f *Food) error
	GetMenu(ctx context.Context, rid int) (*Menu, error)
	SearchRestaurant(ctx context.Context, latitude, longitude float64) ([]*Restaurant, error)
}

type repository struct {
	db *sql.DB
}

func NewRestaurantRepository(dsn string) (Repository, error) {
	db, err := sql.Open("postgre", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) Close() {
	r.db.Close()
}

func (r *repository) CreateRestaurant(ctx context.Context, rest *Restaurant) (int, error) {
	stmt := `INSERT INTO restaurant (name, description, location, coordinate) 
		VALUES ($1, $2, $3, ST_GeomFromText($4, 4326));`

	point := fmt.Sprintf("Point(%f %f)", rest.Latitude, rest.Longitude)
	result, err := r.db.ExecContext(ctx, stmt, rest.Name, rest.Description, rest.Location, point)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (r *repository) CreateFood(ctx context.Context, f *Food) error {
	stmt := `INSERT INTO foods (rid, name, description, price)
	VALUES ($1, $2, $3, $4);`

	_, err := r.db.ExecContext(ctx, stmt, f.Rid, f.Name, f.Description, f.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetMenu(ctx context.Context, rid int) (*Menu, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT (fid, rid, name, description, price) 
		FROM foods 
		WHERE rid = $1
		ORDER BY foods.fid;`,
		rid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	foods := []Food{}
	for rows.Next() {
		f := &Food{}
		err = rows.Scan(&f.Fid, &f.Rid, &f.Name, &f.Description, &f.Price)
		if err != nil {
			return nil, err
		}

		foods = append(foods, *f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &Menu{rid, foods}, nil
}

func (r *repository) SearchRestaurant(ctx context.Context, latitude, longtitude float64) ([]*Restaurant, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, description, location, ST_AsEWKT(coordinate)
		FROM restaurant
		WHERE ST_DWithin(coordinate, ST_GeogFromText('POINT($1, $2)'), 5000)
		ORDER BY id;`,
		latitude,
		longtitude,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := []*Restaurant{}
	for rows.Next() {
		r := &Restaurant{}
		point := ewkb.Point{}
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.Location, &point)
		if err != nil {
			return nil, err
		}

		geometry, err := geojson.Marshal(point.Point)
		if err != nil {
			return nil, err
		}

		// TODO: fixed this
		log.Println(string(geometry))

		restaurants = append(restaurants, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return restaurants, nil
}
