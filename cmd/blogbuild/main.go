package main

import "log"

func main() {
	if err := buildBlog(); err != nil {
		log.Fatal(err)
	}

	log.Println("blog build complete")
}