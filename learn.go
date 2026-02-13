package main

import (
	"fmt"
	"math"
)

// This is a simple Go program that demonstrates basic syntax, variable declarations, and loops.

func printHelloAndLoop() {
	fmt.Println("Hello world")
	// This is a simple loop that prints the iteration number
	// := is the implicit variable declaration and assignment operator in Go
	// with it, we don't need to do something like var i int and then i = 0, we can just do i := 0
	for i := 0; i < 5; i++ {
		fmt.Println("Iteration: ", i)
	}
}

func variableDeclarations() {
	// --------- VARIABLE DECLARATIONS ---------
	var name string = "Erika"
	// this is a SIGNED integer type, meaning it can represent both positive and negative numbers
	var age int = 23
	var isStudent bool = false

	fmt.Printf("Name: %s, Age: %d, Is Student: %t\n", name, age, isStudent)

	// We can also use implicit variable declaration for these types
	city := "New York"
	temperature := 25.5
	isRaining := true

	fmt.Printf("City: %s, Temperature: %.1f°C, Is Raining: %t\n", city, temperature, isRaining)

	// theres also uint which is an unsigned integer type, meaning it can only represent non-negative numbers
	var unsignedInt uint = 100
	fmt.Printf("Unsigned Integer: %d\n", unsignedInt)

	var floatEx1 float32 = 3.14
	var floatEx2 float64 = 3.141592653589793
	fmt.Printf("Float32: %.2f, Float64: %.15f\n", floatEx1, floatEx2)

	var byteEx byte = 'A' // byte is an ALIAS for uint8, it can represent a single ASCII character
	// %c is the format specifier for characters, it will print the character representation of the byte value
	fmt.Printf("Byte: %c\n", byteEx)

	var runeEx rune = '世' // rune is an ALIAS for int32, it can represent a single Unicode character
	// %c is also the format specifier for characters, it will print the character representation of the rune value
	fmt.Printf("Rune: %c\n", runeEx)
}

func declarationTypes() {
	// -------- TYPES OF DECLARATIONS --------
	// 1. Explicit Type Declaration
	var explicitInt int = 42
	fmt.Printf("Explicit Int: %d\n", explicitInt)

	// 2. Implicit Type Declaration
	implicitFloat := 3.14
	fmt.Printf("Implicit Float: %.2f\n", implicitFloat)

	// 3. Multiple Variable Declaration
	var x, y, z int = 1, 2, 3
	fmt.Printf("Multiple Variables: x=%d, y=%d, z=%d\n", x, y, z)

	// 4. Short Variable Declaration in a Loop
	for i := 0; i < 3; i++ {
		fmt.Printf("Loop Variable: %d\n", i)
	}

	// 5. CONST
	const Pi = 3.14159
	fmt.Printf("Constant Pi: %.5f\n", Pi)

	// NOTE:
	// Usually we want to use implicit variable declaration when the type can be easily inferred from the value,
	// and explicit type declaration when we want to specify a particular type or when the type cannot be easily inferred.
}

func typeCasting() {
	// --------- TYPE CASTING ---------
	var intVal int = 42
	var floatVal float64 = float64(intVal) // Explicit type conversion from int to float64
	fmt.Printf("Integer Value: %d, Float Value: %.8f\n", intVal, floatVal)
}

func consoleOutput() {
	// --------- CONSOLE OUTPUT ---------
	// fmt.Printf allows us to format our output using format specifiers
	// %d for integers, %f for floating-point numbers, %s for strings, %t for booleans, and %c for characters
	name := "Erika"
	age := 23
	isErika := true
	fmt.Printf("Name: %s, Age: %d, Is Erika: %t\n", name, age, isErika)

	x := 69.12938094023847
	// %e is the format specifier for scientific notation, it will print the float in the form of a number between 1 and 10 multiplied by a power of 10
	fmt.Printf("Formatted Float: %e\n", x)

	y := 69.12938094023847
	// %.2f is the format specifier for floating-point numbers with 2 decimal places, it will round the float to 2 decimal places
	fmt.Printf("Formatted Float with 2 Decimal Places: %.2f\n", y)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func arithmetic() {
	// --------- ARITHMETIC OPERATIONS ---------
	a := 10
	b := 3

	sum := a + b
	difference := a - b
	product := a * b
	quotient := a / b
	remainder := a % b

	fmt.Printf("Sum: %d, Difference: %d, Product: %d, Quotient: %d, Remainder: %d\n", sum, difference, product, quotient, remainder)

	// string concat
	str1 := "Hello, "
	str2 := "world!"
	concatenated := str1 + str2
	fmt.Printf("Concatenated String: %s\n", concatenated)

	// math
	fmt.Println(math.Round(4.56546))
}

// MAIN FUNCTION --> ENTRY POINT OF THE PROGRAM
func main() {
	printHelloAndLoop()
	variableDeclarations()
	declarationTypes()
	typeCasting()
	consoleOutput()
	fmt.Println(fibonacci(10))
	arithmetic()
}