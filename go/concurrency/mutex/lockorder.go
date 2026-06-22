package mutex

import (
	"errors"
	"math"
	"sync"
)

var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrOverflow          = errors.New("balance overflow")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSelfTransfer      = errors.New("cannot transfer to the same account")
)

type Account struct {
	id      uint64
	mu      sync.Mutex
	balance int64
}

func NewAccount(id uint64, balance int64) *Account {
	return &Account{id: id, balance: balance}
}

func (a *Account) Balance() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func (a *Account) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance > math.MaxInt64-amount {
		return ErrOverflow
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	return nil
}

func Transfer(from, to *Account, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if from == to {
		return ErrSelfTransfer
	}
	first, second := from, to
	if from.id > to.id {
		first, second = to, from
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()
	if from.balance < amount {
		return ErrInsufficientFunds
	}
	if to.balance > math.MaxInt64-amount {
		return ErrOverflow
	}
	from.balance -= amount
	to.balance += amount
	return nil
}
