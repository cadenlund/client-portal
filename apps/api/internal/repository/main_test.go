// Author: Caden Lund
// Created: 4/11/2026
// Last updated: 4/11/2026
// Notes:
// - test main runs once on setup

package repository

import (
	"log"
	"testing"

	"github.com/cadenlund/client-portal/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool that tests use, global - setup once per package
var testPool *pgxpool.Pool

// Main setup function, runs once before all test
func TestMain(m *testing.M) {
	//1. Create testcontainer and return pool
	pool, err := testutil.Setup()
	if err != nil {
		log.Fatalf("Failed to setup container: %v", err)
	}

	//2. Define test pool as global var
	testPool = pool

	//3. Run all tests
	m.Run()

	//4. Teardown test container
	err = testutil.Cleanup()
	if err != nil {
		log.Fatalf("Failed to cleanup container: %v", err)
	}

}
