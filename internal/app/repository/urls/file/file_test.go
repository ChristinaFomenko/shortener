package file

import (
	"context"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

const (
	filePath      = "storage.dat"
	defaultUserID = "abcde"
)

func TestFileRepo_Add(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo("storage.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	action, err := repo.Add(ctx, "qwe", defaultUserID, "yandex.ru")
	require.NoError(t, err, action)
}

func TestFileRepo_Get(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo("storage.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	act, err := repo.Add(ctx, "abc", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	act, err = repo.Get(ctx, "abc")
	require.NoError(t, err)

	assert.Equal(t, "yandex.ru", act)
}

func TestFileRepo_FetchURls_Success(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	_, err = repo.Add(ctx, "abcde", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	_, err = repo.Add(ctx, "edcba", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	repo, err = NewRepo("storage.dat")
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

	_, err = repo.Add(ctx, "qwerty", defaultUserID, "avito.ru")
	require.NoError(t, err)

	_, err = repo.Add(ctx, "ytrewq", defaultUserID, "yandex.ru")
	require.NoError(t, err)

	repo, err = NewRepo("storage.dat")
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

func Test_fileRepository_AddBatch(t *testing.T) {
	type fields struct {
		store    map[string]map[string]string
		ma       sync.RWMutex
		filePath string
	}
	type args struct {
		in0    context.Context
		urls   []models.UserURL
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &fileRepository{
				store:    tt.fields.store,
				ma:       tt.fields.ma,
				filePath: tt.fields.filePath,
			}
			tt.wantErr(t, r.AddBatch(tt.args.in0, tt.args.urls, tt.args.userID), fmt.Sprintf("AddBatch(%v, %v, %v)", tt.args.in0, tt.args.urls, tt.args.userID))
		})
	}
}
