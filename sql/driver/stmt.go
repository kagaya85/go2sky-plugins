//
// Copyright 2021 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package driver

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/SkyAPM/go2sky"
)

// stmt is a tracing wrapper for driver.Stmt
type stmt struct {
	stmt   driver.Stmt
	tracer *go2sky.Tracer

	opts *options
	// query defines the statement query
	query string
}

func (s *stmt) Close() error {
	return s.stmt.Close()
}

func (s *stmt) NumInput() int {
	return s.stmt.NumInput()
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.stmt.Exec(args)
}

// ExecContext implements driver.StmtExecContext ExecContext
// If the underlying Stmt does not implements
// driver.StmtExecContext interface, this method
// will use Exec instead.
func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	span, err := createSpan(ctx, s.tracer, s.opts, "execute")
	if err != nil {
		return nil, err
	}
	defer span.End()

	span.Tag(tagDbType, string(s.opts.dbType))
	span.Tag(tagDbInstance, s.opts.peer)
	if s.opts.reportQuery {
		span.Tag(tagDbStatement, s.query)
	}
	if s.opts.reportParam {
		span.Tag(tagDbSqlParameters, namedValueToValueString(args))
	}

	if execerContext, ok := s.stmt.(driver.StmtExecContext); ok {
		res, err := execerContext.ExecContext(ctx, args)
		if err != nil {
			span.Error(time.Now(), err.Error())
		}
		return res, err
	}

	values, err := namedValueToValue(args)
	if err != nil {
		span.Error(time.Now(), err.Error())
		span.End()
		return nil, err
	}
	res, err := s.Exec(values)
	if err != nil {
		span.Error(time.Now(), err.Error())
		return res, err
	}
	return res, err
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.stmt.Query(args)
}

// QueryContext implements driver.StmtQueryContext QueryContext
// If the underlying Stmt does not implements
// driver.StmtQueryContext interface, this method
// will use Query instead.
func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	span, err := createSpan(ctx, s.tracer, s.opts, "query")
	if err != nil {
		return nil, err
	}
	defer span.End()

	span.Tag(tagDbType, string(s.opts.dbType))
	span.Tag(tagDbInstance, s.opts.peer)
	if s.opts.reportQuery {
		span.Tag(tagDbStatement, s.query)
	}
	if s.opts.reportParam {
		span.Tag(tagDbSqlParameters, namedValueToValueString(args))
	}

	if queryer, ok := s.stmt.(driver.StmtQueryContext); ok {
		rows, err := queryer.QueryContext(ctx, args)
		if err != nil {
			span.Error(time.Now(), err.Error())
		}
		return rows, err
	}

	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	rows, err := s.Query(values)
	if err != nil {
		span.Error(time.Now(), err.Error())
	}
	return rows, err
}
