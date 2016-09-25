// Copyright 2016 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import (
	"bytes"
	"testing"
)

func TestTextTemplate(t *testing.T) {
	data := struct{ Name string }{"Francesc"}
	tests := []struct {
		text    string
		out     string
		invalid bool
	}{
		{text: "{{.Name}}", out: "Francesc"},
		{text: "", out: ""},
		{text: "{{", invalid: true},
		{text: "Hello, {{.Name}}", out: "Hello, Francesc"},
	}

	for _, test := range tests {
		var tt textTemplate
		if err := tt.Set(test.text); err != nil {
			if !test.invalid {
				t.Errorf("parsing %s failed unexpectedly: %v", test.text, err)
			}
			continue
		}
		if test.invalid {
			t.Errorf("parsing %s should have failed", test.text)
			continue
		}

		if tt.text != test.text {
			t.Errorf("the text template should be %s; got %s", test.text, tt.text)
		}

		var buf bytes.Buffer
		if err := tt.t.Execute(&buf, data); err != nil {
			t.Errorf("unexpected error executing %s: %v", test.text, err)
			continue
		}
		if got := buf.String(); got != test.out {
			t.Errorf("expected text was %s; got %s", test.out, got)
			continue
		}
	}
}
