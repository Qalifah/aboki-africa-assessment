package aboki_africa_assessment

import (
	"context"
	"time"
)

const bonus = 50

type User struct {
	ID 			string		`json:"id"`
	Name		string		`json:"name"`
	Email		string		`json:"email"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	time.Time	`json:"deleted_at"`
}

type ReferralCode struct {
	ID			string		`json:"id"`
	UserID		string		`json:"user_id"`
	Code		string		`json:"code"`
	CreatedAt   time.Time	`json:"created_at"`
	DeletedAt	time.Time	`json:"deleted_at"`
}

type Referral struct {
	ID			string		`json:"id"`
	ReferrerID  string      `json:"referrer_id"` 
	RefereeID   string      `json:"referee_id"` 
	CreatedAt   time.Time	`json:"created_at"`
	DeletedAt	time.Time	`json:"deleted_at"`
}

type Point struct {
	ID						string 		`json:"id"`
	UserID					string		`json:"user_id"`
	Points					int			`json:"points"`
	NumberOfReferredUsers	int			`json:"number_of_referred_users"`
	Bonus					int			`json:"bonus"`
	Paid					bool		`json:"paid"`
	CreatedAt   			time.Time	`json:"created_at"`
	UpdatedAt				time.Time	`json:"updated_at"`
	DeletedAt				time.Time	`json:"deleted_at"`
}

func(p *Point) Deduct(points int) {
	p.Points -= points
}

func(p *Point) Add(points int) {
	p.Points += points
}

func(p *Point) AddBonus() {
	p.Bonus += bonus
}

func(p *Point) IncreaseUserReferrals() {
	p.NumberOfReferredUsers++
}

type Transaction struct {
	ID              string     `json:"id"`
	SenderID        string     `json:"sender_id"`
	RecipientID 	string     `json:"recipient_id"`
	Points			int		   `json:"points"`
	Type			string	   `json:"type"`
	CreatedAt   	time.Time  `json:"created_at"`
	DeletedAt		time.Time  `json:"deleted_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByReferralCode(ctx context.Context, code string) (*User, error)
}

type ReferralCodeRepository interface {
	CreateReferralCode(ctx context.Context, uRefCode *ReferralCode) error
	FindReferralCodeByUserID(ctx context.Context, userID string) (*ReferralCode, error)
}

type ReferralRepository interface {
	CreateReferral(ctx context.Context, referral *Referral) error
}

type PointRepository interface {
	CreatePoint(ctx context.Context, Point *Point) error
	FindPointByUserID(ctx context.Context, userID string) (*Point, error)
	UpdatePoint(ctx context.Context, Point *Point) error
	GetPointsBalance(ctx context.Context, userID string) (int, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	// ClaimReferrerBonus(ctx context.Context, userID string) error
}