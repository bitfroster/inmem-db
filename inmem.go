package inmemdb

import (
	"errors"
)

type Transaction[T any] struct {
	data       map[string]*T
	parent     *Transaction[T]
	committed  bool
	rolledBack bool
}

func newTransaction[T any](parent *Transaction[T]) *Transaction[T] {
	data := make(map[string]*T)
	if parent != nil {
		for k, v := range parent.data {
			data[k] = v
		}
	}
	return &Transaction[T]{data: data, parent: parent}
}

type InMemoryDatabase[T any] struct {
	data         map[string]T
	transactions []*Transaction[T]
}

func NewInMemoryDatabase[T any]() *InMemoryDatabase[T] {
	return &InMemoryDatabase[T]{
		data:         make(map[string]T),
		transactions: make([]*Transaction[T], 0),
	}
}

func (db *InMemoryDatabase[T]) Get(key string) (T, error) {
	var zero T
	if len(db.transactions) == 0 {
		val, ok := db.data[key]
		if !ok {
			return zero, errors.New("key not found")
		}
		return val, nil
	}

	tx := db.transactions[len(db.transactions)-1]
	for tx != nil {
		if val, ok := tx.data[key]; ok {
			if val == nil {
				return zero, errors.New("key not found")
			}
			return *val, nil
		}
		tx = tx.parent
	}

	val, ok := db.data[key]
	if !ok {
		return zero, errors.New("key not found")
	}
	return val, nil
}

func (db *InMemoryDatabase[T]) Set(key string, value T) error {
	if len(db.transactions) == 0 {
		return errors.New("no transaction started")
	}
	tx := db.transactions[len(db.transactions)-1]
	tx.data[key] = &value
	return nil
}

func (db *InMemoryDatabase[T]) Delete(key string) error {
	if len(db.transactions) == 0 {
		return errors.New("no transaction started")
	}
	tx := db.transactions[len(db.transactions)-1]
	tx.data[key] = nil
	return nil
}

func (db *InMemoryDatabase[T]) StartTransaction() {
	tx := newTransaction[T](nil)
	if len(db.transactions) > 0 {
		parent := db.transactions[len(db.transactions)-1]
		tx = newTransaction(parent)
	}
	db.transactions = append(db.transactions, tx)
}

func (db *InMemoryDatabase[T]) Commit() error {
	if len(db.transactions) == 0 {
		return errors.New("no transaction started")
	}
	tx := db.transactions[len(db.transactions)-1]
	if tx.committed {
		return errors.New("transaction already committed")
	}
	if tx.rolledBack {
		return errors.New("transaction already rolled back")
	}
	tx.committed = true
	if tx.parent == nil {
		db.data = make(map[string]T)
		for k, v := range tx.data {
			if v != nil {
				db.data[k] = *v
			}
		}
	} else {
		for k, v := range tx.data {
			tx.parent.data[k] = v
		}
	}
	db.transactions = db.transactions[:len(db.transactions)-1]
	return nil
}

func (db *InMemoryDatabase[T]) RollBack() error {
	if len(db.transactions) == 0 {
		return errors.New("no transaction started")
	}
	tx := db.transactions[len(db.transactions)-1]
	if tx.committed {
		return errors.New("transaction already committed")
	}
	if tx.rolledBack {
		return errors.New("transaction already rolled back")
	}
	tx.rolledBack = true
	db.transactions = db.transactions[:len(db.transactions)-1]
	return nil
}
