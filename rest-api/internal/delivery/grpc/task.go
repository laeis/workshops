package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/services"
	"workshops/rest-api/internal/validators"
	"workshops/rest-api/pb"
)

type TaskHandler struct {
	service *services.TaskService
	pb.UnimplementedTaskServiceServer
}

func NewTask(service *services.TaskService) TaskHandler {
	return TaskHandler{
		service: service,
	}
}

func (t *TaskHandler) Get(ctx context.Context, tr *pb.TaskRequest) (*pb.Task, error) {
	task, err := t.service.Get(ctx, int(tr.Id))

	if err != nil {
		return nil, status.Error(codes.NotFound, "Server error")
	}

	pt := taskMapper(task)

	return pt, nil
}

func (t *TaskHandler) Fetch(ctx context.Context, tr *pb.FetchTaskRequest) (*pb.Tasks, error) {
	validator := validators.TaskValidator{}
	params := prepareMapParams(tr)

	queryFilters := filters.ValidatedTaskFilter(&validator, params)
	tasks, err := t.service.Fetch(ctx, &queryFilters)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Server error")
	}

	pbTasks := tasksMapper(&tasks)

	return pbTasks, nil
}

func (t *TaskHandler) Create(ctx context.Context, tr *pb.CreateTaskRequest) (*pb.Task, error) {
	newTask := entities.Task{
		Title:       tr.GetTitle(),
		Description: tr.GetDescription(),
		Category:    tr.GetCategory().String(),
		Date:        tr.GetDate().AsTime(),
	}

	authId, ok := ctx.Value(config.CtxAuthId).(string)

	if !ok || authId == "" {
		return nil, status.Error(codes.Unauthenticated, appError.NotAuthorized.Error())
	}

	task, err := t.service.Create(ctx, &newTask, authId)

	if err != nil {
		message := "Task not created"
		return nil, status.Error(codes.Unknown, message)
	}

	pt := taskMapper(task)

	return pt, nil
}

func (t *TaskHandler) Update(ctx context.Context, tr *pb.UpdateTaskRequest) (*pb.Task, error) {
	newTask := entities.Task{
		Title:       tr.GetTitle(),
		Description: tr.GetDescription(),
		Category:    tr.GetCategory().String(),
		Date:        tr.GetDate().AsTime(),
	}

	authId, ok := ctx.Value(config.CtxAuthId).(string)

	if !ok || authId == "" {
		return nil, status.Error(codes.Unauthenticated, appError.NotAuthorized.Error())
	}

	task, err := t.service.Update(ctx, int(tr.GetId()), &newTask)

	if err != nil {
		message := "Task not created"
		return nil, status.Error(codes.Unknown, message)
	}

	pt := taskMapper(task)

	return pt, nil
}

func (t *TaskHandler) Delete(ctx context.Context, tr *pb.TaskRequest) (*pb.Task, error) {
	success, err := t.service.Delete(ctx, int(tr.GetId()))

	if err != nil || !success {
		message := "Task didnt delete"
		return nil, status.Error(codes.Unknown, message)
	}

	return &pb.Task{}, nil
}

func taskMapper(task *entities.Task) *pb.Task {
	return &pb.Task{
		Id:          uint32(task.Id),
		UserId:      task.UserId,
		Title:       task.Title,
		Description: task.Description,
		Category:    task.Category,
		Date:        timestamppb.New(task.Date),
	}
}

func tasksMapper(tasks *entities.Tasks) *pb.Tasks {
	ts := &pb.Tasks{
		Tasks: make([]*pb.Task, len(*tasks)),
	}
	for i, v := range *tasks {
		ts.Tasks[i] = taskMapper(&v)
	}
	return ts
}

func prepareMapParams(tr *pb.FetchTaskRequest) map[string][]string {
	params := make(map[string][]string)
	if tr.GetCategory().String() != "" {
		params["category"] = []string{tr.GetCategory().String()}
	}
	if tr.GetPeriod().String() != "" {
		params["period"] = []string{tr.GetPeriod().String()}
	}
	if tr.GetOrder().String() != "" {
		params["order"] = []string{tr.GetOrder().String()}
	}
	if tr.GetOrderBy().String() != "" {
		params["order_by"] = []string{tr.GetOrderBy().String()}
	}
	return params
}
