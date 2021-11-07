package tests

import (
	"bytes"
	"encoding/json"
	"github.com/renbou/dontstress/internal/dto"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func Test_CreateLabUnauthorized(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"name": "Test lab",
	})
	response, err := http.Post(baseUrl+"/labs", contentType, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "401 Unauthorized", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	assert.Equal(t, "Unauthorized request", sb)
}

func Test_CreateTaskUnauthorized(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"name": "Test lab",
	})
	response, err := http.Post(baseUrl+"/lab/"+defaultLabId+"/tasks", contentType, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "401 Unauthorized", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	assert.Equal(t, "Unauthorized request", sb)
}

func Test_GetLabsIsArrayOfLabsOrEmpty(t *testing.T) {
	response, err := http.Get(baseUrl + "/labs")
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "200 OK", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var arr []dto.LabDTO
	err = json.Unmarshal(body, &arr)
	assert.Equal(t, nil, err)
}

func Test_CreateLabThenDelete(t *testing.T) {
	client := &http.Client{}

	// Create lab
	postBody, _ := json.Marshal(map[string]string{
		"name": "Test lab",
	})
	req, err := http.NewRequest("POST", baseUrl+"/labs", bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", validToken)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "200 OK", response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var labid string
	err = json.Unmarshal(body, &labid)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isValidUUID(labid))

	// Check if lab exists in DynamoDB
	response, err = http.Get(baseUrl + "/labs")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var labs []dto.LabDTO
	err = json.Unmarshal(body, &labs)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, labExists(labs, labid))

	// Delete lab
	req, err = http.NewRequest("DELETE", baseUrl+"/lab/"+labid, nil)
	req.Header.Add("Authorization", validToken)
	response, err = client.Do(req)
	assert.Equal(t, "200 OK", response.Status)

	// Check if lab does not exist in DynamoDB
	response, err = http.Get(baseUrl + "/labs")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &labs)
	assert.Equal(t, false, labExists(labs, labid))
}

func Test_CreateLabThenCreateTaskThenDeleteBoth(t *testing.T) {
	client := &http.Client{}

	// Create lab
	postBody, _ := json.Marshal(map[string]string{
		"name": "Test lab",
	})
	req, err := http.NewRequest("POST", baseUrl+"/labs", bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", validToken)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "200 OK", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var labid string
	err = json.Unmarshal(body, &labid)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isValidUUID(labid))

	// Check if lab exists in DynamoDB
	response, err = http.Get(baseUrl + "/labs")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var labs []dto.LabDTO
	err = json.Unmarshal(body, &labs)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, labExists(labs, labid))

	// Create task
	postBody, _ = json.Marshal(map[string]string{
		"name": "2 станка",
	})
	req, err = http.NewRequest("POST", baseUrl+"/lab/"+labid+"/tasks", bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", validToken)
	response, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "200 OK", response.Status)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var task dto.TaskDTO
	err = json.Unmarshal(body, &task)
	assert.Equal(t, nil, err)

	// Check if task exists in DynamoDB
	response, err = http.Get(baseUrl + "/lab/" + labid + "/tasks")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var tasks []dto.TaskDTO
	err = json.Unmarshal(body, &tasks)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, taskExists(tasks, task.Id))

	// Delete task
	req, err = http.NewRequest("DELETE", baseUrl+"/lab/"+labid+"/task/"+strconv.Itoa(task.Id), nil)
	req.Header.Add("Authorization", validToken)
	response, err = client.Do(req)
	assert.Equal(t, "200 OK", response.Status)

	// Check if task does not exist in DynamoDB
	response, err = http.Get(baseUrl + "/lab/" + labid + "/tasks")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &tasks)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, taskExists(tasks, task.Id))

	// Delete lab
	req, err = http.NewRequest("DELETE", baseUrl+"/lab/"+labid, nil)
	req.Header.Add("Authorization", validToken)
	response, err = client.Do(req)
	assert.Equal(t, "200 OK", response.Status)

	// Check if lab does not exist in DynamoDB
	response, err = http.Get(baseUrl + "/labs")
	assert.Equal(t, "200 OK", response.Status)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &labs)
	assert.Equal(t, false, labExists(labs, labid))
}

func Test_LangValidationUnsupported(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"lang": "bebra++",
		"data": "bebra main(){ return bebra; }",
	})
	println(baseUrl + "/lab/" + defaultLabId + "/task/" + defaultTaskId + "/test")
	response, err := http.Post(baseUrl+"/lab/"+defaultLabId+"/task/"+defaultTaskId+"/test", contentType, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "400 Bad Request", response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var payload struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &payload)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Unsupported language", payload.Message)
}

// This test uploads file to S3 bucket. There is no endpoint for deleting or auto delete,
// so don't use it until automatic deletion is not implemented.

//func Test_LangValidationSupported(t *testing.T) {
//	postBody, _ := json.Marshal(map[string]string{
//		"lang": "G++",
//		"data": "int main(){ return 0; }",
//	})
//	response, err := http.Post(baseUrl+"/lab/"+defaultLabId+"/task/"+defaultTaskId+"/test", contentType, bytes.NewBuffer(postBody))
//	if err != nil {
//		log.Fatalln(err)
//	}
//	assert.Equal(t, "200 OK", response.Status)
//}
