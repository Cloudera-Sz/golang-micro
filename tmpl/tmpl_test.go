// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package tmpl

import (
	"testing"
)

func TestRender(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	var tests = []struct {
		tmpl, expect string
		data         interface{}
	}{
		{
			tmpl:   `Hello, {{.Name}}`,
			expect: `Hello, Neo`,
			data:   map[string]string{"Name": "Neo"},
		},
		{
			tmpl:   `Hello, {{.Name}}, Age = {{.Age}}`,
			expect: `Hello, Neo, Age = 21`,
			data:   map[string]interface{}{"Name": "Neo", "Age": 21},
		},

		{
			tmpl:   `Hello, {{.Name}}`,
			expect: `Hello, Neo`,
			data:   Person{Name: "Neo"},
		},
		{
			tmpl:   `Hello, {{.Name}}, Age = {{.Age}}`,
			expect: `Hello, Neo, Age = 21`,
			data:   Person{Name: "Neo", Age: 21},
		},
	}

	for i, tt := range tests {
		got, err := Render(tt.tmpl, tt.data)
		Assertf(t, err == nil, "err: %v", err)
		Assertf(t, got == tt.expect, "%d: expect = %q, got = %q", i, tt.expect, got)
	}
}
