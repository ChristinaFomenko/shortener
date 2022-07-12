package file

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	filePath      = "storage.dat"
	defaultUserID = "abcde"
)

func TestFileRepo_Add(t *testing.T) {
	repo, err := NewRepo("store.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("qwe", defaultUserID, "yandex.ru")
	require.NoError(t, err)
}

func TestFileRepo_Get(t *testing.T) {
	repo, err := NewRepo("storage.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("abc", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	act, err := repo.Get("abc")
	require.NoError(t, err)

	assert.Equal(t, "yandex.ru", act)
}

func TestFileRepo_GetList_Success(t *testing.T) {
	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("abcde", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	err = repo.Add("edcba", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	repo, err = NewRepo("storage.dat")
	require.NoError(t, err)

	act, err := repo.GetList(defaultUserID)
	require.NoError(t, err)

	assert.Len(t, act, 2)
}

func TestFileRepo_GetList_NotFound(t *testing.T) {
	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("qwerty", defaultUserID, "avito.ru")
	require.NoError(t, err)

	err = repo.Add("ytrewq", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	repo, err = NewRepo("storage.dat")
	require.NoError(t, err)

	act, err := repo.GetList("fake")
	require.NoError(t, err)

	assert.Len(t, act, 0)
}
