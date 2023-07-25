package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type FindMemoCommentMessage struct {
	ID int

	// Domain specific fields
	ContentSearch []string

	MemoID *int

	// Pagination
	Limit            *int
	Offset           *int
	OrderByUpdatedTs bool
}

type MemoCommentMessage struct {
	ID int

	// Standard fields
	CreatedTs int64
	UpdatedTs int64

	// Domain specific fields
	Content string

	// Info fields
	Username string

	MemoID int
}

type DeleteMemoCommentMessage struct {
	ID     int
	MemoID int
}

func (s *Store) DeleteMemoComment(ctx context.Context, delete *DeleteMemoCommentMessage) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	where, args := []string{"id = ?", "memo_id = ?"}, []any{delete.ID, delete.MemoID}
	stmt := `DELETE FROM memo_comment WHERE ` + strings.Join(where, " AND ")
	result, err := tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if err := s.vacuumImpl(ctx, tx); err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (s *Store) CreateCommentMemo(ctx context.Context, create *MemoCommentMessage) (*MemoCommentMessage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if create.CreatedTs == 0 {
		create.CreatedTs = time.Now().Unix()
	}

	query := `
		INSERT INTO memo_comment (
			created_ts,
			content,
		    username,
			memo_id,
		)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_ts, updated_ts
	`
	if err := tx.QueryRowContext(
		ctx,
		query,
		create.CreatedTs,
		create.Content,
		create.Username,
		create.MemoID,
	).Scan(
		&create.ID,
		&create.CreatedTs,
		&create.UpdatedTs,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return create, nil
}

func (s *Store) GetMemoComments(ctx context.Context, find *FindMemoCommentMessage) ([]*MemoCommentMessage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	list, err := listMemoComments(ctx, tx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, err
	}

	return list, nil
}

func listMemoComments(ctx context.Context, tx *sql.Tx, find *FindMemoCommentMessage) ([]*MemoCommentMessage, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.MemoID; v != nil {
		where, args = append(where, "memo_comment.memo_id = ?"), append(args, *v)
	}

	if v := find.ContentSearch; len(v) != 0 {
		for _, s := range v {
			where, args = append(where, "memo_comment.content LIKE ?"), append(args, "%"+s+"%")
		}
	}

	query := `
	SELECT
		memo_comment.id AS id,
		memo_comment.created_ts AS created_ts,
		memo_comment.updated_ts AS updated_ts,
		memo_comment.content AS content,
		memo_comment.username AS username,
		memo_comment.memo_id AS memo_id
	FROM
		memo_comment 
	WHERE ` + strings.Join(where, " AND ") + `
	`
	if find.Limit != nil {
		query = fmt.Sprintf("%s LIMIT %d", query, *find.Limit)
		if find.Offset != nil {
			query = fmt.Sprintf("%s OFFSET %d", query, *find.Offset)
		}
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memoCommentMessageList := make([]*MemoCommentMessage, 0)
	for rows.Next() {
		var memoCommentMessage MemoCommentMessage
		if err := rows.Scan(
			&memoCommentMessage.ID,
			&memoCommentMessage.CreatedTs,
			&memoCommentMessage.UpdatedTs,
			&memoCommentMessage.Content,
			&memoCommentMessage.Username,
			&memoCommentMessage.MemoID,
		); err != nil {
			return nil, err
		}
		memoCommentMessageList = append(memoCommentMessageList, &memoCommentMessage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memoCommentMessageList, nil
}
