package test

import (
	testifyRequire "github.com/stretchr/testify/require"
	"money_share/test_tool/database"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	require := testifyRequire.New(t)
	_, err := database.Connect()
	require.Nil(err)
}
