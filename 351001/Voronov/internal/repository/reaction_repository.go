package repository

import (
	apperrors "Voronov/internal/errors"
	"Voronov/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgReactionRepository struct {
	pool *pgxpool.Pool
}

func NewReactionRepository(pool *pgxpool.Pool) ReactionRepository {
	return &pgReactionRepository{pool: pool}
}

func (r *pgReactionRepository) FindByID(ctx context.Context, id int64) (*model.Reaction, error) {
	var rx model.Reaction
	err := r.pool.QueryRow(ctx,
		"SELECT id, issue_id, content FROM distcomp.tbl_reaction WHERE id = $1", id,
	).Scan(&rx.ID, &rx.IssueID, &rx.Content)
	if err == pgx.ErrNoRows {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rx, nil
}

func (r *pgReactionRepository) FindAll(ctx context.Context, opts *QueryOptions) ([]*model.Reaction, int64, error) {
	if opts == nil {
		opts = NewQueryOptions()
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM distcomp.tbl_reaction").Scan(&total); err != nil {
		return nil, 0, err
	}

	orderField := "id"
	orderDir := "ASC"
	if opts.Sort != nil {
		if opts.Sort.Field != "" {
			orderField = opts.Sort.Field
		}
		if opts.Sort.Direction == "DESC" {
			orderDir = "DESC"
		}
	}

	offset := (opts.Pagination.Page - 1) * opts.Pagination.PageSize
	query := fmt.Sprintf(
		"SELECT id, issue_id, content FROM distcomp.tbl_reaction ORDER BY %s %s LIMIT $1 OFFSET $2",
		orderField, orderDir,
	)

	rows, err := r.pool.Query(ctx, query, opts.Pagination.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*model.Reaction
	for rows.Next() {
		var rx model.Reaction
		if err := rows.Scan(&rx.ID, &rx.IssueID, &rx.Content); err != nil {
			return nil, 0, err
		}
		items = append(items, &rx)
	}
	if items == nil {
		items = []*model.Reaction{}
	}
	return items, total, nil
}

func (r *pgReactionRepository) Create(ctx context.Context, reaction *model.Reaction) (*model.Reaction, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO distcomp.tbl_reaction (issue_id, content) VALUES ($1, $2) RETURNING id",
		reaction.IssueID, reaction.Content,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, apperrors.ErrDuplicate
		}
		return nil, apperrors.FromDBError(err)
	}
	reaction.ID = id
	return reaction, nil
}

func (r *pgReactionRepository) Update(ctx context.Context, id int64, reaction *model.Reaction) (*model.Reaction, error) {
	result, err := r.pool.Exec(ctx,
		"UPDATE distcomp.tbl_reaction SET issue_id = $1, content = $2 WHERE id = $3",
		reaction.IssueID, reaction.Content, id,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, apperrors.ErrDuplicate
		}
		return nil, apperrors.FromDBError(err)
	}
	if result.RowsAffected() == 0 {
		return nil, apperrors.ErrNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *pgReactionRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM distcomp.tbl_reaction WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *pgReactionRepository) FindByIssueID(ctx context.Context, issueID int64) ([]*model.Reaction, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, issue_id, content FROM distcomp.tbl_reaction WHERE issue_id = $1", issueID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.Reaction
	for rows.Next() {
		var rx model.Reaction
		if err := rows.Scan(&rx.ID, &rx.IssueID, &rx.Content); err != nil {
			return nil, err
		}
		items = append(items, &rx)
	}
	if items == nil {
		items = []*model.Reaction{}
	}
	return items, nil
}
