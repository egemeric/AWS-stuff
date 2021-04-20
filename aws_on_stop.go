package main

import (
	"aws_api/models"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
)

func main() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-central-1")})
	if err != nil {
		fmt.Println("there was an error in creating session:", err.Error())
		log.Fatal(err.Error())
	}
	// Create new EC2 client
	ec2Svc := ec2.New(sess)
	Instance := GetInstances(ec2Svc)
	fmt.Println(Instance)
	//InstanceStart(Instance[0].Id, ec2Svc)
	InstanceStop(Instance[0].Id, ec2Svc)
}
func GetInstances(ec2Svc *ec2.EC2) (custom_inst []models.CustomInstance) {
	var ip_addr string
	resp, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	for idx, res := range resp.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID:", *inst.InstanceId)
			fmt.Println("                     State: ", *inst.State.Name)
			if *inst.State.Name == "running" {
				ip_addr = *inst.PublicIpAddress
				fmt.Println("                 public ip: ", *inst.PublicIpAddress)
			}
			fmt.Println("                  key-name: ", *inst.KeyName)
			custom_inst = append(custom_inst, models.CustomInstance{Ip: ip_addr, Id: *inst.InstanceId, Status: *inst.State.Name, KeyName: *inst.KeyName, Name: *inst.PublicDnsName})
		}
	}
	return custom_inst
}

func InstanceStart(I_id string, ec2Svc *ec2.EC2) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(I_id),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2Svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = ec2Svc.StartInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Success", result.StartingInstances)
		}
	} else { // This could be due to a lack of permissions
		fmt.Println("Error", err)
	}
}

func InstanceStop(I_id string, ec2Svc *ec2.EC2) {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(I_id),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2Svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2Svc.StopInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Success", result.StoppingInstances)
		}
	} else {
		fmt.Println("Error", err)
	}
}
