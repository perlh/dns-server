package util

import (
	"fmt"
	"testing"
	"time"
)

func TestReverse(t *testing.T) {
	arr := []string{"1", "2", "3"}
	ArrayReverse(arr)
	fmt.Println(arr)
}

func TestArrayBubbleSort(t *testing.T) {
	micro := time.Now().UnixMicro()
	array := []int{
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12,
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12}
	//a-b是正排序，b-a是到排序
	ArrayBubbleSort(array, func(a int, b int) bool {
		return a-b > 0
	})
	fmt.Println(array)
	fmt.Println(time.Now().UnixMicro() - micro)
}

func TestArrayInsertSort(t *testing.T) {
	micro := time.Now().UnixMicro()
	array := []int{
		2, 3, 1, 2, 11, 6, 7, 4, 5, 10, 100, 52, 25, 35, 45, 88, 14, 12}
	ArrayInsertSort(array, func(a int, b int) int {
		return a - b
	})
	fmt.Println(array)
	fmt.Println(time.Now().UnixMicro() - micro)
}