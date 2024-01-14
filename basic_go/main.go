package main

import (
	"encoding/json"
	"fmt"
)

type ProxyStatus map[string]string
type CommonProxyStatus map[string]string

type QueryResult struct {
	Status   ProxyStatus
	SqlStats []CommonProxyStatus
}

type Result struct {
	Ret   *QueryResult
	Error string
}

func main() {
	// probe写回 []*result
	byteSlice := []byte(`[{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""},{"Ret":{"Status":null,"SQLStats":null},"Error":""}]`)
	var results []Result
	if err := json.Unmarshal(byteSlice, &results); err != nil {
		panic(err)
	}
	for idx, r := range results {
		if r.Ret == nil || r.Ret.SqlStats == nil {
			fmt.Printf("idx= %d, sql stats is null", idx)
		}
	}
}
