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

package pkg

import (
	"context"
	"time"
)

// User user
type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

// UserProvider user provider service
type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
}

// Reference rpc service id or reference id
func (u *UserProvider) Reference() string {
	return "UserProvider"
}

// JavaClassName got a go struct's Java Class package name which should be a POJO class
func (User) JavaClassName() string {
	return "org.apache.dubbo.User"
}
