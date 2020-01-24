package graph

import (
	"container/list"
	"fmt"
	"github.com/sotomskir/jira-cli/jiraApi/models"
)

type Graph struct {
	VerticesCount uint
	Vertices      map[string]*Vertex
}

type Vertex struct {
	Id      string
	Friends map[string]*Vertex
	Data    interface{}
}

// Constructor
func New(v uint) *Graph {
	graph := &Graph{
		VerticesCount: v,
		Vertices:      make(map[string]*Vertex),
	}
	return graph
}

func NewFromWorkflow(workflow *models.Workflow) *Graph {
	graph := New(uint(len(workflow.Layout.Statuses)))
	for _, status := range workflow.Layout.Statuses {
		graph.AddVertex(&Vertex{
			Id:      status.Id,
			Friends: make(map[string]*Vertex),
			Data:    status,
		})
	}
	for _, transition := range workflow.Layout.Transitions {
		graph.AddEdge(transition.SourceId, transition.TargetId)
	}
	return graph
}

// Function to add an edge into the graph
func (g Graph) AddEdge(fromId string, toId string) {
	g.Vertices[fromId].Friends[toId] = g.Vertices[toId]
}

// Function to add a vertex into the graph
func (g Graph) AddVertex(vertex *Vertex) {
	g.Vertices[vertex.Id] = vertex
}

// prints BFS traversal from a given source s
func (g Graph) BFS(fromId string, toId string) {
	startVertex := g.Vertices[fromId]
	// Mark all the vertices as not visited(By default
	// set as false)
	visited := make(map[string]bool)

	// Create a queue for BFS
	queue := list.New()
	queue.PushBack(startVertex)
	path := list.New()
	path.PushBack(startVertex)
	// Mark the current node as visited and enqueue it
	visited[fromId] = true
	for queue.Len() > 0 {
		// Dequeue a vertex from queue and print it
		qnode := queue.Front()

		path.PushBack(qnode.Value)
		// iterate through all of its friends
		// mark the visited nodes; enqueue the non-visted
		for id, vertex := range qnode.Value.(*Vertex).Friends {
			if _, ok := visited[id]; !ok {
				visited[id] = true
				queue.PushBack(vertex)
				if id == toId {
					path.PushBack(vertex)
					queue = list.New()
					break
				}
			}
		}
		path.PushBack(qnode.Value)
		queue.Remove(qnode)
	}

	for path.Len() > 0 {
		vertex := path.Front()
		fmt.Println(vertex.Value.(*Vertex).Id)
		path.Remove(vertex)
	}
}
