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

func TestSplitModuleName(t *testing.T) {

	var testCases = []struct {
		modulename     string
		expect_owner   string
		expect_module  string
		expect_version string
	}{
		{"puppetlabs-stdlib-9.4.1.tar.gz", "puppetlabs", "stdlib", "9.4.1"},
		{"puppetlabs-stdlib-9.4.1.targz", "", "", ""},
		{"puppetlabs-stdlib-9.4.1.tar.g", "", "", ""},
		{"puppetlabs-stdlib-9.4.1", "puppetlabs", "stdlib", "9.4.1"},
		{"puppetlabs-stdlib-", "", "", ""},
		{"puppetlabs-stdlib", "", "", ""},
		{"puppetlabs-", "", "", ""},
		{"puppetlabs", "", "", ""},
		{"p", "", "", ""},
		{"", "", "", ""},
	}

	for _, test := range testCases {

		result_owner, result_module, result_version, _ := SplitModuleName(test.modulename)

		if result_owner != test.expect_owner && result_module != test.expect_module && result_version != test.expect_version {
			t.Errorf("TestSplitModuleNmae %v: got %v-%v-%v, expected %v-%v-%v", test.modulename,
				result_owner, result_module, result_version,
				test.expect_owner, test.expect_module, test.expect_version)
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

			DownloadModuleRelease(response, request)

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

func TestFetchModuleRelease(t *testing.T) {
	var testCases = []struct {
		path          string
		statusCode    int
		contentLength int
		response      string
	}{
		{"/v3/releases/puppetlabs-stdlib-9.4.1", 200, -1, "ok"},
		{"/v3/releases/puppetlabs-stdlib-9.4.0", 200, -1, "ok"},
		{"/v3/releases/puppetlabs-stdlib-1.0.0", 404, -1, "ok"},
		{"/v3/releases/puppetlabs-stdlib-1.0.0", 404, -1, "ok"},
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

func TestForge_ip(t *testing.T) {

	// test default ip
	expect := DefaultForgeIp
	got := Forge_ip()

	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}

	// test that env overrides default ip
	expect = "127.0.0.2"
	t.Setenv("FORGE_IP", expect)

	got = Forge_ip()
	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

func TestForge_port(t *testing.T) {

	// test default port
	expect := DefaultForgePort
	got := Forge_port()

	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}
	// test that env overrides default ip
	expect = "8081"
	t.Setenv("FORGE_PORT", expect)

	got = Forge_port()
	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

func TestForge_cache(t *testing.T) {

	// test default cache
	expect := DefaultCacheRoot
	got := Forge_cache()

	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}
	// test that env overrides default cache
	expect = "CACHE"
	t.Setenv("FORGE_CACHE", expect)

	got = Forge_cache()
	if expect != got {
		t.Errorf("expected %v, got %v", expect, got)
	}
}
