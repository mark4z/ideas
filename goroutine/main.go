package main

/*
#include "stdio.h"

static void SayHello(const char* name) {
	puts(name);
}
*/
import "C"

func main() {
	C.SayHello(C.CString("Hello, World!"))
}
