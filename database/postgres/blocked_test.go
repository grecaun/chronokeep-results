package postgres

import "testing"

var (
	blockedEmails []string
	blockedPhones []string
)

func setupBlockedTests() {
	if len(blockedEmails) < 1 {
		blockedEmails = []string{
			"banned@test.com",
			"bannedtwo@test.com",
		}
	}
	if len(blockedPhones) < 1 {
		blockedPhones = []string{
			"3865551234",
			"8215559876",
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
	// Verify we don't have any added
	nums, _ := db.GetBlockedPhones()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	// Test different adds
	err = db.AddBlockedPhone(blockedPhones[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 1, len(nums))
	}
	if nums[0] != blockedPhones[0] {
		t.Errorf("Expected to find %v but found %v.", blockedPhones[0], nums[0])
	}
	// To make sure we add multiples
	err = db.AddBlockedPhone(blockedPhones[1])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
	foundNums := 0
	for _, outer := range blockedPhones {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != 2 {
		t.Errorf("Expected to find %v number matches, found %v.", 2, foundNums)
	}
	// Verify we don't add a phone multiple times
	err = db.AddBlockedPhone(blockedPhones[1])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phone: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
}

func TestAddBlockedPhones(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Verify we don't have any added
	nums, _ := db.GetBlockedPhones()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	// Test adding slice
	err = db.AddBlockedPhones(blockedPhones[1:])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phones: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != len(blockedPhones)-1 {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedPhones)-1, len(nums))
	}
	// Test adding full array
	err = db.AddBlockedPhones(blockedPhones)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phones: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != len(blockedPhones) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedPhones), len(nums))
	}
	foundNums := 0
	for _, outer := range blockedPhones {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != len(blockedPhones) {
		t.Errorf("Expected to find %v number matches, found %v.", len(blockedPhones), foundNums)
	}
	// Ensure we don't get doubles
	err = db.AddBlockedPhones(blockedPhones)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked phones: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != len(blockedPhones) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedPhones), len(nums))
	}
}

func TestGetBlockedPhones(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Verify we don't have any added
	nums, err := db.GetBlockedPhones()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked phones: %v", err)
	}
	// Test getting single
	_ = db.AddBlockedPhone(blockedPhones[0])
	nums, err = db.GetBlockedPhones()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v numbers, found %v.", 1, len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked phones: %v", err)
	}
	foundNums := 0
	for _, outer := range blockedPhones {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != 1 {
		t.Errorf("Expected to find %v number matches, found %v.", 1, foundNums)
	}
	// Test getting multiple
	_ = db.AddBlockedPhones(blockedPhones)
	nums, err = db.GetBlockedPhones()
	if len(nums) != len(blockedPhones) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedPhones), len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked phones: %v", err)
	}
	// verify they're all there
	foundNums = 0
	for _, outer := range blockedPhones {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != len(blockedPhones) {
		t.Errorf("Expected to find %v number matches, found %v.", len(blockedPhones), foundNums)
	}
}

func TestUnblockPhone(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Add and verify
	_ = db.AddBlockedPhones(blockedPhones)
	nums, _ := db.GetBlockedPhones()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
	// Verify we can unblock a number
	err = db.UnblockPhone(blockedPhones[0])
	if err != nil {
		t.Errorf("Unexpected error when unblocking blocked phone: %v", err)
	}
	nums, _ = db.GetBlockedPhones()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 1, len(nums))
	}
	found := false
	for _, phone := range nums {
		if blockedPhones[0] == phone {
			found = true
		}
	}
	if found {
		t.Errorf("Found number that should have been unblocked.")
	}
}

func TestAddBlockedEmail(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Verify we don't have any added
	nums, _ := db.GetBlockedEmails()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	// Test different adds
	err = db.AddBlockedEmail(blockedEmails[0])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked email: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 1, len(nums))
	}
	if nums[0] != blockedEmails[0] {
		t.Errorf("Expected to find %v but found %v.", blockedEmails[0], nums[0])
	}
	// To make sure we add multiples
	err = db.AddBlockedEmail(blockedEmails[1])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked email: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
	foundNums := 0
	for _, outer := range blockedEmails {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != 2 {
		t.Errorf("Expected to find %v number matches, found %v.", 2, foundNums)
	}
	// Verify we don't add an email multiple times
	err = db.AddBlockedEmail(blockedEmails[1])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked email: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
}

func TestAddBlockedEmails(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Verify we don't have any added
	nums, _ := db.GetBlockedEmails()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	// Test adding slice
	err = db.AddBlockedEmails(blockedEmails[1:])
	if err != nil {
		t.Errorf("Unexpected error when adding blocked emails: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != len(blockedEmails)-1 {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedEmails)-1, len(nums))
	}
	// Test adding full array
	err = db.AddBlockedEmails(blockedEmails)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked emails: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != len(blockedEmails) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedEmails), len(nums))
	}
	foundNums := 0
	for _, outer := range blockedEmails {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != len(blockedEmails) {
		t.Errorf("Expected to find %v number matches, found %v.", len(blockedEmails), foundNums)
	}
	// Ensure we don't get doubles
	err = db.AddBlockedEmails(blockedEmails)
	if err != nil {
		t.Errorf("Unexpected error when adding blocked emails: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != len(blockedEmails) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedEmails), len(nums))
	}
}

func TestGetBlockedEmails(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Verify we don't have any added
	nums, err := db.GetBlockedEmails()
	if len(nums) > 0 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 0, len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked emails: %v", err)
	}
	// Test getting single
	_ = db.AddBlockedEmail(blockedEmails[0])
	nums, err = db.GetBlockedEmails()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v numbers, found %v.", 1, len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked emails: %v", err)
	}
	foundNums := 0
	for _, outer := range blockedEmails {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != 1 {
		t.Errorf("Expected to find %v number matches, found %v.", 1, foundNums)
	}
	// Test getting multiple
	_ = db.AddBlockedEmails(blockedEmails)
	nums, err = db.GetBlockedEmails()
	if len(nums) != len(blockedEmails) {
		t.Errorf("Expected to find %v numbers, found %v.", len(blockedEmails), len(nums))
	}
	if err != nil {
		t.Errorf("Unexpected error getting blocked emails: %v", err)
	}
	// verify they're all there
	foundNums = 0
	for _, outer := range blockedEmails {
		for _, inner := range nums {
			if outer == inner {
				foundNums += 1
				break
			}
		}
	}
	if foundNums != len(blockedEmails) {
		t.Errorf("Expected to find %v number matches, found %v.", len(blockedEmails), foundNums)
	}
}

func TestUnblockEmail(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupBlockedTests()
	// Add and verify
	_ = db.AddBlockedEmails(blockedEmails)
	nums, _ := db.GetBlockedEmails()
	if len(nums) != 2 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 2, len(nums))
	}
	// Verify we can unblock a number
	err = db.UnblockEmail(blockedEmails[0])
	if err != nil {
		t.Errorf("Unexpected error when unblocking blocked email: %v", err)
	}
	nums, _ = db.GetBlockedEmails()
	if len(nums) != 1 {
		t.Errorf("Expected to find %v blocked numbers, found %v.", 1, len(nums))
	}
	found := false
	for _, email := range nums {
		if blockedEmails[0] == email {
			found = true
		}
	}
	if found {
		t.Errorf("Found number that should have been unblocked.")
	}
}
