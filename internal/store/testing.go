package store

import (
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()

	config := NewConfig()
	config.DatabaseURL = databaseURL
	st := New(config)
	if err := st.Open(); err != nil {
		t.Fatal(err)
	}

	return st, func(tables ...string) {
		if len(tables) > 0 {
			query := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))
			if _, err := st.db.Exec(query); err != nil {
				t.Fatal(err)
			}
		}

		st.Close()
	}
}
