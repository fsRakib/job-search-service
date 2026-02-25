package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "job-search-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client pb.JobServiceClient
	reader *bufio.Reader
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client = pb.NewJobServiceClient(conn)
	reader = bufio.NewReader(os.Stdin)

	fmt.Println("===========================================")
	fmt.Println("   Job Search Service - Interactive CLI")
	fmt.Println("===========================================")

	for {
		showMenu()
		choice := readInput("\nEnter your choice: ")

		switch choice {
		case "1":
			createJob()
		case "2":
			searchJobsByQuery()
		case "3":
			getJobByID()
		case "4":
			searchByLocation()
		case "5":
			searchBySkills()
		case "6":
			deleteJob()
		case "7":
			verifyDeletion()
		case "8":
			fmt.Println("\nThank you for using Job Search Service!")
			return
		default:
			fmt.Println("\nInvalid choice. Please try again.")
		}

		fmt.Println("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}

func showMenu() {
	fmt.Println("\n===========================================")
	fmt.Println("              MAIN MENU")
	fmt.Println("===========================================")
	fmt.Println("1. Create Job Listing")
	fmt.Println("2. Search Jobs by Query")
	fmt.Println("3. Get Job Details by ID")
	fmt.Println("4. Search Jobs by Location")
	fmt.Println("5. Search Jobs by Skills")
	fmt.Println("6. Delete a Job")
	fmt.Println("7. List All Jobs (Verify)")
	fmt.Println("8. Exit")
	fmt.Println("===========================================")
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func createJob() {
	fmt.Println("\n--- Create New Job Listing ---")

	title := readInput("Job Title: ")
	description := readInput("Description: ")
	company := readInput("Company: ")
	location := readInput("Location: ")
	skillsInput := readInput("Skills (comma-separated): ")
	salaryStr := readInput("Salary: ")

	salary, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		fmt.Printf("Invalid salary: %v\n", err)
		return
	}

	skills := []string{}
	if skillsInput != "" {
		skills = strings.Split(skillsInput, ",")
		for i := range skills {
			skills[i] = strings.TrimSpace(skills[i])
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateJob(ctx, &pb.CreateJobRequest{
		Title:       title,
		Description: description,
		Company:     company,
		Location:    location,
		Skills:      skills,
		Salary:      salary,
	})
	if err != nil {
		fmt.Printf("Failed to create job: %v\n", err)
		return
	}

	fmt.Println("\n✓ Job created successfully!")
	fmt.Printf("  Job ID: %s\n", resp.Id)
	fmt.Printf("  %s\n", resp.Message)
}

func searchJobsByQuery() {
	fmt.Println("\n--- Search Jobs by Query ---")

	query := readInput("Enter search query: ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SearchJobs(ctx, &pb.SearchJobsRequest{
		Query: query,
	})
	if err != nil {
		fmt.Printf("Failed to search jobs: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d job(s):\n", resp.Total)
	fmt.Println("-------------------------------------------")
	for i, job := range resp.Jobs {
		fmt.Printf("\n%d. %s", i+1, job.Title)
		if job.Score > 0 {
			fmt.Printf(" [Relevance: %.2f]\n", job.Score)
		} else {
			fmt.Printf("\n")
		}
		fmt.Printf("   Company: %s\n", job.Company)
		fmt.Printf("   Location: %s\n", job.Location)
		fmt.Printf("   Salary: $%.0f\n", job.Salary)
		fmt.Printf("   Skills: %v\n", job.Skills)
		fmt.Printf("   ID: %s\n", job.Id)
	}
}

func getJobByID() {
	fmt.Println("\n--- Get Job Details by ID ---")

	jobID := readInput("Enter Job ID: ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetJob(ctx, &pb.GetJobRequest{
		Id: jobID,
	})
	if err != nil {
		fmt.Printf("Failed to get job: %v\n", err)
		return
	}

	job := resp.Job
	fmt.Println("\n--- Job Details ---")
	fmt.Printf("ID:          %s\n", job.Id)
	fmt.Printf("Title:       %s\n", job.Title)
	fmt.Printf("Company:     %s\n", job.Company)
	fmt.Printf("Location:    %s\n", job.Location)
	fmt.Printf("Salary:      $%.0f\n", job.Salary)
	fmt.Printf("Skills:      %v\n", job.Skills)
	fmt.Printf("Description: %s\n", job.Description)
	fmt.Printf("Created At:  %s\n", job.CreatedAt)
}

func searchByLocation() {
	fmt.Println("\n--- Search Jobs by Location ---")

	location := readInput("Enter location: ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SearchJobs(ctx, &pb.SearchJobsRequest{
		Location: location,
	})
	if err != nil {
		fmt.Printf("Failed to search jobs: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d job(s) in %s:\n", resp.Total, location)
	fmt.Println("-------------------------------------------")
	for i, job := range resp.Jobs {
		fmt.Printf("\n%d. %s", i+1, job.Title)
		if job.Score > 0 {
			fmt.Printf(" [Relevance: %.2f]\n", job.Score)
		} else {
			fmt.Printf("\n")
		}
		fmt.Printf("   Company: %s\n", job.Company)
		fmt.Printf("   Salary: $%.0f\n", job.Salary)
		fmt.Printf("   ID: %s\n", job.Id)
	}
}

func searchBySkills() {
	fmt.Println("\n--- Search Jobs by Skills ---")

	skillsInput := readInput("Enter skills (comma-separated): ")

	skills := []string{}
	if skillsInput != "" {
		skills = strings.Split(skillsInput, ",")
		for i := range skills {
			skills[i] = strings.TrimSpace(skills[i])
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SearchJobs(ctx, &pb.SearchJobsRequest{
		Skills: skills,
	})
	if err != nil {
		fmt.Printf("Failed to search jobs: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d job(s) requiring: %v\n", resp.Total, skills)
	fmt.Println("-------------------------------------------")
	for i, job := range resp.Jobs {
		fmt.Printf("\n%d. %s", i+1, job.Title)
		if job.Score > 0 {
			fmt.Printf(" [Relevance: %.2f]\n", job.Score)
		} else {
			fmt.Printf("\n")
		}
		fmt.Printf("   Company: %s\n", job.Company)
		fmt.Printf("   Location: %s\n", job.Location)
		fmt.Printf("   Salary: $%.0f\n", job.Salary)
		fmt.Printf("   Skills: %v\n", job.Skills)
		fmt.Printf("   ID: %s\n", job.Id)
	}
}

func deleteJob() {
	fmt.Println("\n--- Delete Job ---")

	jobID := readInput("Enter Job ID to delete: ")
	confirm := readInput("Are you sure you want to delete this job? (yes/no): ")

	if strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.DeleteJob(ctx, &pb.DeleteJobRequest{
		Id: jobID,
	})
	if err != nil {
		fmt.Printf("Failed to delete job: %v\n", err)
		return
	}

	fmt.Println("\n✓ Job deleted successfully!")
	fmt.Printf("  %s\n", resp.Message)
}

func verifyDeletion() {
	fmt.Println("\n--- All Jobs (Verification) ---")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SearchJobs(ctx, &pb.SearchJobsRequest{
		Query: "",
	})
	if err != nil {
		fmt.Printf("Failed to list jobs: %v\n", err)
		return
	}

	fmt.Printf("\nTotal jobs in database: %d\n", resp.Total)
	fmt.Println("-------------------------------------------")
	for i, job := range resp.Jobs {
		fmt.Printf("\n%d. %s\n", i+1, job.Title)
		fmt.Printf("   Company: %s\n", job.Company)
		fmt.Printf("   Location: %s\n", job.Location)
		fmt.Printf("   ID: %s\n", job.Id)
	}
}
