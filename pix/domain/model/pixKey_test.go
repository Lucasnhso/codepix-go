package model_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/Lucasnhso/codepix-go/domain/model"
	"github.com/stretchr/testify/require"
)

func TestModel_NewPixKey(t *testing.T) {
	code := "001"
	name := "test bank"
	bank, err := model.NewBank(code, name)

	accountNumber := "abcnumber"
	ownerName := "Lucas"
	account, err := model.NewAccount(bank, accountNumber, ownerName)

	kind := "email"
	key := "j@j.com"
	pixKey, err := model.NewPixKey(kind, account, key)

	require.NotEmpty(t, uuid.FromStringOrNil(pixKey.ID))
	require.Equal(t, pixKey.Kind, kind)
	require.Equal(t, pixKey.Status, "active")

	kind = "cpf"
	_, err = model.NewPixKey(kind, account, key)
	require.Nil(t, err)

	_, err = model.NewPixKey("nome", account, key)
	require.NotNil(t, err)
}