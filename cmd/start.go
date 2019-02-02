package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

// StartInstance starts EC2 instance
func StartInstance(c *cli.Context) {
	ec2svc, err := createEC2Service(region(c), profile(c))
	exitIfError(err)

	out, err := ec2svc.DescribeInstances(nil)
	exitIfError(err)

	//	https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	instances := instancesByStatus(out.Reservations, "stopped")
	if len(instances) == 0 {
		fmt.Println("There are no stopped instances.")
		return
	}

	name := selectInstancePrompt(instances)
	if len(name) == 0 {
		fmt.Println("The command has been canceled.")
		return
	}

	fmt.Printf("Would you like to start the %s (Y/n)?\n", name)
	if !ask4confirm() {
		return
	}

	id := findIDByName(instances, name)
	params := &ec2.StartInstancesInput{InstanceIds: []*string{id}}
	startOut, err := ec2svc.StartInstances(params)
	exitIfError(err)

	for _, r := range startOut.StartingInstances {
		fmt.Printf("Starting %s.\n", *r.InstanceId)
	}
}
