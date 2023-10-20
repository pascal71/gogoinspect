package gogoinspect

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Visualize2Slices[T any](slice1, slice2 []T) {
	header1 := (*reflect.SliceHeader)(unsafe.Pointer(&slice1))
	header2 := (*reflect.SliceHeader)(unsafe.Pointer(&slice2))

	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)

	// ASCII representation header
	fmt.Println("    +--------------------+-------+-------+              +--------------------+-------+-------+")
	fmt.Println("    | Slice1             |  Len  |  Cap  |              | Slice2             |  Len  |  Cap  |")
	fmt.Println("    +--------------------+-------+-------+              +--------------------+-------+-------+")
	fmt.Printf("    | 0x%016x | %3d   | %3d   |              | 0x%016x | %3d   | %3d   |\n",
		header1.Data, header1.Len, header1.Cap,
		header2.Data, header2.Len, header2.Cap)
	fmt.Println("    +--------------------+-------+-------+              +--------------------+-------+-------+")
	fmt.Println("        |                                                    |")
	fmt.Println("        v                                                    v")
	fmt.Println("        |                                                    |")
	fmt.Println("    +--------------------+-------------------------+    +--------------------+-------------------------+")

	// Displaying the contents of both slices side-by-side
	maxLen := max(header1.Len, header2.Len)
	for i := 0; i < maxLen; i++ {
		printElem := func(val reflect.Value, header *reflect.SliceHeader, idx int) (address uintptr, display string) {
			if idx < val.Len() {
				elem := val.Index(idx)
				address = header.Data + uintptr(idx)*uintptr(elem.Type().Size())
				display = fmt.Sprintf("(%-6s) %v", elem.Type(), elem.Interface())
			}
			return
		}

		address1, display1 := printElem(val1, header1, i)
		address2, display2 := printElem(val2, header2, i)

		if i == header2.Len {
			fmt.Printf("    | 0x%016x | %-22s  |    +--------------------+-------------------------+\n", address1, display1)
		} else if i == header1.Len {
			fmt.Printf("    +--------------------+-------------------------+    | 0x%016x | %-22s  |\n", address2, display2)
		} else if i > header2.Len {
			fmt.Printf("    | 0x%016x | %-22s  |\n", address1, display1)
		} else if i > header1.Len {
			fmt.Printf("                                                        | 0x%016x | %-22s  |\n", address2, display2)
		} else {
			fmt.Printf("    | 0x%016x | %-22s  |    | 0x%016x | %-22s  |\n", address1, display1, address2, display2)
		}
	}

	switch {
	case header1.Len == header2.Len:
		fmt.Println("    +--------------------+-------------------------+    +--------------------+-------------------------+")
	case header1.Len > header2.Len:
		fmt.Println("    +---------------------+------------------------+")
	case header2.Len > header1.Len:
		fmt.Println("                                                        +--------------------+-------------------------+")
	}

	fmt.Println()
}

func VisualizeSlice[T any](slice []T) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))

	elemType := reflect.TypeOf(slice).Elem().Name()

	// Header of the ASCII representation
	fmt.Println("    +----------------------+-------+")
	fmt.Println("    | Slice        |  Len  |  Cap  |")
	fmt.Println("    +----------------------+-------+")
	fmt.Printf("    | 0x%8x | %3d   | %3d   |\n", header.Data, header.Len, header.Cap)
	fmt.Println("    +----------------------+-------+")
	fmt.Println("        |")
	fmt.Println("        v")
	fmt.Println("        |")

	// Underlying array
	elemSize := uintptr(unsafe.Sizeof(slice[0]))
	for i := 0; i < header.Cap; i++ {
		address := header.Data + uintptr(i)*elemSize
		offset := uintptr(i) * elemSize
		if i < header.Len {
			fmt.Printf("    | (%s) %8v | 0x%8x (+%3d)\n", elemType, slice[i], address, offset)
		} else {
			fmt.Printf("    |    UNUSED      | 0x%8x (+%3d)\n", address, offset)
		}
		fmt.Println("    +----------------+")
	}

	fmt.Println()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func VisualizeString(s string) {
	// Convert the string to its header representation
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// Header of the ASCII representation
	fmt.Println("    +----------------------+-------+")
	fmt.Println("    | String               |  Len  |")
	fmt.Println("    +----------------------+-------+")
	fmt.Printf("    | 0x%016x   | %4d  |\n", header.Data, header.Len)
	fmt.Println("    +----------------------+-------+")
	fmt.Println("        |")
	fmt.Println("        v")
	fmt.Println("        |")
	fmt.Println("    +---------------------------------+")

	// Underlying byte array
	for i := 0; i < header.Len; i++ {
		address := header.Data + uintptr(i)
		char := s[i]
		// fmt.Printf("    | %4q ('%c') | 0x%016x |\n", char, char, address)
		fmt.Printf("    | 0x%016x | %4q ('%c') |\n", address, char, char)
		fmt.Println("    +---------------------------------+")
	}

	fmt.Println()
}
