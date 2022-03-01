package usecase_task_v1

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
)

type TaskV1CreateTaskUseCase interface {
	Execute(task *model.Task) (*model.Task, error)
}

type taskV1CreateTaskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskV1CreateTaskUseCase(taskRepository repository.TaskRepository) TaskV1CreateTaskUseCase {
	return &taskV1CreateTaskUseCase{
		taskRepository: taskRepository,
	}
}

func (u taskV1CreateTaskUseCase) Execute(task *model.Task) (*model.Task, error) {
	return u.taskRepository.Create(task)
}
