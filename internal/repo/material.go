package repo

import (
	"context"
	"strconv"

	"github.com/MXkodo/Management-of-School-materials/model"
	"github.com/jackc/pgx/v5"
)

type MaterialRepository interface {
	CreateMaterial(ctx context.Context, mat *model.Material) (string, error)
	GetMaterial(ctx context.Context, uuid string) (*model.Material, error)
	UpdateMaterial(ctx context.Context, mat *model.Material) error
	GetAllMaterials(ctx context.Context, filter model.MaterialFilter, limit, offset int) ([]model.Material, error)
}

type PgMaterialRepository struct {
	conn *pgx.Conn
}

func NewPgMaterialRepository(conn *pgx.Conn) *PgMaterialRepository {
	return &PgMaterialRepository{conn: conn}
}

func (r *PgMaterialRepository) CreateMaterial(ctx context.Context, mat *model.Material) (string, error) {
	var uuid string
	err := r.conn.QueryRow(ctx, `INSERT INTO materials (uuid, type, status, title, content, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING uuid`, mat.UUID, mat.Type, mat.Status, mat.Title, mat.Content, mat.CreatedAt, mat.UpdatedAt).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (r *PgMaterialRepository) GetMaterial(ctx context.Context, uuid string) (*model.Material, error) {
	var mat model.Material
	err := r.conn.QueryRow(ctx, `SELECT uuid, type, status, title, content, created_at, updated_at
	FROM materials WHERE uuid=$1`, uuid).Scan(&mat.UUID, &mat.Type, &mat.Status, &mat.Title, &mat.Content, &mat.CreatedAt, &mat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &mat, nil
}

func (r *PgMaterialRepository) UpdateMaterial(ctx context.Context, mat *model.Material) error {
	_, err := r.conn.Exec(ctx, `UPDATE materials SET status=$1, title=$2, content=$3, updated_at=$4 
	WHERE uuid=$5`, mat.Status, mat.Title, mat.Content, mat.UpdatedAt, mat.UUID)
	return err
}

func (r *PgMaterialRepository) GetAllMaterials(ctx context.Context, filter model.MaterialFilter, limit, offset int) ([]model.Material, error) {
	query := `SELECT uuid, type, status, title, content, created_at, updated_at FROM materials`
	params := []interface{}{}
	paramIndex := 1

	if filter.Status != "" || filter.Type != "" || !filter.DateFrom.IsZero() || !filter.DateTo.IsZero() {
		query += ` WHERE`

		if filter.Status != "" {
			query += ` status=$` + strconv.Itoa(paramIndex)
			params = append(params, filter.Status)
			paramIndex++
		}

		if filter.Type != "" {
			if len(params) > 0 {
				query += ` AND`
			}
			query += ` type=$` + strconv.Itoa(paramIndex)
			params = append(params, filter.Type)
			paramIndex++
		}

		if !filter.DateFrom.IsZero() {
			if len(params) > 0 {
				query += ` AND`
			}
			query += ` created_at >= $` + strconv.Itoa(paramIndex)
			params = append(params, filter.DateFrom)
			paramIndex++
		}

		if !filter.DateTo.IsZero() {
			if len(params) > 0 {
				query += ` AND`
			}
			query += ` created_at <= $` + strconv.Itoa(paramIndex)
			params = append(params, filter.DateTo)
			paramIndex++
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)
	params = append(params, limit, offset)

	rows, err := r.conn.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []model.Material
	for rows.Next() {
		var mat model.Material
		if err := rows.Scan(&mat.UUID, &mat.Type, &mat.Status, &mat.Title, &mat.Content, &mat.CreatedAt, &mat.UpdatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, mat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}
