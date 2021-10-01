package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	core "github.com/Qalifah/aboki-africa-assessment"
	"github.com/Qalifah/aboki-africa-assessment/errors"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

const (
	Fail = "Transfer Failed"
	Success = "Transfer Successful"
	transfer = "TRANSFER"
	bonus = "BONUS"
)

type Handler struct {
	userRepository         	core.UserRepository
	referralCodeRepository	core.ReferralCodeRepository
	referralRepository 		core.ReferralRepository
	pointRepository    		core.PointRepository
	transactionRepository 	core.TransactionRepository
	beginTxFunc            func() (pgx.Tx, error)
}

func New(userRepository core.UserRepository, referralCodeRepository	core.ReferralCodeRepository, referralRepository core.ReferralRepository,
	pointRepository core.PointRepository, transactionRepository core.TransactionRepository, beginTxFunc func() (pgx.Tx, error)) *Handler {
		return &Handler{
			userRepository: userRepository,
			referralRepository: referralRepository,
			referralCodeRepository: referralCodeRepository,
			pointRepository: pointRepository,
			transactionRepository: transactionRepository,
			beginTxFunc: beginTxFunc,
		}
}

func(h *Handler) RegisterUser(ctx context.Context, input *UserRequest, logger *log.Entry) (*core.User, error) {
	tx, err := h.beginTxFunc()
	if err != nil {
		logger.WithError(err).Error("failed to start transaction")
		return nil, errors.ErrGeneric
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, core.TxContextKey, tx)
	user := &core.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = h.userRepository.CreateUser(ctx, user)
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return nil, errors.ErrCreateUserFailed
	}

	userRefCode := &core.ReferralCode{
		UserID: user.ID,
		Code: GenReferralCode(7),
	}

	err = h.referralCodeRepository.CreateReferralCode(ctx, userRefCode)
	if err != nil {
		logger.WithError(err).Error("failed to create user referral code")
		return nil, errors.ErrGeneric
	}

	userPoint := &core.Point{
		UserID: user.ID,
		Points: 0,
	}

	err = h.pointRepository.CreatePoint(ctx, userPoint)
	if err != nil {
		logger.WithError(err).Error("failed to create user point")
		return nil, errors.ErrGeneric
	}

	if input.ReferralCode != nil {
		referrer, err := h.userRepository.FindUserByReferralCode(ctx, *input.ReferralCode)
		if err != nil {
			logger.WithError(err).Error("failed to find user by referral code")
			return nil, errors.ErrGeneric
		}

		userReferral := &core.Referral{
			ReferrerID: referrer.ID,
			RefereeID: user.ID,
		}

		err = h.referralRepository.CreateReferral(ctx, userReferral)
		if err != nil {
			logger.WithError(err).Error("failed to create user referral")
			return nil, errors.ErrGeneric
		}

		refPoint, err := h.pointRepository.FindPointByUserID(ctx, referrer.ID)
		if err != nil {
			logger.WithError(err).Error("failed to find user point")
			return nil, errors.ErrGeneric
		}

		// increase user's referrals counter and check if the counter is divisible by 3 to add bonus
		refPoint.IncreaseUserReferrals()
		if (refPoint.NumberOfReferredUsers > 0) && (refPoint.NumberOfReferredUsers % 3 == 0) {
			refPoint.AddBonus()
		}

		err = h.pointRepository.UpdatePoint(ctx, refPoint)
		if err != nil {
			logger.WithError(err).Error("failed to update user point")
			return nil, errors.ErrGeneric
		}
	}

	if err = tx.Commit(ctx); err != nil {
		logger.WithError(err).Error("failed to commit transaction")
		return nil, errors.ErrGeneric
	}

	return user, nil
}

func(h *Handler) TransferPoints(ctx context.Context, input *TransferPointsRequest, logger *log.Entry) (string, error) {
	tx, err := h.beginTxFunc()
	if err != nil {
		logger.WithError(err).Error("failed to start transaction")
		return Fail, errors.ErrGeneric
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, core.TxContextKey, tx)
	balance, err := h.pointRepository.GetPointsBalance(ctx, input.SenderID)
	if err != nil {
		logger.WithError(err).Error("failed to get user points balance")
		return Fail, errors.ErrGeneric
	}

	point, err := h.pointRepository.FindPointByUserID(ctx, input.SenderID)
	if err != nil {
		logger.WithError(err).Error("failed to get sender points")
		return Fail, errors.ErrGeneric
	}

	if balance < input.Points {
		return Fail + getBonusBalanceStatement(point), errors.ErrInsufficientFunds
	}

	tran := &core.Transaction{
		SenderID: input.SenderID,
		RecipientID: input.RecipientID,
		Points: input.Points,
		Type: transfer,
	}

	err = h.transactionRepository.CreateTransaction(ctx, tran)
	if err != nil {
		logger.WithError(err).Error("failed transaction")
		return Fail, errors.ErrTransactionFailed
	}

	return Success + getBonusBalanceStatement(point), nil

}

func getBonusBalanceStatement(point *core.Point) string {
	if point.Bonus == 0 {
		return ""
	}
	return fmt.Sprintf(" : You have an unclaimed bonus of %v points", point.Bonus)
}

// GenReferralCode helps to generate reference code.
func GenReferralCode(max int) string {
	table := []byte("abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}