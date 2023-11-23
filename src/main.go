package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"wrapped/datapack"
)

func main() {
	global := js.Global()

	files := global.Get("files")

	filesMap := make(map[string]string)
	err := json.Unmarshal([]byte(files.String()), &filesMap)
	if err != nil {
		panic(err)
	}

	global.Set("files", js.Null())

	dp := datapack.Parse(filesMap)

	filesMap = nil

	fmt.Println(dp.Account.User.Username)
}
