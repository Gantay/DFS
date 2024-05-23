package main

import (
	"bytes"
	"fmt"

	"io/ioutil"

	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathKey := CASPathTransformFunc(key)

	expectedFileName := "be17b32c2870b1c0c73b59949db6a3be7814dd23"
	expectedPathName := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	if pathKey.PathName != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Filename != expectedFileName {
		t.Errorf("Have %s want %s", pathKey.Filename, expectedFileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	id := generatID()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("expected to have key: %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := ioutil.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("want this: %s got this: %s", data, b)
		}

		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}
		if ok := s.Has(id, key); ok {
			t.Errorf("expected to NOT have key: %s", key)
		}
	}

}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
