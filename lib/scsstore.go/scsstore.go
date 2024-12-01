package scsstore

import (
	"fmt"
	"time"
	"vshop/lib/model"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

func NewSCSStore(m *model.Model) (*scs.SessionManager, error) {

	// Check that database can be connected
	// with the raw database connection (sql.DB not sqlx.DB)

	if err := m.DbHandle.DB.Ping(); err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	scs := scs.New()
	scs.Store = mysqlstore.New(m.DbHandle.DB)
	scs.Lifetime = 12 * time.Hour

	return scs, nil
}
