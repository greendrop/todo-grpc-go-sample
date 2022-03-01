package usecase_task_v1

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
	"github.com/greendrop/todo-grpc-go-sample/util/pagination"
)

type TaskV1GetTaskListUseCase interface {
	Execute(page *int32, perPage *int32) (*[]model.Task, error)
}

type taskV1GetTaskListUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskV1GetTaskListUseCase(taskRepository repository.TaskRepository) TaskV1GetTaskListUseCase {
	return &taskV1GetTaskListUseCase{
		taskRepository: taskRepository,
	}
}

func (u taskV1GetTaskListUseCase) Execute(page *int32, perPage *int32) (*[]model.Task, error) {
	return u.taskRepository.GetList(pagination.ParsePageParam(page), pagination.ParsePerPageParam(perPage))
}
