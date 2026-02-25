package grpc

import (
	"context"
	"job-search-service/internal/service"
	pb "job-search-service/proto"
	"log"
)

type JobHandler struct {
	pb.UnimplementedJobServiceServer
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{
		service: service,
	}
}

func (h *JobHandler) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	log.Printf("Creating job: %s", req.Title)

	id, err := h.service.CreateJob(
		ctx,
		req.Title,
		req.Description,
		req.Company,
		req.Location,
		req.Skills,
		req.Salary,
	)
	if err != nil {
		log.Printf("Error creating job: %v", err)
		return nil, err
	}

	return &pb.CreateJobResponse{
		Id:      id,
		Message: "Job created successfully",
	}, nil
}

func (h *JobHandler) SearchJobs(ctx context.Context, req *pb.SearchJobsRequest) (*pb.SearchJobsResponse, error) {
	log.Printf("Searching jobs with query: %s", req.Query)

	jobs, err := h.service.SearchJobs(ctx, req.Query, req.Location, req.Skills)
	if err != nil {
		log.Printf("Error searching jobs: %v", err)
		return nil, err
	}

	pbJobs := make([]*pb.Job, 0, len(jobs))
	for _, job := range jobs {
		pbJobs = append(pbJobs, &pb.Job{
			Id:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Company:     job.Company,
			Location:    job.Location,
			Skills:      job.Skills,
			Salary:      job.Salary,
			CreatedAt:   job.CreatedAt.Format("2006-01-02"),
		})
	}

	return &pb.SearchJobsResponse{
		Jobs:  pbJobs,
		Total: int32(len(pbJobs)),
	}, nil
}

func (h *JobHandler) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	log.Printf("Getting job with ID: %s", req.Id)

	job, err := h.service.GetJob(ctx, req.Id)
	if err != nil {
		log.Printf("Error getting job: %v", err)
		return nil, err
	}

	return &pb.GetJobResponse{
		Job: &pb.Job{
			Id:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Company:     job.Company,
			Location:    job.Location,
			Skills:      job.Skills,
			Salary:      job.Salary,
			CreatedAt:   job.CreatedAt.Format("2006-01-02"),
		},
	}, nil
}

func (h *JobHandler) DeleteJob(ctx context.Context, req *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	log.Printf("Deleting job with ID: %s", req.Id)

	err := h.service.DeleteJob(ctx, req.Id)
	if err != nil {
		log.Printf("Error deleting job: %v", err)
		return nil, err
	}

	return &pb.DeleteJobResponse{
		Message: "Job deleted successfully",
	}, nil
}
