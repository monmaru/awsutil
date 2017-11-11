package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StopInstances stops EC2 instances
func StopInstances(c *cli.Context) {
	ec2svc, err := createEC2Service(c.String("region"))
	exitIfError(err)

	instances := idFromArgs(c)
	params := &ec2.StopInstancesInput{InstanceIds: instances}
	out, err := ec2svc.StopInstances(params)
	exitIfError(err)

	for _, r := range out.StoppingInstances {
		fmt.Printf("Stopping %s.\n", *r.InstanceId)
	}
}
