package main

import (
	"IM/sql"
	"IM/utils"
	"context"
	"fmt"
)

var ctx = context.Background()

func main() {
	utils.Redis_init()
	redget := sql.RedZRevRange("msg:_5", 0, -1)
	fmt.Println(len(redget))
}
