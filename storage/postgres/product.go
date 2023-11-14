package postgres

import (
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"

	"app/models"
	"app/pkg/helper"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(req *models.ProductCreate) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO product(id, name, price, category_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(query,
		id,
		req.Name,
		req.Price,
		helper.NewNullString(req.CategoryId),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *productRepo) GetByID(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		resp  models.Product
		query string
	)

	query = `
		SELECT
			id, 
			name, 
			price, 
			COALESCE(category_id::VARCHAR, ''),
			created_at, 
			updated_at
		FROM product
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&resp.Id,
		&resp.Name,
		&resp.Price,
		&resp.CategoryId,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *productRepo) GetList(req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {

	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id, 
			name, 
			price, 
			category_id, 
			created_at, 
			updated_at
		FROM product
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND title ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			product models.Product
			// parentId sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&product.Id,
			&product.Name,
			&product.Price,
			&product.CategoryId,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// category.ParentID = parentId.String
		resp.Products = append(resp.Products, &product)
	}

	return resp, nil
}

func (c *productRepo) Update(req *models.ProductUpdate) (*models.ProductPrimaryKey, error) {
	var (
		query = `
			UPDATE product
			SET
			name = $1,
			price = $2,
			category_id = $3
			WHERE id = $4
		`
	)

	_, err := c.db.Exec(query, req.Name, req.Price, req.CategoryId, req.Id)
	if err != nil {
		return nil, err
	}
	return &models.ProductPrimaryKey{Id: req.Id}, nil
}

func (c *productRepo) Delete(req *models.ProductPrimaryKey) error {
	var (
		query = `DELETE FROM product WHERE id = $1`
	)

	_, err := c.db.Exec(query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
