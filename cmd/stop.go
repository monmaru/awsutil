package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StopInstance stops EC2 instance
func StopInstance(c *cli.Context) {
	ec2svc, err := createEC2Service(region(c), profile(c))
	exitIfError(err)

	out, err := ec2svc.DescribeInstances(nil)
	exitIfError(err)

	instances := instancesByStatus(out.Reservations, "running")
	if len(instances) == 0 {
		fmt.Println("There are no running instances.")
		return
	}

	name := selectInstancePrompt(instances)
	if len(name) == 0 {
		fmt.Println("The command has been canceled.")
		return
	}

	fmt.Printf("Would you like to stop the %s (Y/n)?\n", name)
	if !ask4confirm() {
		return
	}

	id := findIDByName(instances, name)
	params := &ec2.StopInstancesInput{InstanceIds: []*string{id}}
	stopOut, err := ec2svc.StopInstances(params)
	exitIfError(err)

	for _, r := range stopOut.StoppingInstances {
		fmt.Printf("Stopping %s.\n", *r.InstanceId)
	}
}
