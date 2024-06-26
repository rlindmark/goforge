package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadModuleRelease(t *testing.T) {
	var testCases = []struct {
		path       string
		statusCode int
	}{
		{"/v2/files/puppetlabs", 500},
		{"/v3/files/puppetlabs", 400},
		{"/v3/files/not-found-0.0.1.tar.gz", 404},
	}

	t.Run("returns a forge user", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			DownloadModuleRelease(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("expected statuscode = %d got %d", test.statusCode, got)
			}
		}
	})
}

func TestCreateModuleRelease(t *testing.T) {

	t.Run("create a forge module", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodPost, "/v3/release", nil)
		response := httptest.NewRecorder()

		CreateModuleRelease(response, request)

		got := response.Result().StatusCode
		if got != 401 {
			t.Errorf("Expected %v got %v", 401, got)
		}

		request, _ = http.NewRequest(http.MethodPost, "/v3/release", nil)
		request.Header.Add("Authorization", "Bearer <api_key>")

		response = httptest.NewRecorder()

		CreateModuleRelease(response, request)
		got = response.Result().StatusCode
		if got != 403 {
			t.Errorf("Expected %v got %v", 403, got)
		}

	})
}

func TestFetchModule(t *testing.T) {
	var testCases = []struct {
		path          string
		statusCode    int
		contentLength int
		response      string
	}{
		{"/v3/m/puppetlabs-stdlib-9.4.1", 500, -1, "ok"},
		{"/v3/mmodulesa/puppetlabs-stdlib-9.4.1", 500, -1, "ok"},

		{"/v3/modules/puppetlabs-stdlib-9.4.1", 400, -1, "ok"},

		{"/v3/modules/puppetlabs-stdlib", 404, -1, "ok"},
	}

	t.Run("returns a forge module", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			FetchModule(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("url = %v expected statuscode = %d got %d", test.path, test.statusCode, got)
			}
			//fmt.Printf("response.Result() = %v", response.Result())
			got = int(response.Result().ContentLength)
			if got != test.contentLength {
				t.Errorf("path:%v expected content lengt = %d got %d", test.path, test.contentLength, got)

			}
		}
	})
}

func TestFetchModuleRelease(t *testing.T) {
	var testCases = []struct {
		path          string
		statusCode    int
		contentLength int
		response      string
	}{
		{"/v3/releases/puppetlabs-stdlib-9.4.1", 200, -1, "ok"},
		{"/v3/releases/puppetlabs-stdlib-9.4.0", 200, -1, "ok"},
		{"/v3/releases/puppetlabs-stdlib-0.0.0", 404, -1, "ok"},
	}

	t.Run("returns a forge module", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			FetchModuleRelease(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("url = %v expected statuscode = %d got %d", test.path, test.statusCode, got)
			}
			//fmt.Printf("response.Result() = %v", response.Result())
			got = int(response.Result().ContentLength)
			if got != test.contentLength {
				t.Errorf("path:%v expected content lengt = %d got %d", test.path, test.contentLength, got)

			}
		}
	})
}

func TestListUsers(t *testing.T) {
	var testCases = []struct {
		path          string
		statusCode    int
		contentLength int
		response      string
	}{
		{"/v3/users", 200, 161699, "ok"},
	}

	t.Run("returns a forge user", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			ListUsers(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("expected statuscode = %d got %d", test.statusCode, got)
			}
		}
	})
}

func TestFetchUser(t *testing.T) {
	var testCases = []struct {
		path       string
		statusCode int
	}{
		{"/v3/users/puppetlabs", 200},
		{"/v3/userss/puppetlabs", 500},
		{"/v4/users/puppetlabs", 500},
	}

	t.Run("returns a forge user", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			FetchUser(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("expected statuscode = %d got %d", test.statusCode, got)
			}
		}
	})
}

func TestFetchUserParams(t *testing.T) {
	var testCases = []struct {
		path       string
		statusCode int
	}{
		{"/v3/users/puppetlabs?with_html=true", 200},
		{"/v3/users/puppetlabs?with_html=false", 200},
		{"/v3/users/puppetlabs?with_html=debug", 200},
	}

	t.Run("returns a forge user", func(t *testing.T) {

		for _, test := range testCases {

			request, _ := http.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			FetchUser(response, request)

			//got := response.Body.String()
			got := response.Result().StatusCode
			if got != test.statusCode {
				t.Errorf("expected statuscode = %d got %d", test.statusCode, got)
			}
		}
	})
}
