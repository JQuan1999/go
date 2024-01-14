package history

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestAddTask(t *testing.T) {
	rdb := GetRedisClient()
	ctx, cfn := context.WithCancel(context.Background())
	defer cfn()

	go func() {
		Worker(ctx, rdb)
	}()
	users := []User{{Name: "abc", Age: 1}, {Name: "def", Age: 2}, {Name: "qwz", Age: 3}}
	AddTask(users)

	pipeline := rdb.Pipeline()
	// pipeline get
	for i := range users {
		pipeline.Get(context.Background(), userPrefix+users[i].Name)
	}
	// sleep 1 second ensure worker goroutine has set value
	time.Sleep(time.Second * 1)
	// exec
	cmds, err := pipeline.Exec(context.Background())
	if err != nil {
		t.Fatal("pipeline get key failed, err: ", err)
	}
	for index, cmd := range cmds {
		// change cmd type
		result, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		// get result
		value, err := result.Result()
		if err != nil {
			t.Fatal("exec get key failed, err: ", err)
		}
		// decode from redis value
		var decode User
		if err := json.Unmarshal([]byte(value), &decode); err != nil {
			t.Fatal("decode user failed, err: ", err)
		}
		if decode.Name != users[index].Name || decode.Age != users[index].Age {
			t.Fatal("decode user is not equal input user")
		}
	}
}
