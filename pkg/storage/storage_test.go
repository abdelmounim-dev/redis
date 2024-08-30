package storage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/abdelmounim-dev/redis/pkg/parser"
)

func TestKeyValueStore_Set(t *testing.T) {
	store := &KeyValueStore{
		data: make(map[string]*parser.Token),
	}

	tests := []struct {
		name    string
		key     string
		value   *parser.Token
		wantErr bool
	}{
		{
			name:    "Set simple string",
			key:     "test",
			value:   &parser.Token{Type: parser.SimpleString, Value: "value"},
			wantErr: false,
		},
		{
			name:    "Set integer",
			key:     "int",
			value:   &parser.Token{Type: parser.Integer, Value: int64(42)},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyValueStore.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if the value was actually set
			if !tt.wantErr {
				if val, ok := store.data[tt.key]; !ok || val.Type != tt.value.Type || val.Value != tt.value.Value {
					t.Errorf("KeyValueStore.Set() did not set the correct value for key %s", tt.key)
				}
			}
		})
	}
}

func TestKeyValueStore_Get(t *testing.T) {
	store := &KeyValueStore{
		data: map[string]*parser.Token{
			"existing": {Type: parser.SimpleString, Value: "value"},
		},
	}

	tests := []struct {
		name    string
		key     string
		want    *parser.Token
		wantErr bool
	}{
		{
			name:    "Get existing key",
			key:     "existing",
			want:    &parser.Token{Type: parser.SimpleString, Value: "value"},
			wantErr: false,
		},
		{
			name:    "Get non-existing key",
			key:     "nonexistent",
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyValueStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got == nil || got.Type != tt.want.Type || got.Value != tt.want.Value) {
				t.Errorf("KeyValueStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyValueStore_Delete(t *testing.T) {
	store := &KeyValueStore{
		data: map[string]*parser.Token{
			"existing": {Type: parser.SimpleString, Value: "value"},
		},
	}

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "Delete existing key",
			key:     "existing",
			wantErr: false,
		},
		{
			name:    "Delete non-existing key",
			key:     "nonexistent",
			wantErr: false, // Assuming Delete is idempotent and doesn't error on non-existent keys
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyValueStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if the key was actually deleted
			if _, exists := store.data[tt.key]; exists {
				t.Errorf("KeyValueStore.Delete() did not delete the key %s", tt.key)
			}
		})
	}
}

func TestKeyValueStore_Concurrency(t *testing.T) {
	store := &KeyValueStore{
		data: make(map[string]*parser.Token),
	}

	const numOperations = 1000
	var wg sync.WaitGroup
	wg.Add(numOperations * 3) // Set, Get, and Delete operations

	for i := 0; i < numOperations; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := &parser.Token{Type: parser.Integer, Value: int64(i)}
			_ = store.Set(key, value)
		}(i)

		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			_, _ = store.Get(key)
		}(i)

		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			_ = store.Delete(key)
		}(i)
	}

	wg.Wait()
	// If we've reached here without deadlocks or race conditions, the test passes
}
