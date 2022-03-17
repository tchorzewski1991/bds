package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	service = "service"
	ts      = "ts"
	level   = "level"
	traceID = "trace_id"
	msg     = "msg"
)

func main() {
	var b strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		m := make(map[string]interface{})

		err := json.Unmarshal([]byte(line), &m)
		if err != nil {
			continue
		}

		b.Reset()

		if v, ok := m[service]; ok {
			b.WriteString(fmt.Sprintf("%s | ", v))
		}

		if v, ok := m[ts]; ok {
			b.WriteString(fmt.Sprintf("%s | ", v))
		}

		if v, ok := m[level]; ok {
			b.WriteString(fmt.Sprintf("%s | ", v))
		}

		if v, ok := m[traceID]; ok {
			b.WriteString(fmt.Sprintf("%s | ", v))
		}

		if v, ok := m[msg]; ok {
			b.WriteString(fmt.Sprintf("%s | ", v))
		}

		for k, v := range m {
			switch k {
			case service, ts, level, traceID, msg:
				continue
			}
			b.WriteString(fmt.Sprintf("%s[%v] | ", k, v))
		}

		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
