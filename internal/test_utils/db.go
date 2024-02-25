package test_utils

import (
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
)

type TxManager interface {
	db.TxManager
}

type Transactor interface {
	db.Transactor
}

type Client interface {
	db.Client
}
