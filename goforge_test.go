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
		{"puppetlabs-stdlib-9.4.1.targz", false},
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
	var testCases = []struct {
		path          string
		statusCode    int
		contentLength int
		response      string
	}{
		{"/v3/files/puppetlabs-stdlib-9.4.1.tar.gz", 200, 161699, "ok"},
		{"/v3/files/puppetlabs-stdlib-9.4.0.tar.gz", 200, 162679, "ok"},
		{"/v3/files/puppetlabs-stdlib-1.0.0.tar.gz", 404, -1, "ok"},
		{"/v3/files/puppetlabs-stdlib-1.0.0.tar", 400, -1, "ok"},
	}

	t.Run("returns a forge module", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			handleV3Files(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("expected statuscode = %d got %d", test.statusCode, got)
			}
			got = int(response.Result().ContentLength)
			if got != test.contentLength {
				t.Errorf("expected contentlengt = %d got %d", test.contentLength, got)

			}
		}
	})
}
