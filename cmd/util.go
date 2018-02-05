package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func createEC2Service(region, profile string) (*ec2.EC2, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	conf := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile),
	}
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

func region(c *cli.Context) string {
	return c.String("region")
}

func profile(c *cli.Context) string {
	return c.String("profile")
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ask4confirm() bool {
	var s string
	_, err := fmt.Scanln(&s)
	if err != nil {
		log.Fatal(err)
	}
	s = strings.ToLower(strings.TrimSpace(s))

	if s == "y" {
		return true
	} else if s == "n" {
		return false
	} else {
		fmt.Println("Please type y or n and then press enter:")
		return ask4confirm()
	}
}
