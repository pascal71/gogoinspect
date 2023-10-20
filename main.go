package main

import "gogoinspect/gogoinspect"

func main() {
	slice := []int{1, 2, 3}
	gogoinspect.VisualizeSlice(slice)

	slice2 := slice
	gogoinspect.Visualize2Slices(slice, slice2)

	slice2[1] = 42
	gogoinspect.Visualize2Slices(slice, slice2)

	slice2 = append(slice2, 77)
	gogoinspect.Visualize2Slices(slice, slice2)

	// Demonstrating for string slices
	strSlice := []string{"Hello", "World"}
	strSlice2 := []string{"Hi", "Earth"}
	gogoinspect.Visualize2Slices(strSlice, strSlice2)
	strSlice2 = append(strSlice, strSlice2...)
	gogoinspect.Visualize2Slices(strSlice, strSlice2)

	fltSlice1 := []float64{3.14, 1.41, 2.81}
	fltSlice2 := fltSlice1

	gogoinspect.Visualize2Slices(fltSlice1, fltSlice2)

	str := "Hello, Go!"
	gogoinspect.VisualizeString(str)
	str += "Yeah!"
	gogoinspect.VisualizeString(str)

	buffer1 := make([]byte, 16)
	buffer2 := buffer1

	gogoinspect.Visualize2Slices(buffer1, buffer2)
	buffer1 = []byte(str)
	gogoinspect.Visualize2Slices(buffer1, buffer2)

	buffer2 = append(buffer1, buffer1...)
	gogoinspect.Visualize2Slices(buffer1, buffer2)
}
