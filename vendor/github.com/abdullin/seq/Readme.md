Structural equality library for Golang.

## Story

While we were working on [HappyPancake](http://abdullin.com/happypancake/) project, [Pieter](https://twitter.com/pjvds) always wanted to have a better assertion library for our event-driven specifications.

Months later, awesome folks from [@DDDBE](https://twitter.com/dddbe) community presented me with some Trappist Beer. Thanks to it (and some spare time on the weekend), this assertion library was finally written.

## How it works

You can define expectations on objects (e.g. API responses or expected events) by creating an instance of `seq.Map`, which is provided by the package `github.com/abdullin/seq`.

Maps can be nested or they could have flat paths. Values could be represented with strings, primitive types, instances of `seq.Map` or JSON-serializable objects.

Consider following types:

```go
type Robot struct {
    Legs int    `json:"legs"`
    Arms int    `json:"arms"`
    Name string `json:"name"`
}

type Party struct {
    Rating  []int             `json:"rating"`
    Seating map[string]*Robot `json:"seating"`
}
```
Let's imagine that our JSON API returns `Party` object, which we want to verify. We could define our expectation like this:

```go
  expect := seq.Map{
    // array
    "rating.len": 3,
    "rating[1]":  5,
    // flat path with value terminator
    "seating.front.name": "R2D2",
    "seating.front.arms": "1",
    "seating.front.legs": 3,
    // flat path with map terminator
    "seating.right": seq.Map{
      "name": "C3PO",
    },
    // flat path with object terminator
    "seating.back": &Robot{
      Name: "Marvin",
      Legs: 2,
      Arms: 2,
    },
  }

```
Once you have the expectation, you could compare it with an actual object. Here is an example of a valid object:
```go
  actual := &Party{
    Rating: []int{4, 5, 4},
    Seating: map[string]*Robot{
      "front": &Robot{
        Name: "R2D2",
        Arms: 1,
        Legs: 3,
      },
      "back": &Robot{
        Name: "Marvin",
        Legs: 2,
        Arms: 2,
      },
      "right": &Robot{
        Name: "C3PO",
        Legs: 2,
        Arms: 2,
      },
    },
  }
  result := expect.Test(actual)

```
Result value would contain `Issues []seq.Issue` with any differences and could be checked like this:

```go
if !result.Ok() {
  fmt.Println("Differences")
  for _, v := range result.Issues {
    fmt.Println(v.String())
  }
}
```

If actual object has some invalid or missing properties, then result will have nice error messages. Consider this object:

```go
  actual := &Party{
    Seating: map[string]*Robot{
      "front": &Robot{
        Name: "R2D2",
        Arms: 1,
      },
      "back": &Robot{
        Name: "Marvin",
        Arms: 3,
      },
      "right": &Robot{
        Name: "C4PO",
        Legs: 2,
        Arms: 3,
      },
    },
  }

```

If verified against the original expectation, `result.Issues` would contain these error messages:

```
Expected rating.len to be '3' but got nothing
Expected seating.back.legs to be '2' but got '0'
Expected seating.front.legs to be '3' but got '0'
Expected rating[1] to be '5' but got nothing
Expected seating.back.arms to be '2' but got '3'
Expected seating.right.name to be 'C3PO' but got 'C4PO'
```
Check out the [unit tests](https://github.com/abdullin/seq/blob/master/seq_test.go) for more examples.

## Feedback

Feedback is welcome and appreciated!
