package usecase_task_v1

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
)

type TaskV1DeleteTaskUseCase interface {
	Execute(task *model.Task) (*model.Task, error)
}

type taskV1DeleteTaskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskV1DeleteTaskUseCase(taskRepository repository.TaskRepository) TaskV1DeleteTaskUseCase {
	return &taskV1DeleteTaskUseCase{
		taskRepository: taskRepository,
	}
}

func (u taskV1DeleteTaskUseCase) Execute(task *model.Task) (*model.Task, error) {
	return u.taskRepository.Delete(task)
}
