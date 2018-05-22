package db

import (
	"context"
	"database/sql/driver"
	"errors"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

type wrappedDriver struct {
	parent driver.Driver
	logger zerolog.Logger
}

type wrappedConn struct {
	parent driver.Conn
	logger zerolog.Logger
}

type wrappedTx struct {
	ctx    context.Context
	parent driver.Tx
	start  time.Time
	logger zerolog.Logger
}

type wrappedStmt struct {
	ctx    context.Context
	query  string
	parent driver.Stmt
	logger zerolog.Logger
}

type wrappedResult struct {
	ctx    context.Context
	parent driver.Result
}

type wrappedRows struct {
	ctx    context.Context
	parent driver.Rows
}

// WrapDriver -
func WrapDriver(driver driver.Driver, l zerolog.Logger) driver.Driver {
	return wrappedDriver{parent: driver, logger: l}
}

func (d wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.parent.Open(name)
	if err != nil {
		return nil, err
	}

	return wrappedConn{parent: conn, logger: d.logger}, nil
}

func (c wrappedConn) Prepare(query string) (driver.Stmt, error) {
	parent, err := c.parent.Prepare(query)
	if err != nil {
		return nil, err
	}

	return wrappedStmt{query: query, parent: parent, logger: c.logger}, nil
}

func (c wrappedConn) Close() error {
	return c.parent.Close()
}

func (c wrappedConn) Begin() (driver.Tx, error) {
	tx, err := c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{parent: tx}, nil
}

func (c wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	_, file, line, _ := runtime.Caller(10)
	start := time.Now()

	c.logger.Info().Msgf("[WRAPPED] Begin Tx [%s]: %d, %s", file, line, start)

	if connBeginTx, ok := c.parent.(driver.ConnBeginTx); ok {
		tx, err = connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		return wrappedTx{ctx: ctx, parent: tx, logger: c.logger, start: start}, nil
	}

	tx, err = c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{ctx: ctx, parent: tx, logger: c.logger, start: start}, nil
}

func (c wrappedConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if connPrepareCtx, ok := c.parent.(driver.ConnPrepareContext); ok {
		stmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		return wrappedStmt{ctx: ctx, parent: stmt, logger: c.logger}, nil
	}

	return c.Prepare(query)
}

func (c wrappedConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	start := time.Now()
	c.logger.Info().Msgf("[WRAPPED] Exec [%s]: %s", start, query)

	defer func() {
		c.logger.Info().Msgf("[WRAPPED] Exec Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, query)
	}()

	if execer, ok := c.parent.(driver.Execer); ok {
		res, err := execer.Exec(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{parent: res}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	start := time.Now()
	c.logger.Info().Msgf("[WRAPPED] ExecContext [%s]: %s", start, query)

	defer func() {
		c.logger.Info().Msgf("[WRAPPED] ExecContext Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, query)
	}()

	if execContext, ok := c.parent.(driver.ExecerContext); ok {
		res, err := execContext.ExecContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Exec(query, dargs)
}

func (c wrappedConn) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.parent.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}

	return nil
}

func (c wrappedConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	start := time.Now()
	c.logger.Info().Msgf("[WRAPPED] Query [%s]: %s", start, query)
	c.logger.Info().Msgf("[WRAPPED] Args %v", args)
	defer func() {
		c.logger.Info().Msgf("[WRAPPED] Query Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, query)
	}()

	if queryer, ok := c.parent.(driver.Queryer); ok {
		rows, err := queryer.Query(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{parent: rows}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	start := time.Now()

	ignored := ctx.Value("skip-logging")

	if ignored == nil {
		c.logger.Info().Msgf("[WRAPPED] QueryContext [%s]: %s", start, query)
		c.logger.Info().Msgf("[WRAPPED] Args %v", args)
	}

	defer func() {
		if ignored == nil {
			c.logger.Info().Msgf("[WRAPPED] QueryContext Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, query)
		}
	}()

	if queryerContext, ok := c.parent.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Query(query, dargs)
}

func (t wrappedTx) Commit() (err error) {
	_, file, line, _ := runtime.Caller(5)
	t.logger.Info().Msgf("[WRAPPED] Commit Tx [Took %fms] %s - %d", time.Since(t.start).Seconds()*1000, file, line)

	return t.parent.Commit()
}

func (t wrappedTx) Rollback() (err error) {
	_, file, line, _ := runtime.Caller(6)
	t.logger.Info().Msgf("[WRAPPED] Rollback Tx [Took %fms] %s - %d", time.Since(t.start).Seconds()*1000, file, line)

	return t.parent.Rollback()
}

func (s wrappedStmt) Close() (err error) {
	return s.parent.Close()
}

func (s wrappedStmt) NumInput() int {
	return s.parent.NumInput()
}

func (s wrappedStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	start := time.Now()
	s.logger.Info().Msgf("[WRAPPED] Exec [%s]: %s", start, s.query)

	defer func() {
		s.logger.Info().Msgf("[WRAPPED] Exec Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, s.query)
	}()

	res, err = s.parent.Exec(args)
	if err != nil {
		return nil, err
	}

	return wrappedResult{ctx: s.ctx, parent: res}, nil
}

func (s wrappedStmt) Query(args []driver.Value) (rows driver.Rows, err error) {
	start := time.Now()
	s.logger.Info().Msgf("[WRAPPED] Query [%s]: %s", start, s.query)
	s.logger.Info().Msgf("[WRAPPED] Args %v", args)
	defer func() {
		s.logger.Info().Msgf("[WRAPPED] Query Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, s.query)
	}()

	rows, err = s.parent.Query(args)
	if err != nil {
		return nil, err
	}

	return wrappedRows{ctx: s.ctx, parent: rows}, nil
}

func (s wrappedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	start := time.Now()
	s.logger.Info().Msgf("[WRAPPED] ExecContext [%s]: %s", start, s.query)

	for _, arg := range args {
		s.logger.Info().Msgf("[WRAPPED] Name: %s, Ordinal: %d, Value: %v", arg.Name, arg.Ordinal, arg.Value)
	}

	defer func() {
		s.logger.Info().Msgf("[WRAPPED] ExecContext Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, s.query)
	}()

	if stmtExecContext, ok := s.parent.(driver.StmtExecContext); ok {
		res, err := stmtExecContext.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Exec(dargs)
}

func (s wrappedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	start := time.Now()
	s.logger.Info().Msgf("[WRAPPED] QueryContext [%s]: %s", start, s.query)

	for _, arg := range args {
		s.logger.Info().Msgf("[WRAPPED] Name: %s, Ordinal: %d, Value: %v", arg.Name, arg.Ordinal, arg.Value)
	}

	defer func() {
		s.logger.Info().Msgf("[WRAPPED] QueryContext Finished [Took %fms]: %s", time.Since(start).Seconds()*1000, s.query)
	}()

	if stmtQueryContext, ok := s.parent.(driver.StmtQueryContext); ok {
		rows, err := stmtQueryContext.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Query(dargs)
}

func (r wrappedResult) LastInsertId() (id int64, err error) {
	return r.parent.LastInsertId()
}

func (r wrappedResult) RowsAffected() (num int64, err error) {
	return r.parent.RowsAffected()
}

func (r wrappedRows) Columns() []string {
	return r.parent.Columns()
}

func (r wrappedRows) Close() error {
	return r.parent.Close()
}

func (r wrappedRows) Next(dest []driver.Value) (err error) {
	return r.parent.Next(dest)
}

// namedValueToValue is a helper function copied from the database/sql package
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
