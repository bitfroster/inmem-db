package inmemdb

import (
	"testing"
)

func TestNewInMemoryDatabase(t *testing.T) {
	db := NewInMemoryDatabase[string]()
	if len(db.data) != 0 {
		t.Errorf("Expected empty data map, but got %v", db.data)
	}
	if len(db.transactions) != 0 {
		t.Errorf("Expected empty transactions slice, but got %v", db.transactions)
	}
}

func TestGet(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Get when key doesn't exist
	_, err := db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key, but got nil")
	}

	// Test Get after Set
	db.StartTransaction()
	db.Set("key1", "value1")
	db.Commit()
	val, err := db.Get("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}

	// Test Get after Delete
	db.StartTransaction()
	db.Delete("key1")
	db.Commit()
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key after Delete, but got nil")
	}

	// Test Get after nested transactions
	db.StartTransaction()
	db.Set("key1", "value1")
	db.StartTransaction()
	db.Delete("key1")
	db.Commit()
	db.Commit()
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key after nested transactions, but got nil")
	}
}

func TestSet(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Set without transaction
	err := db.Set("key1", "value1")
	if err == nil {
		t.Errorf("Expected error for Set without transaction, but got nil")
	}

	// Test Set within transaction
	db.StartTransaction()
	err = db.Set("key1", "value1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	db.Commit()
	val, err := db.Get("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}
}

func TestDelete(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Delete without transaction
	err := db.Delete("key1")
	if err == nil {
		t.Errorf("Expected error for Delete without transaction, but got nil")
	}

	// Test Delete within transaction
	db.StartTransaction()
	db.Set("key1", "value1")
	err = db.Delete("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	db.Commit()
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key after Delete, but got nil")
	}
}

func TestStartTransaction(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Start transaction
	db.StartTransaction()
	if len(db.transactions) != 1 {
		t.Errorf("Expected 1 transaction after StartTransaction, but got %d", len(db.transactions))
	}

	// Test Start nested transaction
	db.StartTransaction()
	if len(db.transactions) != 2 {
		t.Errorf("Expected 2 transactions after StartTransaction, but got %d", len(db.transactions))
	}
}

func TestCommit(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Commit without transaction
	err := db.Commit()
	if err == nil {
		t.Errorf("Expected error for Commit without transaction, but got nil")
	}

	// Test Commit after transaction
	db.StartTransaction()
	db.Set("key1", "value1")
	err = db.Commit()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	val, err := db.Get("key1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}

	// Test Commit error cases
	db.StartTransaction()
	err = db.Commit()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = db.Commit()
	if err == nil {
		t.Errorf("Expected error for Commit after transaction committed, but got nil")
	}
	err = db.RollBack()
	if err == nil {
		t.Errorf("Expected error for RollBack after transaction committed, but got nil")
	}
}

func TestRollBack(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test RollBack without transaction
	err := db.RollBack()
	if err == nil {
		t.Errorf("Expected error for RollBack without transaction, but got nil")
	}

	// Test RollBack after transaction
	db.StartTransaction()
	db.Set("key1", "value1")
	err = db.RollBack()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key after RollBack, but got nil")
	}

	// Test RollBack error cases
	db.StartTransaction()
	err = db.Commit()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = db.RollBack()
	if err == nil {
		t.Errorf("Expected error for RollBack after transaction committed, but got nil")
	}
}

func TestNestedTransactions(t *testing.T) {
	db := NewInMemoryDatabase[string]()

	// Test Nested Transactions
	db.StartTransaction()
	db.Set("key3", "value3")
	db.StartTransaction()
	db.Set("key3", "value3-modified")
	db.Delete("key1")
	db.Commit()
	val, err := db.Get("key3")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "value3-modified" {
		t.Errorf("Expected value3-modified, got %s", val)
	}
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key, but got nil")
	}
	db.Commit()
	val, err = db.Get("key3")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "value3-modified" {
		t.Errorf("Expected value3-modified, got %s", val)
	}
	_, err = db.Get("key1")
	if err == nil {
		t.Errorf("Expected error for non-existent key after Commit, but got nil")
	}
}
