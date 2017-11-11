package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StartInstances starts EC2 instances
func StartInstances(c *cli.Context) {
	ec2svc, err := createEC2Service(c.String("region"))
	exitIfError(err)

	instances := idFromArgs(c)
	params := &ec2.StartInstancesInput{InstanceIds: instances}
	out, err := ec2svc.StartInstances(params)
	exitIfError(err)

	for _, r := range out.StartingInstances {
		fmt.Printf("Starting %s.\n", *r.InstanceId)
	}
}
