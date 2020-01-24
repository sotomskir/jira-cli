package graph

import (
	"container/list"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi/models"
	"os"
	"strings"
)

type Graph struct {
	VerticesCount uint
	Vertices      map[string]*Vertex
	Transitions   map[string]map[string]*models.Transition
}

type Vertex struct {
	Id      string
	Friends map[string]*Vertex
	Status  models.Status
}

type PathNode struct {
	Status         *models.Status
	NextTransition *models.Transition
}

// Constructor
func New(v uint) *Graph {
	graph := &Graph{
		VerticesCount: v,
		Vertices:      make(map[string]*Vertex),
		Transitions:   make(map[string]map[string]*models.Transition),
	}
	return graph
}

func NewFromWorkflow(workflow *models.Workflow) *Graph {
	graph := New(uint(len(workflow.Layout.Statuses)))
	for _, status := range workflow.Layout.Statuses {
		graph.AddVertex(&Vertex{
			Id:      status.Id,
			Friends: make(map[string]*Vertex),
			Status:  status,
		})
		graph.Transitions[status.Id] = make(map[string]*models.Transition)
	}
	for _, transition := range workflow.Layout.Transitions {
		t := transition
		graph.Transitions[transition.SourceId][transition.TargetId] = &t
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
func (g Graph) BFS(fromId string) map[string]*Vertex {
	startVertex := g.Vertices[fromId]
	// Mark all the vertices as not visited(By default
	// set as false)
	visited := make(map[string]*Vertex)

	// Create a queue for BFS
	queue := list.New()
	queue.PushBack(startVertex)

	// Mark the current node as visited and enqueue it
	visited[fromId] = startVertex
	for queue.Len() > 0 {
		// Dequeue a vertex from queue and print it
		qnode := queue.Front()

		// iterate through all of its friends
		// mark the visited nodes; enqueue the non-visted
		for id, vertex := range qnode.Value.(*Vertex).Friends {
			if _, ok := visited[id]; !ok {
				visited[id] = qnode.Value.(*Vertex)
				queue.PushBack(vertex)
			}
		}
		queue.Remove(qnode)
	}
	return visited
}

func (g Graph) FindPath(fromId string, toId string) *list.List {
	bfs := g.BFS(fromId)
	path := list.New()
	status := g.Vertices[toId].Status
	path.PushFront(PathNode{
		Status:         &status,
		NextTransition: nil,
	})

	currentId := toId
	for currentId != fromId {
		qnode := bfs[currentId]
		currentId = qnode.Id
		toStatus := path.Front().Value.(PathNode).Status
		fromStatus := qnode.Status
		path.PushFront(PathNode{
			Status:         &fromStatus,
			NextTransition: g.Transitions[fromStatus.Id][toStatus.Id],
		})
	}

	return path
}

func (g Graph) FindVertexByName(name string) *Vertex {
	for _, v := range g.Vertices {
		if strings.ToLower(strings.TrimSpace(v.Status.Name)) == strings.ToLower(strings.TrimSpace(name)) {
			return v
		}
	}
	logrus.Errorf(fmt.Sprintf("unknown status: %s", name))
	os.Exit(1)
	return nil
}
