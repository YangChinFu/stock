package main

import (
	"fmt"
	"log"

	"github.com/YangChinFu/stock/pkg/realtime"
)

func main() {
	realTimes, err := realtime.Get([]string{"2887", "2834"}...)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(realTimes)

}
