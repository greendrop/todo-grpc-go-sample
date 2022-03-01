package usecase_task_v1

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
)

type TaskV1GetTaskUseCase interface {
	Execute(id int64) (*model.Task, error)
}

type taskV1GetTaskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskV1GetTaskUseCase(taskRepository repository.TaskRepository) TaskV1GetTaskUseCase {
	return &taskV1GetTaskUseCase{
		taskRepository: taskRepository,
	}
}

func (u taskV1GetTaskUseCase) Execute(id int64) (*model.Task, error) {
	return u.taskRepository.GetById(id)
}
