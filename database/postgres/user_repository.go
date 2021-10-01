package postgres

import (
	"context"

	core "github.com/Qalifah/aboki-africa-assessment"
)

type UserRepository struct {
	client *Client
}

func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{client: client}
}

func(u *UserRepository) CreateUser(ctx context.Context, user *core.User) error {
	tx, err := u.client.GetTx(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, 
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email,
	)

	err = row.Scan(&user.ID)

	return err
}

func(u *UserRepository) FindUserByID(ctx context.Context, id string) (*core.User, error) {
	tx, err := u.client.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL", id)

	user := &core.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func(u *UserRepository) FindUserByReferralCode(ctx context.Context, code string) (*core.User, error) {
	tx, err := u.client.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, `SELECT users.id, users.name, users.email, users.created_at, users.updated_at FROM users 
	INNER JOIN referral_codes ON users.id = referral_codes.user_id WHERE referral_codes.code = $1 AND referral_codes.deleted_at IS NULL`, code)

	user := &core.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}