package ds

import (
	"errors"
	"testing"
)

func errorIf(t *testing.T, err error) {
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestCreateLinkedList(t *testing.T) {
	if LinkedListNew() == nil {
		t.Errorf("should be able to create empty list")
	}
}

func TestLinkedListIsEmpty(t *testing.T) {
	ll := LinkedListNew()
	if !ll.IsEmpty() {
		t.Errorf("should be empty open creation")
	}
}

func TestLinkedListAdd(t *testing.T) {
	ll := LinkedListNew()

	ll.Add(interface{}(1))
	ll.Add(interface{}("hello"))
	ll.Add(interface{}(1.123))

	if ll.Length() != 3 {
		t.Errorf("should have added 3 items")
	}
}

func TestLinkedListFind(t *testing.T) {
	ll := LinkedListNew()
	ll.Add(interface{}(100))
	ll.Add(interface{}(200))
	ll.Add(interface{}(300))
	ll.Add(interface{}(400))

	node, lag, err := ll.Find(300)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if node == nil {
		t.Fatalf("%s", "should have node data")
	}

	value := node.Item.(int)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if value != 300 {
		t.Errorf("did not find the right value")
	}

	if lag.next != node {
		t.Errorf("bad lag")
	}
}

func TestLinkedListFindSmall(t *testing.T) {
	ll := LinkedListNew()
	ll.Add(interface{}(300))
	ll.Add(interface{}(400))

	node, lag, err := ll.Find(300)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if node == nil {
		t.Errorf("should have node data")
		return
	}

	value := node.Item.(int)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if value != 300 {
		t.Errorf("did not find the right value")
	}

	if lag == nil {
		t.Errorf("bad lag")
	}
}

func TestLinkedListDelete(t *testing.T) {
	ll := LinkedListNew()

	ll.Add(interface{}(1))
	ll.Add(interface{}(2))
	ll.Add(interface{}(3))

	if ll.Length() != 3 {
		t.Errorf("bad length")
	}

	errorIf(t, ll.Delete(2))
	if ll.Length() != 2 {
		t.Errorf("bad length")
	}

	errorIf(t, ll.Delete(1))
	if ll.Length() != 1 {
		t.Errorf("bad length")
	}

	errorIf(t, ll.Delete(3))
	if ll.Length() != 0 {
		t.Errorf("bad length")
	}
}

func TestLinkedListDeleteEmpty(t *testing.T) {
	ll := LinkedListNew()
	err := ll.Delete(10)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("should be of error not found")
	}
}

func TestLinkedListPeekEmpty(t *testing.T) {
	_, err := LinkedListNew().Peek()
	if !errors.Is(err, ErrEmptyList) {
		t.Errorf("cannot peek empty list")
	}
}
