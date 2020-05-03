package main

import (
	"fmt"
	"regexp"
)

func main() {
	r := regexp.MustCompile(`^(START|END) (\S*) \"(.*)\"$`)
	o := r.FindStringSubmatch("START snes \"to jest _ próba\"")

	fmt.Printf("%#v\n", o)
	o = r.FindStringSubmatch("START snes \"to jest _ blekota\"")

	fmt.Printf("%#v\n", o)
}
