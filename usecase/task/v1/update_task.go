package usecase_task_v1

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
)

type TaskV1UpdateTaskUseCase interface {
	Execute(task *model.Task) (*model.Task, error)
}

type taskV1UpdateTaskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskV1UpdateTaskUseCase(taskRepository repository.TaskRepository) TaskV1UpdateTaskUseCase {
	return &taskV1UpdateTaskUseCase{
		taskRepository: taskRepository,
	}
}

func (u taskV1UpdateTaskUseCase) Execute(task *model.Task) (*model.Task, error) {
	return u.taskRepository.Update(task)
}
