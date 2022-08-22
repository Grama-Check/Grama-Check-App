package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Grama-Check/Grama-Check-App/util"
	"github.com/stretchr/testify/require"
)

func createRandomCheck(t *testing.T) Check {
	args := CreateCheckParams{
		Nic:     util.RandomID(),
		Name:    util.RandomName(),
		Address: util.RandomAddress(),
		Email:   util.RandomEmail(),
	}

	check, err := testQueries.CreateCheck(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, check)

	require.Equal(t, args.Nic, check.Nic)
	require.Equal(t, args.Name, check.Name)
	require.Equal(t, args.Address, check.Address)
	require.Equal(t, args.Address, check.Address)
	require.Equal(t, false, check.Idcheck)
	require.Equal(t, false, check.Addresscheck)
	require.Equal(t, false, check.Policecheck)
	require.Equal(t, false, check.Failed)

	return check
}

func TestCreateCheck(t *testing.T) {
	createRandomCheck(t)
}

func TestGetCheck(t *testing.T) {
	check := createRandomCheck(t)

	check2, err := testQueries.GetCheck(context.Background(), check.Nic)

	require.NoError(t, err)
	require.NotEmpty(t, check2)

	require.Equal(t, check.Nic, check2.Nic)
	require.Equal(t, check.Name, check2.Name)
	require.Equal(t, check.Address, check2.Address)
	require.Equal(t, check.Email, check2.Email)
	require.Equal(t, check.Idcheck, check2.Idcheck)
	require.Equal(t, check.Addresscheck, check2.Addresscheck)
	require.Equal(t, check.Policecheck, check2.Policecheck)
	require.Equal(t, check.Failed, check2.Failed)
}

func TestUpdateIdentityCheck(t *testing.T) {
	check := createRandomCheck(t)

	testQueries.UpdateIdentityCheck(context.Background(), check.Nic)

	check, err := testQueries.GetCheck(context.Background(), check.Nic)

	require.NoError(t, err)

	require.Equal(t, true, check.Idcheck)
}

func TestUpdateAddressCheck(t *testing.T) {
	check := createRandomCheck(t)

	testQueries.UpdateAddressCheck(context.Background(), check.Nic)

	check, err := testQueries.GetCheck(context.Background(), check.Nic)

	require.NoError(t, err)

	require.Equal(t, true, check.Addresscheck)
}

func TestUpdatePoliceCheck(t *testing.T) {
	check := createRandomCheck(t)

	testQueries.UpdatePoliceCheck(context.Background(), check.Nic)

	check, err := testQueries.GetCheck(context.Background(), check.Nic)

	require.NoError(t, err)

	require.Equal(t, true, check.Policecheck)
}

func TestUpdateFailed(t *testing.T) {
	check := createRandomCheck(t)

	testQueries.UpdateFailed(context.Background(), check.Nic)

	check, err := testQueries.GetCheck(context.Background(), check.Nic)

	require.NoError(t, err)

	require.Equal(t, true, check.Failed)
}

func TestDeleteCheck(t *testing.T) {
	check1 := createRandomCheck(t)
	err := testQueries.DeleteCheck(context.Background(), check1.Nic)
	require.NoError(t, err)

	check2, err := testQueries.GetCheck(context.Background(), check1.Nic)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, check2)
}
