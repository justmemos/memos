package postgres

import (
	"context"
	"strings"

	storepb "github.com/usememos/memos/proto/gen/store"
	"github.com/usememos/memos/store"
)

func (d *DB) CreateReaction(ctx context.Context, create *storepb.Reaction) (*storepb.Reaction, error) {
	fields := []string{"creator_id", "content_id", "reaction_type"}
	args := []interface{}{create.CreatorId, create.ContentId, create.ReactionType.String()}
	stmt := "INSERT INTO reaction (" + strings.Join(fields, ", ") + ") VALUES (" + placeholders(len(args)) + ") RETURNING id, created_ts"
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.Id,
		&create.CreatedTs,
	); err != nil {
		return nil, err
	}

	reaction := create
	return reaction, nil
}

func (d *DB) ListReactions(ctx context.Context, find *store.FindReaction) ([]*storepb.Reaction, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if find.ID != nil {
		where, args = append(where, "id = "+placeholder(len(args)+1)), append(args, *find.ID)
	}
	if find.CreatorID != nil {
		where, args = append(where, "creator_id = "+placeholder(len(args)+1)), append(args, *find.CreatorID)
	}
	if find.ContentID != nil {
		where, args = append(where, "content_id = "+placeholder(len(args)+1)), append(args, *find.ContentID)
	}

	rows, err := d.db.QueryContext(ctx, `
		SELECT
			id,
			created_ts,
			creator_id,
			content_id,
			reaction_type
		FROM reaction
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*storepb.Reaction{}
	for rows.Next() {
		reaction := &storepb.Reaction{}
		var reactionType string
		if err := rows.Scan(
			&reaction.Id,
			&reaction.CreatedTs,
			&reaction.CreatorId,
			&reaction.ContentId,
			&reactionType,
		); err != nil {
			return nil, err
		}
		reaction.ReactionType = storepb.Reaction_Type(storepb.Reaction_Type_value[reactionType])
		list = append(list, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteReaction(ctx context.Context, delete *store.DeleteReaction) error {
	_, err := d.db.ExecContext(ctx, "DELETE FROM reaction WHERE id = $1", delete.ID)
	return err
}
