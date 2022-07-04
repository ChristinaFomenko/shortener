package file

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const filePath = "storage.dat"

func TestFileRepo_Add(t *testing.T) {
	repo, err := NewRepo("storage.dat")
	if err != nil {
		return
	}

	err = repo.Add("hello", "world")
	if err != nil {
		return
	}

	err = os.Remove(filePath)
	if err != nil {
		return
	}
}

func TestFileRepo_Get(t *testing.T) {
	repo, err := NewRepo(filePath)
	if err != nil {
		return
	}

	err = repo.Add("hi", "Go")
	if err != nil {
		return
	}

	err = repo.Add("hi", "Chris")
	if err != nil {
		return
	}

	err = repo.Add("good", "morning")
	if err != nil {
		return
	}

	err = repo.Add("good", "morning")
	if err != nil {
		return
	}

	repo, err = NewRepo("storage.dat")
	if err != nil {
		return
	}

	act, err := repo.Get("good")
	if err != nil {
		return
	}

	if act != "morning" {
		t.Error(act)
	}

	err = os.Remove(filePath)
	if err != nil {
		return
	}
}

func TestFileRepo_GetList(t *testing.T) {
	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("abcde", "yandex.ru")
	require.NoError(t, err)

	err = repo.Add("qwerty", "github.com")
	require.NoError(t, err)

	repo, err = NewRepo("storage.dat")
	require.NoError(t, err)

	act, err := repo.GetList()
	require.NoError(t, err)

	assert.Len(t, act, 2)
}
