package main

import "fmt"
import "reflect"
import "strconv"
import "os"
import "demo/validation/calc"

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
)

func flush(data string) {
	f, err := os.Create("./test/check.cue")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		panic(err2.Error())
	}
}

func FuncAnalyse(m interface{}) string {
	//Reflection type of the underlying data of the interface
	x := reflect.TypeOf(m)

	numIn := x.NumIn()   //Count inbound parameters
	numOut := x.NumOut() //Count outbounding parameters

	fmt.Println("Method:", x.String())
	fmt.Println("Variadic:", x.IsVariadic()) // Used (<type> ...) ?
	fmt.Println("Package:", x.PkgPath())

	var data string
	for i := 0; i < numIn; i++ {
		inV := x.In(i)
		in_Kind := inV.Kind() //func
		fmt.Printf("\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", in_Kind, inV.Name())

		println("\n")
		data = data + fmt.Sprintf("p_%s: %s \n", strconv.Itoa(i), in_Kind)
	}

	flush(data)

	for o := 0; o < numOut; o++ {
		returnV := x.Out(0)
		return_Kind := returnV.Kind()
		fmt.Printf("\nParameter OUT: "+strconv.Itoa(o)+"\nKind: %v\nName: %v\n", return_Kind, returnV.Name())
	}

	return data
}

func printErr(prefix string, err error) {
	if err != nil {
		msg := errors.Details(err, nil)
		fmt.Printf("%s:\n%s\n", prefix, msg)
	}
}

func loose(v cue.Value) error {
	return v.Validate(
		// not final or concrete
		cue.Concrete(false),
		// check minimally
		cue.Definitions(false),
		cue.Hidden(false),
		cue.Optional(false),
	)
}

func inputGenerator(args []string) string {
	var r string
	for i, v := range args {
		r += fmt.Sprintf("p_%d: %s \n", i, v)
	}
	return r
}

func main() {
	// analysis
	a := calc.Add
	schema := FuncAnalyse(a)
	println("schema: ")
	println(schema)

	// valid
	input := inputGenerator([]string{"1", "127", "true", "\"hello\""})

	// out of bound
	// input := inputGenerator([]string{"1", "128", "true", "\"hello\""})

	// invalid
	// input := inputGenerator([]string{"1", "127", "truth", "\"hello\""})

	val := schema + input
	println("target: ")
	println(val)

	c := cuecontext.New()
	v := c.CompileString(val)

	//   try out different validation schemes
	printErr("loose error", loose(v))

	fmt.Printf("\nvalue:\n%#v\n", v)
}
