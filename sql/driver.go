// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package sql

import (
	"context"
	"database/sql/driver"

	"github.com/SkyAPM/go2sky"
)

type DBType string

const (
	MYSQL DBType = "mysql"
	IPV4  DBType = "others"
)

// swDriver is a tracing wrapper for driver.Driver
type swDriver struct {
	driver driver.Driver
	tracer *go2sky.Tracer

	dbType DBType
}

func NewTracerDriver(driver driver.Driver, tracer *go2sky.Tracer, dbType DBType) driver.Driver {
	return &swDriver{
		driver: driver,
		tracer: tracer,
		dbType: dbType,
	}
}

func (d *swDriver) Open(name string) (driver.Conn, error) {
	addr := parseAddr(name, d.dbType)
	s, err := d.tracer.CreateExitSpan(context.Background(), genOpName(d.dbType, "open"), addr, emptyInjectFunc)
	if err != nil {
		return nil, err
	}
	defer s.End()
	c, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &conn{
		conn:   c,
		tracer: d.tracer,
		addr:   addr,
	}, nil
}
