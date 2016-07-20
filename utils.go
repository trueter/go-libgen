package main

import (
  "log"
)

// Util log functions
func scream(err error) {
    if err != nil {
      log.Fatal(err)
      panic(err)
    }
}
