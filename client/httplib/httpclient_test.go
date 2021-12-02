// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httplib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/djimenez/iconv-go"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	Init(&Options{})
	code := m.Run()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	var body []byte
	_, err := New(&Options{}).Get("http://www.163.com").BindBytes(&body).doInternal()
	assert.Nil(t, err)
}

func TestHttpClientGet(t *testing.T) {
	body := &testJsonBody{
		A: "fdstt",
		B: 5323,
		C: 43.54,
	}
	b, err := json.Marshal(body)
	assert.Nil(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write(b)
		assert.Nil(t, err)
	}))
	defer s.Close()

	c := New(&Options{})

	var obj testJsonBody
	statusCode, err := c.Get(s.URL).BindJson(&obj).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, &obj)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestHttpClientPost(t *testing.T) {
	body := &testJsonBody{
		A: "fdstt",
		B: 5323,
		C: 43.54,
	}
	b, err := json.Marshal(body)
	assert.Nil(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, err = w.Write(b)
		assert.Nil(t, err)
	}))
	defer s.Close()

	c := New(&Options{})

	var obj testJsonBody
	statusCode, err := c.Post(s.URL).BindJson(&obj).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, &obj)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestInit(t *testing.T) {
	body := []byte("ok")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(body)
		assert.Nil(t, err)
	}))
	defer s.Close()

	defaultHttpClient = nil
	var res1 []byte
	Init(&Options{})
	statusCode, err := Get(s.URL).BindBytes(&res1).doInternal()
	assert.EqualValues(t, body, res1)
	assert.Equal(t, http.StatusOK, statusCode)

	var res2 []byte
	c := New(&Options{})
	statusCode, err = c.Get(s.URL).BindBytes(&res2).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, res2)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestCharset(t *testing.T) {
	body := []byte("{\"a\": \"成功\"}")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.URL.Query()["charset"]
		if c != nil || len(c) > 0 {
			w.Header().Set("Content-Type", "application/json;charset="+c[0])
			output, err := iconv.ConvertString(string(body), "utf-8", c[0])
			assert.Nil(t, err)
			_, err = w.Write([]byte(output))
			assert.Nil(t, err)
		} else {
			_, err := w.Write(body)
			assert.Nil(t, err)
		}
	}))
	defer s.Close()

	c := New(&Options{})

	var res0 string
	statusCode, err := c.Get(s.URL).BindString(&res0).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, string(body), res0)
	assert.Equal(t, http.StatusOK, statusCode)

	var res1 string
	statusCode, err = c.Get(s.URL + "?charset=gbk").BindString(&res1).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, string(body), res1)
	assert.Equal(t, http.StatusOK, statusCode)

	var res2 testJsonBody
	statusCode, err = c.Get(s.URL + "?charset=gbk").BindJson(&res2).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, "成功", res2.A)
	assert.Equal(t, http.StatusOK, statusCode)

}

func TestCharset_gbk(t *testing.T) {
}

func TestCharset_utf8(t *testing.T) {
}

func TestCharset_String2JsonErr(t *testing.T) {
}
