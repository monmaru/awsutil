package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StopInstances stops EC2 instances
func StopInstances(c *cli.Context) {
	if c.NArg() == 0 {
		fmt.Println("Please specify the instance id")
		return
	}

	fmt.Println("Would you like to stop the EC2 instances (Y/n)?")
	if !ask4confirm() {
		return
	}

	ec2svc, err := createEC2Service(region(c))
	exitIfError(err)

	instances := idFromArgs(c)
	params := &ec2.StopInstancesInput{InstanceIds: instances}
	out, err := ec2svc.StopInstances(params)
	exitIfError(err)

	for _, r := range out.StoppingInstances {
		fmt.Printf("Stopping %s.\n", *r.InstanceId)
	}
}
