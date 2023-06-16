package main

import "IM/router"

func main() {
	r := router.Router()
	r.Run(":8081")
}
