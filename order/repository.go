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

	row := tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (restaurant_id, user_id, created_at) 
		VALUES($1, $2, $3)
		RETURNING id;`,
		order.Rid,
		order.Uid,
		order.CreatedAt,
	)
	if err != nil {
		return
	}

	if err = row.Scan(&id); err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO order_products(order_id, fid, name, price, quantity) 
		VALUES ($1, $2, $3, $4, $5)`,
	)
	for _, p := range order.Products {
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
	row := r.db.QueryRowContext(
		ctx,
		`SELECT o.id, o.user_id, o.restaurant_id, o.created_at
		FROM orders AS o
		WHERE o.id = $1`,
		id,
	)
	err := row.Scan(&order.Id, &order.Uid, &order.Rid, &order.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrOrderInvalid
		default:
			return nil, err
		}
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.fid, o.name, o.price, o.quantity 
		FROM order_products as o 
		WHERE o.order_id = $1
		ORDER BY o.fid`,
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

	order.Products = products
	return order, nil
}

func (r *repository) GetOrderForUser(ctx context.Context, id int) (*[]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id, o.restaurant_id, o.user_id, o.created_at, op.fid, op.name, op.price, op.quantity
		FROM orders AS o 
		JOIN order_products AS op ON (o.id = op.order_id)
		WHERE o.user_id = $1
		ORDER BY o.id`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := Order{}
	orders := []Order{}
	lastOrder := Order{Id: -1}
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err = rows.Scan(&order.Id, &order.Rid, &order.Uid, &order.CreatedAt, &p.Fid,
			&p.Name, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}

		if lastOrder.Id != -1 && order.Id != lastOrder.Id {
			lastOrder.Products = products
			orders = append(orders, lastOrder)
			products = products[:0]
		}
		products = append(products, p)
		lastOrder = order
	}
	// last item
	if lastOrder.Id != -1 {
		order.Products = products
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &orders, nil
}
