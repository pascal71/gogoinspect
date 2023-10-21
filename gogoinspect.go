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
	fmt.Println(
		"    +--------------------+-------+-------+              +--------------------+-------+-------+",
	)
	fmt.Println(
		"    | Slice1             |  Len  |  Cap  |              | Slice2             |  Len  |  Cap  |",
	)
	fmt.Println(
		"    +--------------------+-------+-------+              +--------------------+-------+-------+",
	)
	fmt.Printf("    | 0x%016x | %3d   | %3d   |              | 0x%016x | %3d   | %3d   |\n",
		header1.Data, header1.Len, header1.Cap,
		header2.Data, header2.Len, header2.Cap)
	fmt.Println(
		"    +--------------------+-------+-------+              +--------------------+-------+-------+",
	)
	fmt.Println("        |                                                    |")
	fmt.Println("        v                                                    v")
	fmt.Println("        |                                                    |")
	fmt.Println(
		"    +--------------------+-------------------------+    +--------------------+-------------------------+",
	)

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
			fmt.Printf(
				"    | 0x%016x | %-22s  |    +--------------------+-------------------------+\n",
				address1,
				display1,
			)
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
		fmt.Println(
			"    +--------------------+-------------------------+    +--------------------+-------------------------+",
		)
	case header1.Len > header2.Len:
		fmt.Println("    +---------------------+------------------------+")
	case header2.Len > header1.Len:
		fmt.Println(
			"                                                        +--------------------+-------------------------+",
		)
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
	fmt.Println("    +----------------+")
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
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))

	fmt.Println("    +----------------------+-------+")
	fmt.Println("    | String               |  Len  |")
	fmt.Println("    +----------------------+-------+")
	fmt.Printf("    | 0x%016x   | %4d  |\n", header.Data, header.Len)
	fmt.Println("    +----------------------+-------+")
	fmt.Println("        |")
	fmt.Println("        v")
	fmt.Println("        |")
	fmt.Println("    +--------------------+--------------+")
	i := 0
	for i < len(s) {
		address := header.Data + uintptr(i)
		charDisplay, bytesRead := byteSequenceRepresentation(s[i:])
		if len(charDisplay) == 3 {
			fmt.Printf("    | 0x%016x | %-12s |\n", address, charDisplay)
		} else {
			fmt.Printf("    | 0x%016x | %-11s |\n", address, charDisplay)
		}
		for j := 1; j < bytesRead; j++ {
			fmt.Printf("    | 0x%016x | %-12s |\n", address+uintptr(j), "")
		}
		fmt.Println("    +--------------------+--------------+")
		i += bytesRead
	}
	fmt.Println()
}

func byteSequenceRepresentation(subs string) (charDisplay string, bytesRead int) {
	b := subs[0]
	switch {
	case b < 0x80:
		return fmt.Sprintf("'%c'", b), 1
	case b >= 0xC0 && b < 0xE0:
		return fmt.Sprintf("\"%s\"", subs[:2]), 2
	case b >= 0xE0 && b < 0xF0:
		return fmt.Sprintf("\"%s\"", subs[:3]), 3
	case b >= 0xF0:
		return fmt.Sprintf("\"%s\"", subs[:4]), 4
	}
	return "", 0
}

func VisualizeStruct(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct!")
		return
	}

	structName := val.Type().Name()
	fmt.Printf("[%s struct]\n", structName)

	fmt.Printf(
		"+--------------------------+--------------------+-----------+------------------+----------+------------+\n",
	)
	fmt.Printf(
		"| Field Name               | Memory Address     | Value     | Type             | Size     | Padding    |\n",
	)
	fmt.Printf(
		"+--------------------------+--------------------+-----------+------------------+----------+------------+\n",
	)

	var prevFieldEnd uintptr
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		size := field.Type().Size()
		addr := field.UnsafeAddr()
		padding := uintptr(0)

		if i != 0 {
			padding = addr - prevFieldEnd
		}
		prevFieldEnd = addr + size

		valueStr := fmt.Sprint(field.Interface())
		if len(valueStr) > 8 {
			valueStr = valueStr[:5] + "..."
		}

		fmt.Printf("| %-24s | 0x%016x | %-9s | %-16s | %8d | %10d |\n",
			fieldType.Name, addr, valueStr, field.Type(), size, padding)
	}

	fmt.Printf(
		"+--------------------------+--------------------+-----------+------------------+----------+------------+\n",
	)
}
