# argjoy

[![Build Status](https://travis-ci.org/lunixbochs/argjoy.svg?branch=master)](https://travis-ci.org/lunixbochs/argjoy)
[![GoDoc](https://godoc.org/github.com/lunixbochs/argjoy?status.svg)](http://godoc.org/github.com/lunixbochs/argjoy)

Golang module allowing you to call a method using callbacks to translate arguments. Also allows optional arguments.

Extremely basic example:

    package main

    import (
        "fmt"
        "github.com/lunixbochs/argjoy"
    )

    // optC is just a variable name. The opt prefix is not required.
    func test(a, b, optC int) int {
        return a + b + optC
    }

    func main() {
        aj := argjoy.NewArgjoy(argjoy.StrToInt)
        // Enables optional arguments, where unpassed arguments are zeroed.
        aj.Optional = true

        // The following is effectively: out := test(1, 2, 0)
        out, err := aj.Call(test, "1", "2")
        if err != nil {
            panic(err)
        }
        // out is []interface{} so you need to do a type assert
        fmt.Println(out[0].(int))
    }

Custom argument decoder example (use [`argjoy.StrToInt`](https://github.com/lunixbochs/argjoy/blob/master/codecs.go#L9) for a more robust version of this):

    func strToInt(arg interface{}, vals []interface{}) (err error) {
        if v, ok := val[0].(string); ok {
            if a, ok := arg.(*int); ok {
                *a, err = strconv.Atoi(v)
                return
            }
        }
        return argjoy.NoMatch
    }

## Why?

Reduces duplicate decoding logic. Repeated decoding like `v, err := strconv.Atoi(); if err != nil` is completely eliminated.

Type-safe command-line parsing. Flag, Cobra, etc only reduce to a list of strings, so the rest of the decoding was up to you.

Syscall argument decoding via metaprogramming (used in [Usercorn](https://github.com/lunixbochs/usercorn)).
