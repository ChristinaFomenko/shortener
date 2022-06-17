package file

import (
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
