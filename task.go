package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Task is a model of todoist project entity
type Task struct {
	ID           string `json:"id"`
	CommentCount int    `json:"comment_count"`
	Completed    bool   `json:"completed"`
	Content      string `json:"content"`
	Indent       int    `json:"indent"`
	LabelIDs     []int  `json:"label_ids"`
	Order        int    `json:"order"`
	Priority     int    `json:"priority"`
	ProjectID    string `json:"project_id"`
	Due          Due    `json:"due"`
}

// Due is a model of todoist project entity
type Due struct {
	String   string    `json:"string"`
	Date     string    `json:"date"`
	Datetime time.Time `json:"datetime,omitempty"`
	Timezone string    `json:"timezone"`
}

func (t Task) taskSave() taskSave {
	return taskSave{
		t.Content,
		t.ProjectID,
		t.Order,
		t.LabelIDs,
		t.Priority,
		t.Due.String,
		t.Due.Datetime,
		"en",
	}
}

func decodeTask(body io.ReadCloser) (Task, error) {
	defer body.Close()
	decoder := json.NewDecoder(body)
	var task Task

	if err := decoder.Decode(&task); err != nil {
		return Task{}, err
	}
	return task, nil
}

// QueryParam is a map[string]string to build http query
type QueryParam map[string]string

func (qp QueryParam) String() string {
	if len(qp) == 0 {
		return ""
	}

	ret := "?"
	for key, value := range qp {
		if ret != "?" {
			ret = ret + "&"
		}
		ret = ret + key + "=" + value
	}

	return ret
}

// ListTask return all task, you can filter using QueryParam
// See documentation: https://developer.todoist.com/rest/v8/#get-tasks
func ListTask(qp QueryParam) ([]Task, error) {
	path := fmt.Sprintf("tasks%s", qp)
	res, err := makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return []Task{}, err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var tasks []Task

	if err := decoder.Decode(&tasks); err != nil {
		return []Task{}, err
	}

	return tasks, nil
}

// GetTask return a task by id
func GetTask(id string) (Task, error) {
	path := fmt.Sprintf("tasks/%s", id)
	res, err := makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return Task{}, err
	}

	return decodeTask(res.Body)
}

// CreateTask create a new task
func CreateTask(task Task) (Task, error) {
	res, err := makeRequest(http.MethodPost, "tasks", task.taskSave())
	if err != nil {
		return Task{}, err
	}

	return decodeTask(res.Body)
}

// Delete remove a task
func (t Task) Delete() error {
	path := fmt.Sprintf("tasks/%s", t.ID)
	_, err := makeRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return nil
}

// Update a task
func (t Task) Update() error {
	path := fmt.Sprintf("tasks/%s", t.ID)
	_, err := makeRequest(http.MethodPost, path, t.taskSave())
	if err != nil {
		return err
	}

	return nil
}

// Close mask task as done
func (t Task) Close() error {
	path := fmt.Sprintf("tasks/%s/close", t.ID)
	_, err := makeRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return nil
}

// Reopen a task
func (t Task) Reopen() error {
	path := fmt.Sprintf("tasks/%s/reopen", t.ID)
	_, err := makeRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return nil
}
