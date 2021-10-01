package postgres

import (
	"context"

	core "github.com/Qalifah/aboki-africa-assessment"
)

type ReferralCodeRepository struct {
	client *Client
}

func NewReferralCodeRepository(client *Client) *ReferralCodeRepository {
	return &ReferralCodeRepository{
		client: client,
	}
}

func(rc *ReferralCodeRepository) CreateReferralCode(ctx context.Context, uRefCode *core.ReferralCode) error {
	tx, err := rc.client.GetTx(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, 
		"INSERT INTO referral_codes (user_id, code) VALUES ($1, $2) RETURNING id", uRefCode.UserID, uRefCode.Code,
	)

	err = row.Scan(&uRefCode.ID)

	return err
}

func(rc *ReferralCodeRepository) FindReferralCodeByUserID(ctx context.Context, userID string) (*core.ReferralCode, error) {
	tx, err := rc.client.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, "SELECT * FROM referral_codes WHERE user_id = $1 AND deleted_at IS NULL", userID)

	rCode := &core.ReferralCode{}
	err = row.Scan(&rCode.ID, &rCode.UserID, &rCode.Code, &rCode.CreatedAt)
	if err != nil {
		return nil, err
	}

	return rCode, nil
}