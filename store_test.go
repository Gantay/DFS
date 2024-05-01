package main

import (
	"bytes"

	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathKey := CASPathTransformFunc(key)

	expectedOriginalKey := "be17b32c2870b1c0c73b59949db6a3be7814dd23"
	expectedPathName := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	if pathKey.PathName != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Original != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.Original, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
