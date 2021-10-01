package postgres

import (
	"context"

	core "github.com/Qalifah/aboki-africa-assessment"
)

type TransactionRepository struct {
	client *Client
}

func NewTransactionRepository(client *Client) *TransactionRepository {
	return &TransactionRepository{
		client: client,
	}
}

func(t *TransactionRepository) CreateTransaction(ctx context.Context, transaction *core.Transaction) error {
	tx, err := t.client.GetTx(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, 
		"INSERT INTO transactions (sender_id, recipient_id, points) VALUES ($1, $2, $3, $4) RETURNING id", 
		transaction.SenderID, transaction.RecipientID, transaction.Points, transaction.Type,
	)

	err = row.Scan(&transaction.ID)
	
	return err
}