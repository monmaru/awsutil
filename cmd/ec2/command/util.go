package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/urfave/cli"
)

func createEC2Service(region, profile string) (ec2iface.EC2API, error) {
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

func instancesByStatus(reservations []*ec2.Reservation, status string) []*ec2.Instance {
	var is []*ec2.Instance
	for _, r := range reservations {
		for _, i := range r.Instances {
			if *i.State.Name == status {
				is = append(is, i)
			}
		}
	}
	return is
}

func findIDByName(instances []*ec2.Instance, name string) *string {
	for _, i := range instances {
		for _, t := range i.Tags {
			if *t.Key == "Name" && name == *t.Value {
				return i.InstanceId
			}
		}
	}
	return nil
}

func idFromArgs(c *cli.Context) []*string {
	var ids []*string
	for _, arg := range c.Args() {
		id := string(arg)
		ids = append(ids, &id)
	}
	return ids
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
