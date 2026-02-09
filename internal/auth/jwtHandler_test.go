package auth

import (
	"strconv"
	"testing"

	"github.com/grongoglongo/chatter-go/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateJwt_Success(t *testing.T) {
	jwtHandler := NewJwtHandler(utils.GenerateKey())

	userDto := utils.CreateUserDto()
	_, err := jwtHandler.CreateJwt(userDto)

	require.NoError(t, err)
}

func TestCreateJwt_MissingKey(t *testing.T) {
	jwtHandler := NewJwtHandler("")

	userDto := utils.CreateUserDto()
	_, err := jwtHandler.CreateJwt(userDto)

	require.Error(t, err)
}

func TestDecodeJwt_Success(t *testing.T) {
	jwtHandler := NewJwtHandler(utils.GenerateKey())

	userDto := utils.CreateUserDto()
	jwt, err := jwtHandler.CreateJwt(userDto)
	require.NoError(t, err)

	claims, err := jwtHandler.DecryptJwt(jwt)
	require.NoError(t, err)

	require.Equal(t, strconv.FormatInt(userDto.ID, 10), claims.Id)
}
