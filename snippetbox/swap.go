package main

import "fmt"

// My point of confusion was the dual meaning of *
// func (*TYPE) means a parameter is a POINTER (mem address) to some value of TYPE
// *x means derefencing the pointer, or GO TO the mem address and get the value.

func swap(x, y *int) {
	// EX 1: Successful swap because saved the value referenced by x in a new variable
	//temp := *x
	//*x = *y
	//*y = temp

	// EX 2: Unsuccessful swap because temp refers to same memory address as x
	temp := x  // temp is a pointer, not a value
	*x = *y    // value that x AND temp references is changed
	*y = *temp // go to the address at y and set it to the value that temp points to.

}
func main() {
	//square(&x)
	x := 1
	y := 2
	swap(&x, &y)

	fmt.Printf("x after swap is: %d\n", x)
	fmt.Printf("y after swap is: %d\n", y)
}