package persistence

import (
	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	"github.com/greendrop/todo-grpc-go-sample/domain/repository"
)

type taskPersistence struct{}

func NewTaskPersistence() repository.TaskRepository {
	return &taskPersistence{}
}

func (p taskPersistence) GetList(page int, perPage int) (*[]model.Task, error) {
	var tasks []model.Task
	offset := (page - 1) * perPage
	err := GormDB.Order("id asc").Offset(offset).Limit(perPage).Find(&tasks).Error
	return &tasks, err
}

func (p taskPersistence) GetById(id int64) (*model.Task, error) {
	var task model.Task
	err := GormDB.Where(&model.Task{Id: id}).Take(&task).Error
	return &task, err
}

func (p taskPersistence) Create(task *model.Task) (*model.Task, error) {
	err := GormDB.Create(task).Error
	return task, err
}

func (p taskPersistence) Update(task *model.Task) (*model.Task, error) {
	err := GormDB.Model(task).Updates(model.Task{Title: task.Title, Description: task.Description, Done: task.Done}).Error
	return task, err
}

func (p taskPersistence) Delete(task *model.Task) (*model.Task, error) {
	err := GormDB.Delete(task).Error
	return task, err
}
