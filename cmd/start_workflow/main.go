package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	taskQueue := os.Getenv("TEMPORAL_TASK_QUEUE")
	query := "What is the weather in Louisiana and what is my current ip address? Where am I located?"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	c, err := client.Dial(client.Options{HostPort: temporalAddress})
	if err != nil {
		log.Fatalf("create temporal client: %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "agentic-loop-id-" + uuid.NewString(),
		TaskQueue: taskQueue,
	}

	run, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "AgentWorkflow", query)
	if err != nil {
		log.Fatalf("execute workflow: %v", err)
	}

	var result string
	if err := run.Get(context.Background(), &result); err != nil {
		log.Fatalf("get workflow result: %v", err)
	}

	fmt.Printf("Result: %s\n", result)
}
