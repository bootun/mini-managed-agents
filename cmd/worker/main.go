package main

import (
	"log"
	"os"

	"github.com/bootun/mini-managed-agents/activities"
	"github.com/bootun/mini-managed-agents/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	taskQueue := os.Getenv("TEMPORAL_TASK_QUEUE")
	c, err := client.Dial(client.Options{HostPort: temporalAddress})
	if err != nil {
		log.Fatalf("create temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})
	// 注册workflow和activity
	w.RegisterWorkflowWithOptions((&workflows.AgentWorkflow{}).Run, workflow.RegisterOptions{Name: "AgentWorkflow"})
	w.RegisterActivity(activities.CreateOpenAIResponse)
	w.RegisterActivity(activities.InvokeTool)

	log.Printf("worker listening on %s with task queue %s", temporalAddress, taskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
