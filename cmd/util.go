package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func createEC2Service(region string) (*ec2.EC2, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err

	}
	conf := &aws.Config{Region: aws.String(region)}
	return ec2.New(session, conf), nil
}

func idFromArgs(c *cli.Context) []*string {
	var instances []*string
	for _, arg := range c.Args() {
		id := string(arg)
		instances = append(instances, &id)
	}
	return instances
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
