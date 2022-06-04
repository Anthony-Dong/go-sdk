package main

import (
	"fmt"
	"regexp"
)

func main() {
	//mustCompile := regexp.MustCompile(`[\S]+`)
	//fmt.Println( mustCompile.MatchString("\t\t"))
	//fmt.Println( mustCompile.MatchString(""))
	//fmt.Println( mustCompile.MatchString("11"))
	//
	//compile := regexp.MustCompile(`@title(\s*:\s*)(.*)`)
	//
	//submatch := compile.FindStringSubmatch("@title: Get .....", -1)
	//fmt.Println(len(submatch))
	//fmt.Println()
	//for _,elem := range submatch {
	//	fmt.Printf("%#v\n",elem)
	//}

	/**
	Code:
	re := regexp.MustCompile(`a(x*)b(y|z)c`)
	fmt.Printf("%q\n", re.FindStringSubmatch("-axxxbyc-"))
	fmt.Printf("%q\n", re.FindStringSubmatch("-abzc-"))
	Output:
	["axxxbyc" "xxx" "y"]
	["abzc" "" "z"]
	*/

	/**
	Code:
	re := regexp.MustCompile(`a(x*)b`)
	fmt.Printf("%q\n", re.FindAllStringSubmatch("-ab-", -1))
	fmt.Printf("%q\n", re.FindAllStringSubmatch("-axxb-", -1))
	fmt.Printf("%q\n", re.FindAllStringSubmatch("-ab-axb-", -1))
	fmt.Printf("%q\n", re.FindAllStringSubmatch("-axxb-ab-", -1))
	Output:
	[["ab" ""]]
	[["axxb" "xx"]]
	[["ab" ""] ["axb" "x"]]
	[["axxb" "xx"] ["ab" ""]]
	*/

	/**
	`(?P<name>[a-zA-Z]+)\s+(?P<age>\d+)\s+(?P<email>\w+@\w+(?:\.\w+)+)`
	(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})
	*/
	re := regexp.MustCompile(`(?P<name>[a-zA-Z]+)\s+(?P<age>\d+)\s+(?P<email>\w+@\w+(?:\.\w+)+)`)
	names := re.SubexpNames()

	fmt.Printf("%#v\n", names)
	for _, elem := range re.FindAllStringSubmatch("tom 18 127@gmail.com", -1) {
		fmt.Printf("%#v\n", elem)
	}
	for _, elem := range re.FindAllStringSubmatch("-axxb-ab-", -1) {
		fmt.Printf("all: %#v\n", elem)
	}

	fmt.Printf("%#v\n", re.FindStringSubmatch("-axxb-ab-"))
}
