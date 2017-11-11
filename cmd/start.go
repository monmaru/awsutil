package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StartInstances starts EC2 instances
func StartInstances(c *cli.Context) {
	if c.NArg() == 0 {
		fmt.Println("Please specify the instance id")
		return
	}

	fmt.Println("Would you like to start the EC2 instances (Y/n)?")
	if !ask4confirm() {
		return
	}

	ec2svc, err := createEC2Service(region(c))
	exitIfError(err)

	instances := idFromArgs(c)
	params := &ec2.StartInstancesInput{InstanceIds: instances}
	out, err := ec2svc.StartInstances(params)
	exitIfError(err)

	for _, r := range out.StartingInstances {
		fmt.Printf("Starting %s.\n", *r.InstanceId)
	}
}
