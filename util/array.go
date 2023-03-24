package util

func ArrayReverse(arr []string) {
	i, j := 0, len(arr)-1
	for {
		if i < len(arr)/2 {
			arr[i], arr[j] = arr[j], arr[i]
			i++
			j--
		} else {
			break
		}
	}
}

func ArrayContains[T comparable](arr []T, item T) bool {
	for i := range arr {
		if arr[i] == item {
			return true
		}
	}
	return false
}

// ArrayInsertSort 插入排序
// 从2-array.length开始，比较i与0..i=j之前的所有数，如果找到比i小的数，则将j的数据开始向后推，然后i放入j的位置
func ArrayInsertSort[T interface{}](array []T, comparable func(a, b T) bool) {
	//如果数组长度等于1，则不需要排序
	if len(array) == 1 {
		return
	}
	//遍历从1..arr.length
	for i := 1; i < len(array); i++ {
		item := array[i]
		j := i - 1
		//如果比较通过，则将j的值后推一位
		for j >= 0 && comparable(item, array[j]) {
			array[j+1] = array[j]
			j--
		}
		//当比较完毕，找到要插入的位置，插入
		array[j+1] = item
	}
}

// ArrayBubbleSort 冒泡排序
func ArrayBubbleSort[T interface{}](array []T, comparable func(a T, b T) bool) {
	for i := 0; i < len(array)-1; i++ {
		for j := i + 1; j < len(array); j++ {
			if comparable(array[i], array[j]) {
				array[i], array[j] = array[j], array[i]
			}
		}
	}
}

func ArrayFilter[T interface{}](array []T, filter func(item T) bool) []T {
	var result []T
	for i := range array {
		if filter(array[i]) {
			result = append(result, array[i])
		}
	}
	return result
}
