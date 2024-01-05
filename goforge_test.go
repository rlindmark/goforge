package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidModuleReleaseFilename(t *testing.T) {

	var testCases = []struct {
		moduleReleaseFilename string
		expect                bool
	}{
		{"puppetlabs-stdlib-9.4.1.tar.gz", true},
		{"puppetlabs-stdlib-9.4.1.tar.g", false},
		{"puppetlabs-stdlib-9.4.1", false},
		{"puppetlabs-stdlib-", false},
		{"puppetlabs-stdlib", false},
		{"puppetlabs-", false},
		{"puppetlabs", false},
		{"p", false},
		{"", false},
	}

	for _, test := range testCases {

		result, _ := validModuleReleaseFilename(test.moduleReleaseFilename)

		if result != test.expect {
			t.Errorf("TestValidModuleReleaseFilename %v: got %v, expected %v", test.moduleReleaseFilename, result, test.expect)
		}
	}

}

func TestV3File(t *testing.T) {
	t.Run("returns a forge module", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v3/files/puppetlabs-stdlib-9.4.0.tar.gz", nil)
		response := httptest.NewRecorder()

		handleV3Files(response, request)

		//got := response.Body.String()
		got := response.Result().StatusCode
		if got != 200 {
			t.Errorf("expected 200 got %d", got)
		}
		got = int(response.Result().ContentLength)
		if got != 161699 {
			t.Errorf("expected contentlengt = 161699 got %d", got)

		}
	})
}
