package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Token save the personal token from todoist
var Token string
var todoistURL = "https://api.todoist.com/rest/v2/"

func makeRequest(method, endpoint string, data interface{}) (*http.Response, error) {
	url := todoistURL + endpoint
	body := bytes.NewBuffer([]byte{})

	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(json)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", Token)
	req.Header.Add("Authorization", bearer)

	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		defer res.Body.Close()
		str, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(str))
	}

	return res, nil
}

type taskSave struct {
	Content     string     `json:"content"`
	ProjectID   string     `json:"project_id,omitempty"`
	Order       int        `json:"order,omitempty"`
	LabelIDs    []int      `json:"label_ids,omitempty"`
	Priority    int        `json:"priority,omitempty"`
	DueString   string     `json:"due_string,omitempty"`
	DueDateTime time.Time  `json:"due_datetime,omitempty"`
	DueLang     string     `json:"due_lang,omitempty"`
}

func (ts taskSave) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	if ts.Content == "" {
		return nil, fmt.Errorf("Content is empty")
	}
	buffer.WriteString(fmt.Sprintf("\"content\":\"%s\"", ts.Content))

	if ts.ProjectID != "" {
		buffer.WriteString(fmt.Sprintf(",\"project_id\":%s", ts.ProjectID))
	}

	if ts.Order != 0 {
		buffer.WriteString(fmt.Sprintf(",\"order\":%d", ts.Order))
	}

	if !ts.DueDateTime.IsZero() {
		buffer.WriteString(",\"due_datetime\":")
		json, err := json.Marshal(ts.DueDateTime)
		if err != nil {
			return nil, err
		}
		buffer.Write(json)
	}

	if len(ts.LabelIDs) != 0 {
		buffer.WriteString(",\"label_ids\":")
		json, err := json.Marshal(ts.LabelIDs)
		if err != nil {
			return nil, err
		}
		buffer.Write(json)
	}

	if ts.Priority != 0 {
		buffer.WriteString(fmt.Sprintf(",\"priority\":%d", ts.Priority))
	}

	if ts.DueString != "" {
		buffer.WriteString(fmt.Sprintf(",\"due_string\":\"%s\"", ts.DueString))
	}

	if ts.DueLang != "" {
		buffer.WriteString(fmt.Sprintf(",\"due_lang\":\"%s\"", ts.DueLang))
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
