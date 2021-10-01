package handler

type UserRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	ReferralCode *string `json:"referral_code"`
}

type TransferPointsRequest struct {
	SenderID          string `json:"sender_id"`
	RecipientID 	  string `json:"recipient_id"`
	Points            int    `json:"points"`
}