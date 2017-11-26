# GOPTER

... the GOlang Property TestER
[![Build Status](https://travis-ci.org/leanovate/gopter.svg?branch=master)](https://travis-ci.org/leanovate/gopter)
[![codecov.io](https://codecov.io/github/leanovate/gopter/coverage.svg?branch=master)](https://codecov.io/github/leanovate/gopter?branch=master)
[![GoDoc](https://godoc.org/github.com/leanovate/gopter?status.png)](https://godoc.org/github.com/leanovate/gopter)

## Synopsis

Gopter tries to bring the goodness of [ScalaCheck](https://www.scalacheck.org/) (and impliticly the goodness of [QuickCheck](http://hackage.haskell.org/package/QuickCheck)) to Go.
It can be also seen as a more sophisticated version of the testing/quick package.

Main differences to ScalaCheck:

* It is Go ... duh
* ... nevertheless: Do not expect the same typesafety and elegance as in ScalaCheck.
* For simplicity [Shrink](https://www.scalacheck.org/files/scalacheck_2.11-1.12.5-api/index.html#org.scalacheck.Shrink) has become part of the generators. They can be still easily changed if necessary.
* There is no [Pretty](https://www.scalacheck.org/files/scalacheck_2.11-1.12.5-api/index.html#org.scalacheck.util.Pretty) ... so far gopter feels quiet comfortable being ugly.
* A generator for regex matches
* No parallel commands ... yet?

Main differences to the testing/quick package:

* Much tighter control over generators
* Shrinkers, i.e. automatically find the minimum value falsifying a property
* A generator for regex matches (already mentioned that ... but it's cool)

## Documentation

Current godocs:

* [gopter](https://godoc.org/github.com/leanovate/gopter):  Main interfaces
* [gopter/gen](https://godoc.org/github.com/leanovate/gopter/gen): All commonly used generators
* [gopter/prop](https://godoc.org/github.com/leanovate/gopter/prop): Common helpers to create properties from a condition function and an specific generators
* [gopter/arbitrary](https://godoc.org/github.com/leanovate/gopter/arbitrary): Helpers automatically combine generators for arbitrary types
* [gopter/commands](https://godoc.org/github.com/leanovate/gopter/commands): Helpers to create stateful tests based on arbirtrary commands
* [gopter/convey](https://godoc.org/github.com/leanovate/gopter/convey): Helpers to used gopter inside goconvey tests

## License

[MIT Licence](http://opensource.org/licenses/MIT)
