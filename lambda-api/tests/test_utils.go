package tests

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/renbou/dontstress/internal/dto"
	"os"
)

const (
	contentType = "application/json"
)

var (
	baseUrl       = getEnv("BASE_URL")
	defaultLabId  = getEnv("DEFAULT_LAB_ID")
	defaultTaskId = getEnv("DEFAULT_TASK_ID")
	validToken    = getEnv("VALID_TOKEN")
)

func getEnv(key string) string {
	_ = godotenv.Load(".env")
	return os.Getenv(key)
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func labExists(labs []dto.LabDTO, labid string) bool {
	for _, l := range labs {
		if l.Id == labid {
			return true
		}
	}
	return false
}

func taskExists(tasks []dto.TaskDTO, taskid int) bool {
	for _, l := range tasks {
		if l.Id == taskid {
			return true
		}
	}
	return false
}
