#argjoy

[![Build Status](https://travis-ci.org/lunixbochs/argjoy.svg?branch=master)](https://travis-ci.org/lunixbochs/argjoy)
[![GoDoc](https://godoc.org/github.com/lunixbochs/argjoy?status.svg)](http://godoc.org/github.com/lunixbochs/argjoy)

Golang module allowing you to call a method using callbacks to translate arguments. Also allows optional arguments.

You could write a function `test(a, b, c int)`, call it with `[]string{"1", "2"}`, which would effectively call `test(1, 2, 0)`.

Basic example:

    package main

    import (
        "fmt"
        "github.com/lunixbochs/argjoy"
        "strconv"
    )

    func test(a, b, optC int) int {
        return a + b + optC
    }

    func main() {
        aj := argjoy.NewArgjoy()
        // this is just to demonstrate making a simple codec
        // if you actually need string to int, argjoy.StrToInt is more robust
        aj.Register(func(arg, val interface{}) (err error) {
            if v, ok := val.(string); ok {
                if a, ok := arg.(*int); ok {
                    *a, err = strconv.Atoi(v)
                    return
                }
            }
            return argjoy.NoMatchErr
        })
        aj.Optional = true

        // Argjoy.Call returns an interface slice: []interface{}
        // the following is effectively: out, err := test(1, 2, 0)
        out, err := aj.Call(test, "1", "2")
        if err != nil {
            panic(err)
        }
        fmt.Println(out[0].(int))
    }

##Why?

Reduces duplicate decoding logic. Repeated decoding like `v, err := strconv.Atoi(); if err != nil` is completely eliminated.

Type-safe command-line parsing. Flag, Cobra, etc only reduce to a list of strings, so the rest of the decoding was up to you.

Syscall argument decoding via metaprogramming (used in [Usercorn](https://github.com/lunixbochs/usercorn)).
