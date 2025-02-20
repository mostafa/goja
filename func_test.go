package goja

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFuncProto(t *testing.T) {
	const SCRIPT = `
	"use strict";
	function A() {}
	A.__proto__ = Object;
	A.prototype = {};

	function B() {}
	B.__proto__ = Object.create(null);
	var thrown = false;
	try {
		delete B.prototype;
	} catch (e) {
		thrown = e instanceof TypeError;
	}
	thrown;
	`
	testScript(SCRIPT, valueTrue, t)
}

func TestFuncPrototypeRedefine(t *testing.T) {
	const SCRIPT = `
	let thrown = false;
	try {
		Object.defineProperty(function() {}, "prototype", {
			set: function(_value) {},
		});
	} catch (e) {
		if (e instanceof TypeError) {
			thrown = true;
		} else {
			throw e;
		}
	}
	thrown;
	`

	testScript(SCRIPT, valueTrue, t)
}

func TestFuncExport(t *testing.T) {
	vm := New()
	typ := reflect.TypeOf((func(FunctionCall) Value)(nil))

	f := func(expr string, t *testing.T) {
		v, err := vm.RunString(expr)
		if err != nil {
			t.Fatal(err)
		}
		if actualTyp := v.ExportType(); actualTyp != typ {
			t.Fatalf("Invalid export type: %v", actualTyp)
		}
		ev := v.Export()
		if actualTyp := reflect.TypeOf(ev); actualTyp != typ {
			t.Fatalf("Invalid export value: %v", ev)
		}
	}

	t.Run("regular function", func(t *testing.T) {
		f("(function() {})", t)
	})

	t.Run("arrow function", func(t *testing.T) {
		f("(()=>{})", t)
	})

	t.Run("method", func(t *testing.T) {
		f("({m() {}}).m", t)
	})

	t.Run("class", func(t *testing.T) {
		f("(class {})", t)
	})
}

func ExampleAssertConstructor() {
	vm := New()
	res, err := vm.RunString(`
		(class C {
			constructor(x) {
				this.x = x;
			}
		})
	`)
	if err != nil {
		panic(err)
	}
	if ctor, ok := AssertConstructor(res); ok {
		obj, err := ctor(nil, vm.ToValue("Test"))
		if err != nil {
			panic(err)
		}
		fmt.Print(obj.Get("x"))
	} else {
		panic("Not a constructor")
	}
	// Output: Test
}
