package postgres

import (
	"context"

	core "github.com/Qalifah/aboki-africa-assessment"
)

type ReferralRepository struct {
	client *Client
}

func NewReferralRepository(client *Client) *ReferralRepository {
	return &ReferralRepository{
		client: client,
	}
}

func (r *ReferralRepository) CreateReferral(ctx context.Context, referral *core.Referral) error {
	tx, err := r.client.GetTx(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, 
		"INSERT INTO referrals (referrer_id, referee_id) VALUES ($1, $2) RETURNING id", referral.ReferrerID, referral.RefereeID,
	)

	err = row.Scan(&referral.ID)

	return err
}