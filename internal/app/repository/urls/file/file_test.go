package file

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	filePath      = "storage.dat"
	defaultUserID = "user"
)

func TestFileRepo_Add(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo("storage.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwe", "yandex.ru", defaultUserID)
	require.NoError(t, err)
}

func TestFileRepo_Get(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "abc", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	act, err := repo.Get(ctx, "abc")

	assert.NoError(t, err)
	assert.Equal(t, "yandex.ru", act)
}

func TestFileRepo_FetchURls_Success(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "avito.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepo(filePath)
	require.NoError(t, err)

	act, err := repo.FetchURLs(ctx, defaultUserID)
	require.NoError(t, err)

	assert.Len(t, act, 2)
}

func TestFileRepo_FetchURls_NotFound(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepo(filePath)
	require.NoError(t, err)

	act, err := repo.FetchURLs(ctx, "fake")
	require.NoError(t, err)

	assert.Len(t, act, 0)
}

func TestFileRepo_Ping(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Ping(ctx)
	assert.NoError(t, err)
}
