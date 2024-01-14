package history

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ngaut/log"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	ttl   time.Duration
	key   string
	value string
}

type User struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

var (
	tasksCh chan []Task
)

const userPrefix = "pipeline_user_"

func init() {
	tasksCh = make(chan []Task, 512)
}

func AddTask(users []User) {
	// 设置key数组
	keys := make([]string, 0)
	for _, user := range users {
		keys = append(keys, userPrefix+user.Name)
	}

	ttl := time.Minute * 10

	var tasks []Task
	for i := 0; i < len(users); i++ {
		// 序列化user设置为value
		data, err := json.Marshal(users[i])
		if err != nil {
			log.Errorf("marshal users failed, err: %v", err)
			continue
		}
		// append到task数组
		tasks = append(tasks, Task{key: keys[i], value: string(data), ttl: ttl})
	}

	// 添加到taskCh
	tasksCh <- tasks
}

func Worker(ctx context.Context, rds *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		case tasks := <-tasksCh:
			if err := PipelineSet(ctx, rds, tasks); err != nil {
				log.Errorf("pipeline set failed, err: %v", err)
			}
		}
	}
}

func PipelineSet(ctx context.Context, rds *redis.Client, tasks []Task) error {
	setCtx, cfn := context.WithTimeout(ctx, time.Second*4)
	defer cfn()

	// 使用pipelined自动调用exec函数
	cmds, err := rds.Pipelined(setCtx, func(p redis.Pipeliner) error {
		for index := range tasks {
			item := tasks[index]
			p.Set(ctx, item.key, item.value, item.ttl)
		}
		return nil
	})

	if err != nil {
		log.Errorf("pipeline execute failed, err: %v", err)
	}
	for index, cmd := range cmds {
		if cmd.Err() != nil {
			log.Errorf("pipeline key %v error: %v", tasks[index].key, cmd.Err())
		}
	}
	return nil
}
