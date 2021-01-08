package stack

import "container/list"

type Stack struct {
	list *list.List
}

func New() *Stack{
	return &Stack{list: list.New()}
}

func (s *Stack) Push(v interface{}) interface{}{
	s.list.PushBack(v)
	return v
}

func (s *Stack) Pop() interface{}{
	if s.Len() == 0{
		return nil
	}
	e := s.list.Back()
	s.list.Remove(e)
	return e.Value
}

func (s *Stack) Peek() interface{}{
	if s.Len() == 0{
		return nil
	}
	return s.list.Back().Value
}

func (s *Stack) Len() int{
	return s.list.Len()
}

func (s *Stack) Init() *Stack{
	s.list.Init()
	return s
}

