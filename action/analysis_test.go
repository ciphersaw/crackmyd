package action

import (
	"testing"
)

func TestAnalyseMariadb_10_3(t *testing.T) {
	PwdMode = "assign"
	PwdFile = "../example/password.txt"
	SufMode = "assign"
	SufFile = "../example/suffix.txt"

	obj := "../example/user-mariadb-10.3.MYD"
	result, err := analyseFile(obj)
	if err != nil {
		t.Errorf("analyseFile test user-mariadb-10.3.MYD fails: %v", err)
	}

	for _, r := range result {
		switch r.user {
		case "kali":
			if r.plaintext != "qwerty" {
				t.Errorf("expected: qwerty, got: %v", r.plaintext)
			}
		case "kalinew":
			if r.plaintext != "q1w2e3r4" {
				t.Errorf("expected: q1w2e3r4, got: %v", r.plaintext)
			}
		case "app":
			if r.plaintext != "app123" {
				t.Errorf("expected: app123, got: %v", r.plaintext)
			}
		case "appnew":
			if r.plaintext != "appnew@gmail.com" {
				t.Errorf("expected: appnew@gmail.com, got: %v", r.plaintext)
			}
		case "crackmyd":
			if r.plaintext != "crackmyd" {
				t.Errorf("expected: crackmyd, got: %v", r.plaintext)
			}
		case "stronger":
			if r.plaintext != "" {
				t.Errorf("expected: <empty string>, got: %v", r.plaintext)
			}
		default:
			t.Errorf("unknown user: %v", r.user)
		}
	}
}

func TestAnalyseMysql_5_7(t *testing.T) {
	PwdMode = "assign"
	PwdFile = "../example/password.txt"
	SufMode = "assign"
	SufFile = "../example/suffix.txt"

	obj := "../example/user-mysql-5.7.MYD"
	result, err := analyseFile(obj)
	if err != nil {
		t.Errorf("analyseFile test user-mysql-5.7.MYD fails: %v", err)
	}

	for _, r := range result {
		switch r.user {
		case "centos":
			if r.plaintext != "qwerty" {
				t.Errorf("expected: qwerty, got: %v", r.plaintext)
			}
		case "centosnew":
			if r.plaintext != "q1w2e3r4" {
				t.Errorf("expected: q1w2e3r4, got: %v", r.plaintext)
			}
		case "app":
			if r.plaintext != "app123" {
				t.Errorf("expected: app123, got: %v", r.plaintext)
			}
		case "appnew":
			if r.plaintext != "appnew@gmail.com" {
				t.Errorf("expected: appnew@gmail.com, got: %v", r.plaintext)
			}
		case "crackmyd":
			if r.plaintext != "crackmyd" {
				t.Errorf("expected: crackmyd, got: %v", r.plaintext)
			}
		default:
			t.Errorf("unknown user: %v", r.user)
		}
	}
}
