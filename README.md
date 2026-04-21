Mini Managed Agents - 永不失败的Agent
===

我在模型调用和工具调用中增加了随机的异常:
```go
func randomExit() {
	if time.Now().Nanosecond()%2 == 0 {
		debug.PrintStack()
		os.Exit(1)
	}
}
```
但绝大多数情况下,Agent都能成功的完成任务

## 环境需求
- Go: 1.25
- Docker(Docker Compose)


## 启动方式
1. 复制一份`.env`文件并修改为你自己的API信息
```sh
cp .env.example .env
```

2. 启动Agent Worker和Temporal组件
```sh
docker compose up -d
```

3. 向AgentWorker发送任务
```sh
set -a && source .env && set +a && go run ./cmd/start_workflow
```

4. 进入Temporal UI查看Workflow状态: http://localhost:8233/namespaces/default/workflows