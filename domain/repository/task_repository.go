package repository

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
)

type TaskRepository interface {
	GetList(page int, parPage int) (*[]model.Task, error)
	GetById(id int64) (*model.Task, error)
	Create(task *model.Task) (*model.Task, error)
	Update(task *model.Task) (*model.Task, error)
	Delete(task *model.Task) (*model.Task, error)
}
