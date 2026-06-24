package mutex_test

import (
	"errors"
	"math"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/mutex"
)

func TestAccount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		initial int64
		f       func(*mutex.Account) error
		want    int64
		wantErr error
	}{
		{
			name:    "deposit",
			initial: 100,
			f:       func(a *mutex.Account) error { return a.Deposit(50) },
			want:    150,
			wantErr: nil,
		},
		{
			name:    "deposit non-positive",
			initial: 100,
			f:       func(a *mutex.Account) error { return a.Deposit(0) },
			want:    100,
			wantErr: mutex.ErrInvalidAmount,
		},
		{
			name:    "deposit overflow",
			initial: math.MaxInt64,
			f:       func(a *mutex.Account) error { return a.Deposit(1) },
			want:    math.MaxInt64,
			wantErr: mutex.ErrOverflow,
		},
		{
			name:    "withdraw",
			initial: 100,
			f:       func(a *mutex.Account) error { return a.Withdraw(40) },
			want:    60,
			wantErr: nil,
		},
		{
			name:    "withdraw non-positive",
			initial: 100,
			f:       func(a *mutex.Account) error { return a.Withdraw(-5) },
			want:    100,
			wantErr: mutex.ErrInvalidAmount,
		},
		{
			name:    "withdraw insufficient",
			initial: 100,
			f:       func(a *mutex.Account) error { return a.Withdraw(200) },
			want:    100,
			wantErr: mutex.ErrInsufficientFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := mutex.NewAccount(1, tt.initial)
			if err := tt.f(a); !errors.Is(err, tt.wantErr) {
				t.Errorf("%s error = %v, want %v", tt.name, err, tt.wantErr)
			}
			if got := a.Balance(); got != tt.want {
				t.Errorf("Balance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		fromID      uint64
		toID        uint64
		fromBalance int64
		toBalance   int64
		amount      int64
		wantFrom    int64
		wantTo      int64
		wantErr     error
	}{
		{
			name:        "success",
			fromID:      1,
			toID:        2,
			fromBalance: 100,
			toBalance:   50,
			amount:      30,
			wantFrom:    70,
			wantTo:      80,
			wantErr:     nil,
		},
		{
			name:        "success reversed ids",
			fromID:      2,
			toID:        1,
			fromBalance: 100,
			toBalance:   50,
			amount:      30,
			wantFrom:    70,
			wantTo:      80,
			wantErr:     nil,
		},
		{
			name:        "non-positive amount",
			fromID:      1,
			toID:        2,
			fromBalance: 100,
			toBalance:   50,
			amount:      0,
			wantFrom:    100,
			wantTo:      50,
			wantErr:     mutex.ErrInvalidAmount,
		},
		{
			name:        "insufficient funds",
			fromID:      1,
			toID:        2,
			fromBalance: 100,
			toBalance:   50,
			amount:      200,
			wantFrom:    100,
			wantTo:      50,
			wantErr:     mutex.ErrInsufficientFunds,
		},
		{
			name:        "recipient overflow",
			fromID:      1,
			toID:        2,
			fromBalance: 100,
			toBalance:   math.MaxInt64,
			amount:      1,
			wantFrom:    100,
			wantTo:      math.MaxInt64,
			wantErr:     mutex.ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			from := mutex.NewAccount(tt.fromID, tt.fromBalance)
			to := mutex.NewAccount(tt.toID, tt.toBalance)
			if err := mutex.Transfer(from, to, tt.amount); !errors.Is(err, tt.wantErr) {
				t.Errorf("Transfer(from, to, %d) error = %v, want %v", tt.amount, err, tt.wantErr)
			}
			if got := from.Balance(); got != tt.wantFrom {
				t.Errorf("from.Balance() = %v, want %v", got, tt.wantFrom)
			}
			if got := to.Balance(); got != tt.wantTo {
				t.Errorf("to.Balance() = %v, want %v", got, tt.wantTo)
			}
		})
	}

	t.Run("self transfer", func(t *testing.T) {
		t.Parallel()
		a := mutex.NewAccount(1, 100)
		if err := mutex.Transfer(a, a, 10); !errors.Is(err, mutex.ErrSelfTransfer) {
			t.Errorf("Transfer(a, a, 10) error = %v, want %v", err, mutex.ErrSelfTransfer)
		}
		if got := a.Balance(); got != 100 {
			t.Errorf("Balance() = %v, want %v", got, 100)
		}
	})
}
