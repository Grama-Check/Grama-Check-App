package db

import (
	"context"
	"testing"

	"github.com/Grama-Check/Grama-Check-App/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Nic:     util.RandomID(),
		Name:    util.RandomName(),
		Address: util.RandomAddress(),
		Email:   util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Nic, user.Nic)
	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Address, user.Address)
	require.Equal(t, args.Address, user.Address)
	require.Equal(t, false, user.Idcheck)
	require.Equal(t, false, user.Addresscheck)
	require.Equal(t, false, user.Policecheck)
	require.Equal(t, false, user.Failed)

	return user

}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetPerson(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user.Nic)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Nic, user2.Nic)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Address, user2.Address)
	require.Equal(t, user.Email, user2.Email)

}

func TestSetIdentityCheck(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.UpdateID(context.Background(), user.Nic)

	user, err := testQueries.GetUser(context.Background(), user.Nic)

	require.NoError(t, err)

	require.Equal(t, true, user.Idcheck)
}
func TestSetPoliceCheck(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.UpdatePolice(context.Background(), user.Nic)

	user, err := testQueries.GetUser(context.Background(), user.Nic)

	require.NoError(t, err)

	require.Equal(t, true, user.Policecheck)
}
func TestSetAddressCheck(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.UpdateAddress(context.Background(), user.Nic)

	user, err := testQueries.GetUser(context.Background(), user.Nic)

	require.NoError(t, err)

	require.Equal(t, true, user.Addresscheck)
}
func TestSetFailed(t *testing.T) {
	user := CreateRandomUser(t)

	testQueries.UpdateFailed(context.Background(), user.Nic)

	user, err := testQueries.GetUser(context.Background(), user.Nic)

	require.NoError(t, err)

	require.Equal(t, true, user.Failed)
}
