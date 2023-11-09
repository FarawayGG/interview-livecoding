//go:build integration
// +build integration

package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/farawaygg/wisdom/internal/psql"
	"github.com/farawaygg/wisdom/internal/storage"
)

var db *sql.DB

func init() {
	time.Local = time.UTC
	var err error
	db, err = psql.Open(os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
}

func getStorage(t *testing.T) *Storage {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping test in short mode")
		return nil
	}

	s, err := New(sqlx.NewDb(db, "pgx"))
	require.NoError(t, err)

	s.db.MustExec(`DELETE FROM wisdoms`)
	s.db.MustExec(`DELETE FROM authors`)

	return s
}

func TestStorage_AllQueriesAreValid(t *testing.T) {
	s := getStorage(t)
	assert.NotNil(t, s)
}

func TestStorage_GetWisdoms(t *testing.T) {
	s := getStorage(t)
	assert.NotNil(t, s)

	var (
		ctx = context.Background()

		wisdoms = []storage.Wisdom{
			wisdomgen(),
			wisdomgen(),
			wisdomgen(),
			wisdomgen(),
		}
	)

	for _, w := range wisdoms {
		err := s.CreateWisdom(ctx, w)
		require.NoError(t, err)
	}
	var got []storage.Wisdom
	err := s.GetWisdoms(ctx, func(w storage.Wisdom) error {
		got = append(got, w)
		return nil
	})
	require.NoError(t, err)
	assert.ElementsMatch(t, wisdoms, got)
}

func wisdomgen() storage.Wisdom {
	var w storage.Wisdom
	fuzz.New().Funcs(func(t *time.Time, cont fuzz.Continue) {
		*t = time.Now().Round(time.Second).UTC().Add(-time.Minute)
	}).NilChance(0).Fuzz(&w)

	return w
}
