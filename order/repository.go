package order

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

var (
	ErrOrderInvalid = errors.New("order not found")
)

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, order *Order) (int, error)
	GetOrder(ctx context.Context, id int) (*Order, error)
	GetOrderForUser(ctx context.Context, id int) (*[]Order, error)
}

type repository struct {
	db *sql.DB
}

func NewOrderRepository(dsn string) (Repository, error) {
	log.Println("order repository dsn: ", dsn)
	db, err := sql.Open("postgres", dsn)
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

func (r *repository) CreateOrder(ctx context.Context, order *Order) (id int, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO order(rid, uid, created_at) VALUES($1, $2, $3)",
		order.Rid,
		order.Uid,
		order.CreatedAt,
	)
	if err != nil {
		return
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return
	}

	id = int(id64)
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO order_products(order_id, fid, name, price, quantity) 
		VALUES ($1, $2, $3, $4, $5)`,
	)
	for _, p := range *order.Products {
		_, err = stmt.ExecContext(ctx, id, p.Fid, p.Name, p.Price, p.Quantity)
		if err != nil {
			return
		}
	}
	stmt.Close()
	return
}

func (r *repository) GetOrder(ctx context.Context, id int) (*Order, error) {
	order := &Order{}
	row := r.db.QueryRowContext(ctx, "SELECT * FROM order WHERE id = $1", id)
	switch err := row.Scan(&order.Id, &order.Rid, &order.Uid, &order.CreatedAt); err {
	case sql.ErrNoRows:
		return nil, ErrOrderInvalid
	case nil:
		break
	default:
		return nil, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.fid, o.name, o.price, o.quantity 
		FROM order_products as o 
		WHERE o.order_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.Fid, &p.Name, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	order.Products = &products
	return order, nil
}

func (r *repository) GetOrderForUser(ctx context.Context, id int) (*[]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id, o.rid, o.uid, o.created_at, op.fid, op.name, op.price, op.quantity
		FROM order as o 
		JOIN order_products as op ON (o.id = op.order_id)
		WHERE o.uid = $1
		ORDER BY o.id`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := Order{}
	orders := []Order{}
	lastOrder := &Order{}
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err = rows.Scan(&order.Id, &order.Rid, &order.Uid, &order.CreatedAt, &p.Fid,
			&p.Name, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}

		if order.Id != lastOrder.Id && lastOrder != nil {
			lastOrder.Products = &products
			orders = append(orders, *lastOrder)
			products = products[:0]
		}
		products = append(products, p)
		lastOrder = &order
	}
	// last item
	if lastOrder != nil {
		order.Products = &products
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &orders, nil
}
