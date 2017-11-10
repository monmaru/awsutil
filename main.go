package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

const region = "ap-northeast-1"

func main() {
	go spinner(100 * time.Millisecond)

	session, err := session.NewSession()
	exitIfError(err)
	ec2svc := ec2.New(session, &aws.Config{Region: aws.String(region)})

	out, err := ec2svc.DescribeInstances(nil)
	exitIfError(err)

	header := []string{
		"Name",
		"InstanceId",
		"InstanceType",
		"AZ",
		"PrivateIP",
		"PublicIP",
		"Status",
	}
	var data [][]string

	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			var tagName string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					tagName = *t.Value
				}
			}

			var privateIP, publicIP string
			if i.PublicIpAddress == nil {
				publicIP = "-"
			} else {
				publicIP = *i.PublicIpAddress
			}
			if i.PrivateIpAddress == nil {
				privateIP = "-"
			} else {
				privateIP = *i.PrivateIpAddress
			}

			data = append(data, []string{
				tagName,
				*i.InstanceId,
				*i.InstanceType,
				*i.Placement.AvailabilityZone,
				privateIP,
				publicIP,
				*i.State.Name,
			})
		}
	}

	fmt.Print("\033[2K")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.Render()
}

func spinner(delay time.Duration) {
	fmt.Print("\033[?25l")
	for {
		for _, r := range `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
