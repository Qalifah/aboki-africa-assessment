package postgres

import (
	"context"

	core "github.com/Qalifah/aboki-africa-assessment"
)

type PointRepository struct {
	client *Client
}

func NewPointRepository(client *Client) *PointRepository {
	return &PointRepository{
		client: client,
	}
}

func(p *PointRepository) CreatePoint(ctx context.Context, point *core.Point) error {
	tx, err := p.client.GetTx(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, 
		"INSERT INTO user_points (user_id) VALUES ($1) RETURNING id", point.UserID,
	)

	err = row.Scan(&point.ID)

	return err
}

func(p *PointRepository) FindPointByUserID(ctx context.Context, userID string) (*core.Point, error) {
	tx, err := p.client.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, "SELECT * FROM user_points WHERE user_id = $1 AND deleted_at IS NULL", userID)

	point := &core.Point{}
	err = row.Scan(&point.ID, &point.UserID, &point.Points, &point.NumberOfReferredUsers, &point.Bonus, &point.Paid, &point.CreatedAt, &point.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return point, nil
}

func(p *PointRepository) UpdatePoint(ctx context.Context, point *core.Point) error {
	tx, err := p.client.GetTx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE user_points SET updated_at = CURRENT_TIMESTAMP AND points = $1 AND bonus = $2 AND number_of_referred_users = $3 WHERE id = $4 AND deleted_at IS NULL",
		point.Points, point.Bonus, point.NumberOfReferredUsers, point.ID,
	)

	return err
}

func (u *PointRepository) GetPointsBalance(ctx context.Context, userID string) (int, error) {
	tx, err := u.client.GetTx(ctx)
	if err != nil {
		return 0, err
	}

	var balance int
	row := tx.QueryRow(ctx, "SELECT points from user_points WHERE user_id = $1 AND deleted_at IS NULL", userID)
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}