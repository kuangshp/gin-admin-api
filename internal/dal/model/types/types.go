package types

import (
	"github.com/shopspring/decimal"
	"gorm.io/plugin/optimisticlock"
)

type Decimal = decimal.Decimal

// Version mysql乐观锁
type Version = optimisticlock.Version
