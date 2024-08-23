package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var err error
	for i := 0; i < len(unmarshalTests); i++ {
		stg := &unmarshalTests[i]
		yaml := strings.ReplaceAll(stg.data, "\\n", "\n")
		yaml = strings.ReplaceAll(stg.data, "\\r", "\r")

		var jsonb []byte
		if jsonb, err = json.MarshalIndent(stg.value, "", "\t"); err != nil {
			log.Println("error:", err)
			continue
		}

		if err = os.WriteFile(fmt.Sprintf("testdata/%d.yaml", i), []byte(yaml), 0644); err != nil {
			log.Println("error:", err)
			continue
		}
		if err = os.WriteFile(fmt.Sprintf("testdata/%d.json", i), jsonb, 0644); err != nil {
			log.Println("error:", err)
		}
	}
}
