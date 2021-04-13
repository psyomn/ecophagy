// This is pretty inefficient but wanted to bash out something pretty
// quick and for fun. Don't use this in production.
package ds

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("could not find item")
	ErrEmptyList = errors.New("linked list is empty")
)

type Node struct {
	next *Node
	Item interface{}
}

type LinkedList struct {
	head *Node
	size uint64
}

func LinkedListNew() *LinkedList {
	return &LinkedList{}
}

func (s *LinkedList) Peek() (interface{}, error) {
	if s.head == nil {
		return nil, ErrEmptyList
	}

	return s.head.Item, nil
}

func (s *LinkedList) Add(item interface{}) {
	cursor := s.nextEmpty()
	if cursor == nil {
		s.head = &Node{}
		s.head.Item = item
	} else {
		cursor.next = &Node{}
		cursor = cursor.next
		cursor.Item = item
	}
	s.size++
}

func (s *LinkedList) Delete(item interface{}) error {
	node, lag, err := s.Find(item)
	if err != nil {
		return err
	}

	// pray to the garbage collector, ohmmmmm
	lag.next = node.next

	s.size--
	return nil
}

func (s *LinkedList) Find(item interface{}) (*Node, *Node, error) {
	lag := s.head
	for cursor := s.head; cursor != nil; cursor = cursor.next {
		if cursor.Item == item {
			return cursor, lag, nil
		}
		lag = cursor
	}

	return nil, nil, ErrNotFound
}

func (s *LinkedList) IsEmpty() bool {
	return s.size == 0
}

func (s *LinkedList) Length() uint64 {
	return s.size
}

func (s *LinkedList) nextEmpty() *Node {
	if s.head == nil {
		return nil
	}

	curr := s.head
	for {
		if curr != nil && curr.next == nil {
			return curr
		}

		curr = curr.next
	}
}
