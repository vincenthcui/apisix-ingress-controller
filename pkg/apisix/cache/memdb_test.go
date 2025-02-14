// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

func TestMemDBCacheRoute(t *testing.T) {
	c, err := NewMemDBCache()
	assert.Nil(t, err, "NewMemDBCache")

	r1 := &v1.Route{
		Metadata: v1.Metadata{
			ID:   "1",
			Name: "abc",
		},
	}
	assert.Nil(t, c.InsertRoute(r1), "inserting route 1")

	r, err := c.GetRoute("1")
	assert.Nil(t, err)
	assert.Equal(t, r1, r)

	r2 := &v1.Route{
		Metadata: v1.Metadata{
			ID:   "2",
			Name: "def",
		},
	}
	r3 := &v1.Route{
		Metadata: v1.Metadata{
			ID:   "3",
			Name: "ghi",
		},
	}
	assert.Nil(t, c.InsertRoute(r2), "inserting route r2")
	assert.Nil(t, c.InsertRoute(r3), "inserting route r3")

	r, err = c.GetRoute("3")
	assert.Nil(t, err)
	assert.Equal(t, r3, r)

	assert.Nil(t, c.DeleteRoute(r3), "delete route r3")

	routes, err := c.ListRoutes()
	assert.Nil(t, err, "listing routes")

	if routes[0].Name > routes[1].Name {
		routes[0], routes[1] = routes[1], routes[0]
	}
	assert.Equal(t, routes[0], r1)
	assert.Equal(t, routes[1], r2)

	r4 := &v1.Route{
		Metadata: v1.Metadata{
			ID:   "4",
			Name: "name4",
		},
	}
	assert.Error(t, ErrNotFound, c.DeleteRoute(r4))
}

func TestMemDBCacheSSL(t *testing.T) {
	c, err := NewMemDBCache()
	assert.Nil(t, err, "NewMemDBCache")

	s1 := &v1.Ssl{
		ID: "abc",
	}
	assert.Nil(t, c.InsertSSL(s1), "inserting ssl 1")

	s, err := c.GetSSL("abc")
	assert.Nil(t, err)
	assert.Equal(t, s1, s)

	s2 := &v1.Ssl{
		ID: "def",
	}
	s3 := &v1.Ssl{
		ID: "ghi",
	}
	assert.Nil(t, c.InsertSSL(s2), "inserting ssl 2")
	assert.Nil(t, c.InsertSSL(s3), "inserting ssl 3")

	s, err = c.GetSSL("ghi")
	assert.Nil(t, err)
	assert.Equal(t, s3, s)

	assert.Nil(t, c.DeleteSSL(s3), "delete ssl 3")

	ssl, err := c.ListSSL()
	assert.Nil(t, err, "listing ssl")

	if ssl[0].ID > ssl[1].ID {
		ssl[0], ssl[1] = ssl[1], ssl[0]
	}
	assert.Equal(t, ssl[0], s1)
	assert.Equal(t, ssl[1], s2)

	s4 := &v1.Ssl{
		ID: "id4",
	}
	assert.Error(t, ErrNotFound, c.DeleteSSL(s4))
}

func TestMemDBCacheUpstream(t *testing.T) {
	c, err := NewMemDBCache()
	assert.Nil(t, err, "NewMemDBCache")

	u1 := &v1.Upstream{
		Metadata: v1.Metadata{
			ID:   "1",
			Name: "abc",
		},
	}
	err = c.InsertUpstream(u1)
	assert.Nil(t, err, "inserting upstream 1")

	u, err := c.GetUpstream("1")
	assert.Nil(t, err)
	assert.Equal(t, u1, u)

	u2 := &v1.Upstream{
		Metadata: v1.Metadata{
			Name: "def",
			ID:   "2",
		},
	}
	u3 := &v1.Upstream{
		Metadata: v1.Metadata{
			Name: "ghi",
			ID:   "3",
		},
	}
	assert.Nil(t, c.InsertUpstream(u2), "inserting upstream 2")
	assert.Nil(t, c.InsertUpstream(u3), "inserting upstream 3")

	u, err = c.GetUpstream("3")
	assert.Nil(t, err)
	assert.Equal(t, u3, u)

	assert.Nil(t, c.DeleteUpstream(u3), "delete upstream 3")

	upstreams, err := c.ListUpstreams()
	assert.Nil(t, err, "listing upstreams")

	if upstreams[0].Name > upstreams[1].Name {
		upstreams[0], upstreams[1] = upstreams[1], upstreams[0]
	}
	assert.Equal(t, upstreams[0], u1)
	assert.Equal(t, upstreams[1], u2)

	u4 := &v1.Upstream{
		Metadata: v1.Metadata{
			Name: "name4",
			ID:   "4",
		},
	}
	assert.Error(t, ErrNotFound, c.DeleteUpstream(u4))
}

func TestMemDBCacheReference(t *testing.T) {
	r := &v1.Route{
		Metadata: v1.Metadata{
			Name: "route",
			ID:   "1",
		},
		UpstreamId: "1",
	}
	u := &v1.Upstream{
		Metadata: v1.Metadata{
			ID:   "1",
			Name: "upstream",
		},
	}

	db, err := NewMemDBCache()
	assert.Nil(t, err, "NewMemDBCache")
	assert.Nil(t, db.InsertRoute(r))
	assert.Nil(t, db.InsertUpstream(u))

	assert.Error(t, ErrStillInUse, db.DeleteUpstream(u))
	assert.Nil(t, db.DeleteRoute(r))
	assert.Nil(t, db.DeleteUpstream(u))
}
