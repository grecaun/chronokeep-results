package sqlite

import "testing"

var (
	blockedEmails []string
	blockedPhones []string
)

func setupBlockedTests() {
	if len(blockedEmails) < 1 {
		blockedEmails = []string{
			"banned@test.com",
		}
	}
	if len(blockedPhones) < 1 {
		blockedPhones = []string{
			"3605551234",
		}
	}
}

func TestAddBlockedPhone(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.AddBlockedPhone(blockedPhones[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestAddBlockedPhones(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.AddBlockedPhones(blockedPhones)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phones: %v", err)
	}
}

func TestGetBlockedPhones(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	_, err = db.GetBlockedPhones()
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestUnblockPhone(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.UnblockPhone(blockedPhones[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestAddBlockedEmail(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.AddBlockedEmail(blockedEmails[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestAddBlockedEmails(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.AddBlockedEmails(blockedEmails)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestGetBlockedEmails(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	_, err = db.GetBlockedEmails()
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}

func TestUnblockEmail(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	err = db.UnblockEmail(blockedEmails[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
}
