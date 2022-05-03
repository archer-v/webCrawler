package main

import "testing"

func TestQueue(t *testing.T) {
	queue := UrlsQueue{}
	taskId1 := 1
	arr1 := []string{"one", "two"}

	t.Run("Put to the empty queue", func(t *testing.T) {
		queue.Put(&taskId1, arr1)
		if queue.Len() != 2 {
			t.Fatalf("Wrong queue length")
		}
	})

	taskId2 := 2
	arr2 := []string{"three", "four"}

	t.Run("Put to the not empty queue", func(t *testing.T) {
		queue.Put(&taskId2, arr2)
		if queue.Len() != 4 {
			t.Fatalf("Wrong queue length")
		}
	})

	var arr []string
	arr = append(arr, arr1...)
	arr = append(arr, arr2...)

	t.Run("Get elements from the queue", func(t *testing.T) {
		for i, element := range arr {
			taskId, str, e := queue.Get()
			if e != nil {
				t.Fatalf("Unexpected error is received, %v", e)
			}
			if str != element {
				t.Fatalf("Wrong queue data is received, expected %v, but got %v", element, str)
			}
			if i < len(arr1) && *taskId.(*int) != taskId1 {
				t.Fatalf("Wrong taskId received, expected %v, but got %v", taskId1, *taskId.(*int))
			}
			if i >= len(arr1) && *taskId.(*int) != taskId2 {
				t.Fatalf("Wrong taskId received, expected %v, but got %v", taskId2, *taskId.(*int))
			}
		}
	})

	t.Run("Get elements from the empty queue", func(t *testing.T) {
		_, _, e := queue.Get()
		if e == nil {
			t.Fatalf("Error was expected on empty queue")
		}
	})
}
