package server

import (
	"context"
	"errors"

	"github.com/greendrop/todo-grpc-go-sample/domain/model"
	proto_task_v1 "github.com/greendrop/todo-grpc-go-sample/proto/task/v1"
	usecase_task_v1 "github.com/greendrop/todo-grpc-go-sample/usecase/task/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type taskV1ServiceServer struct {
	taskV1GetTaskListUseCase usecase_task_v1.TaskV1GetTaskListUseCase
	taskV1GetTaskUseCase     usecase_task_v1.TaskV1GetTaskUseCase
	taskV1CreateTaskUseCase  usecase_task_v1.TaskV1CreateTaskUseCase
	taskV1UpdateTaskUseCase  usecase_task_v1.TaskV1UpdateTaskUseCase
	taskV1DeleteTaskUseCase  usecase_task_v1.TaskV1DeleteTaskUseCase
	proto_task_v1.UnimplementedTaskServiceServer
}

func NewTaskV1ServiceServer(
	taskV1GetTaskListUseCase usecase_task_v1.TaskV1GetTaskListUseCase,
	taskV1GetTaskUseCase usecase_task_v1.TaskV1GetTaskUseCase,
	taskV1CreateTaskUseCase usecase_task_v1.TaskV1CreateTaskUseCase,
	taskV1UpdateTaskUseCase usecase_task_v1.TaskV1UpdateTaskUseCase,
	taskV1DeleteTaskUseCase usecase_task_v1.TaskV1DeleteTaskUseCase,
) *taskV1ServiceServer {
	return &taskV1ServiceServer{
		taskV1GetTaskListUseCase: taskV1GetTaskListUseCase,
		taskV1GetTaskUseCase:     taskV1GetTaskUseCase,
		taskV1CreateTaskUseCase:  taskV1CreateTaskUseCase,
		taskV1UpdateTaskUseCase:  taskV1UpdateTaskUseCase,
		taskV1DeleteTaskUseCase:  taskV1DeleteTaskUseCase,
	}
}

func (s taskV1ServiceServer) GetTaskList(ctx context.Context, getTaskListRequest *proto_task_v1.GetTaskListRequest) (*proto_task_v1.GetTaskListResponse, error) {
	tasks, err := s.taskV1GetTaskListUseCase.Execute(getTaskListRequest.Page, getTaskListRequest.PerPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, codes.NotFound.String())
		} else {
			return nil, status.Errorf(codes.Internal, codes.Internal.String())
		}
	}

	var responseTasks []*proto_task_v1.GetTaskListResponse_Task
	for _, task := range *tasks {
		responseTasks = append(responseTasks, &proto_task_v1.GetTaskListResponse_Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Done:        task.Done,
			CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05+09:00"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05+09:00"),
		})
	}

	return &proto_task_v1.GetTaskListResponse{
		Tasks: responseTasks,
	}, nil
}

func (s taskV1ServiceServer) GetTask(ctx context.Context, getTaskRequest *proto_task_v1.GetTaskRequest) (*proto_task_v1.GetTaskResponse, error) {
	task, err := s.taskV1GetTaskUseCase.Execute(getTaskRequest.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, codes.NotFound.String())
		} else {
			return nil, status.Errorf(codes.Internal, codes.Internal.String())
		}
	}

	return &proto_task_v1.GetTaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05+09:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05+09:00"),
	}, nil
}

func (s taskV1ServiceServer) CreateTask(ctx context.Context, createTaskRequest *proto_task_v1.CreateTaskRequest) (*proto_task_v1.CreateTaskResponse, error) {
	task := &model.Task{
		Title:       createTaskRequest.Title,
		Description: createTaskRequest.Description,
		Done:        createTaskRequest.Done,
	}

	task, err := s.taskV1CreateTaskUseCase.Execute(task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &proto_task_v1.CreateTaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05+09:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05+09:00"),
	}, nil
}

func (s taskV1ServiceServer) UpdateTask(ctx context.Context, updateTaskRequest *proto_task_v1.UpdateTaskRequest) (*proto_task_v1.UpdateTaskResponse, error) {
	task := &model.Task{
		Id:          updateTaskRequest.Id,
		Title:       updateTaskRequest.Title,
		Description: updateTaskRequest.Description,
		Done:        updateTaskRequest.Done,
	}

	task, err := s.taskV1UpdateTaskUseCase.Execute(task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &proto_task_v1.UpdateTaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05+09:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05+09:00"),
	}, nil
}

func (s taskV1ServiceServer) DeleteTask(ctx context.Context, deleteTaskRequest *proto_task_v1.DeleteTaskRequest) (*proto_task_v1.DeleteTaskResponse, error) {
	task := &model.Task{
		Id: deleteTaskRequest.Id,
	}

	_, err := s.taskV1DeleteTaskUseCase.Execute(task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &proto_task_v1.DeleteTaskResponse{}, nil
}
