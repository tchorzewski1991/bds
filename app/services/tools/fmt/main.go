package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

		traceID := "00000000"
		if v, ok := m["trace_id"]; ok {
			traceID = v.(string)
		}

		b.Reset()
		b.WriteString(fmt.Sprintf("%s | %s | %s | %s | %s | ",
			m["service"],
			m["ts"],
			m["level"],
			traceID,
			m["msg"],
		))

		for k, v := range m {
			switch k {
			case "service", "ts", "level", "trace_id", "msg":
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
