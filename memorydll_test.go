package memorydll

// //#include<stdlib.h> // for C.free
// import "C" // it cannot test, so it's commented
import (
	_ "embed"
	"testing"
)

/*
//file example.c
gcc -c -fpic example.c

gcc -shared example.o -o example.dll

#include <stdio.h>
/ * Compute the greatest common divisor of positive integers * /

int gcd(int x, int y) {
    int g;
    g = y;
    while (x > 0) {
        g = x;
        x = y % x;
        y = g;
    }
    return g;
}

void print(char * hello){
    printf(hello);
}
*/

// the following dll is example.dll
//
//go:embed testdata/example.dll
var exampleDll []byte

func TestMemoryLoadLibrary(t *testing.T) {
	dll, err := NewDLL(exampleDll, "example.dll")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dll.Release)

	proc, err := dll.FindProc("gcd")
	if err != nil {
		t.Fatal(err)
	}
	result, _, _ := proc.Call(4, 8)
	if n := int(result); n != 4 {
		t.Fatalf("gcd calc error, want 4, got %d", int(n))
	}

	// // go test complains import "C"
	// proc, err = dll.FindProc("print")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// cstr := C.CString("hello,world")
	// defer C.free(unsafe.Pointer(cstr))
	// proc.Call(cstr)
}
