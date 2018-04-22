package axmysql

import "github.com/go-sql-driver/mysql"

const (
	mysqlDeadLock  = 1213 // https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html#error_er_lock_deadlock
	mysqlDuplicate = 1062 // https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html#error_er_dup_entry
)

// isMySQLError returns true if err is a MySQL error with the given code.
func isMySQLError(err error, code uint16) bool {
	e, ok := err.(*mysql.MySQLError)
	return ok && e.Number == code
}

// isDuplicate returns true if err is a MySQL error indicating a duplicate key.
func isDuplicateKey(err error) bool {
	return isMySQLError(err, mysqlDuplicate)
}
