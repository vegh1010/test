package modelinit

import (
	"github.com/jmoiron/sqlx"
	"github.com/vegh1010/test/pkg/model/merchant"
)

// PrepareStatements prepares all of the model's statements.
func PrepareStatements(db *sqlx.DB) {

	merchant.PrepareStatements(db)

}
