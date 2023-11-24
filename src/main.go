package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"syscall/js"
	"time"
	"wrapped/datapack"

	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	global := js.Global()

	files := global.Get("files")
	t1 := time.Now()

	length := files.Length()
	goSlice := make([]byte, length)
	js.CopyBytesToGo(goSlice, files)

	var filesMap map[string]string
	err := msgpack.Unmarshal(goSlice, &filesMap)
	if err != nil {
		panic(err)
	}

	global.Set("files", js.Null())
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))

	dp := datapack.Parse(filesMap)

	filesMap = nil

	fmt.Println(dp.Account.User.Username)

	document := global.Get("document")
	messagesList := document.Call("getElementById", "messages")

	wordPopularity := make(map[string]int)
	emotePopularity := make(map[string]int)

	for _, channel := range dp.Messages {
		for _, message := range channel.Messages {
			if message.Content == " " {
				continue
			}

			li := document.Call("createElement", "li")
			li.Set("innerHTML", message.Content)
			messagesList.Call("appendChild", li)

			words := strings.Split(message.Content, " ")

			for _, word := range words {
				word = strings.Trim(word, ".,!?")
				word = strings.ToLower(word)
				wordPopularity[word]++
			}

			re := regexp.MustCompile(`<:(\w+):(\d+)>`)
			matches := re.FindAllStringSubmatch(message.Content, -1)

			for _, match := range matches {
				emotePopularity[match[2]]++

			}
		}

		// sort words by popularity
		type kv struct {
			Key   string
			Value int
		}

		var ss []kv
		for k, v := range wordPopularity {
			ss = append(ss, kv{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		// print top 10 words
		for i := 0; i < 10; i++ {
			fmt.Println(ss[i].Key, ss[i].Value)
		}

		emotesList := document.Call("getElementById", "emotes")

		// sort emotes by popularity
		var es []kv
		for k, v := range emotePopularity {
			es = append(es, kv{k, v})
		}

		sort.Slice(es, func(i, j int) bool {
			return es[i].Value > es[j].Value
		})

		// print top 10 emotes
		for i := 0; i < 10; i++ {
			li := document.Call("createElement", "li")
			li.Set("innerHTML", fmt.Sprintf("<img src=\"https://cdn.discordapp.com/emojis/%s.png?v=1\" /> %s", es[i].Key, es[i].Value))
			emotesList.Call("appendChild", li)
		}
	}
}
