package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var jstLoc *time.Location

func init() {
	jstLoc, _ = time.LoadLocation("Asia/Tokyo")
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) (string, error) {
	ec2svc, err := initEC2service()
	if err != nil {
		return "initialize error", err
	}

	if err := run(ec2svc); err != nil {
		log.Printf("error: %+v\n", err)
		return "run error", err
	}
	return "success", nil
}

func initEC2service() (ec2iface.EC2API, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	conf := &aws.Config{Region: aws.String("ap-northeast-1")}
	ec2svc := ec2.New(session, conf)
	return ec2svc, nil
}

func run(ec2svc ec2iface.EC2API) error {
	now := time.Now().In(jstLoc)
	log.Printf("Now: %s", now.String())

	// List
	out, err := ec2svc.DescribeInstances(nil)
	if err != nil {
		return err
	}

	// Start
	if weekday(now) {
		if err := startInstances(ec2svc, out.Reservations, now); err != nil {
			return err
		}
	}

	// Stop
	if err := stopInstances(ec2svc, out.Reservations, now); err != nil {
		return err
	}

	return nil
}

func weekday(t time.Time) bool {
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return false
	default:
		return true
	}
}

func startInstances(svc ec2iface.EC2API, reservations []*ec2.Reservation, now time.Time) error {
	log.Println("startInstances called")
	const tag = "PowerOn"
	var ids []*string
	for _, ins := range instancesByStatus(reservations, tag, "stopped") {
		tagValue := getTagValue(ins, tag)
		if isTarget(tagValue, now) {
			ids = append(ids, ins.InstanceId)
		}
	}

	if len(ids) == 0 {
		log.Println("There are no target instances.")
		return nil
	}

	params := &ec2.StartInstancesInput{InstanceIds: ids}
	out, err := svc.StartInstances(params)
	if err != nil {
		return err
	}

	for _, r := range out.StartingInstances {
		log.Printf("Starting %s.\n", *r.InstanceId)
	}
	return nil
}

func stopInstances(svc ec2iface.EC2API, reservations []*ec2.Reservation, now time.Time) error {
	log.Println("stopInstances called")
	const tag = "PowerOff"
	var ids []*string
	for _, ins := range instancesByStatus(reservations, tag, "running") {
		tagValue := getTagValue(ins, tag)
		if isTarget(tagValue, now) {
			ids = append(ids, ins.InstanceId)
		}
	}

	if len(ids) == 0 {
		log.Println("There are no target instances.")
		return nil
	}

	params := &ec2.StopInstancesInput{InstanceIds: ids}
	out, err := svc.StopInstances(params)
	if err != nil {
		return err
	}

	for _, r := range out.StoppingInstances {
		log.Printf("Stopping %s.\n", *r.InstanceId)
	}
	return nil
}

func instancesByStatus(reservations []*ec2.Reservation, tag, status string) []*ec2.Instance {
	var is []*ec2.Instance
	for _, r := range reservations {
		for _, i := range r.Instances {
			for _, t := range i.Tags {
				if *t.Key == tag && *i.State.Name == status {
					is = append(is, i)
				}
			}
		}
	}
	return is
}

func isTarget(tagValue string, now time.Time) bool {
	fewMinutesAgo := now.Add(-40 * time.Minute)
	ss := strings.Split(tagValue, ":")

	if len(ss) != 2 {
		log.Printf("invalid tag value %s\n", tagValue)
		return false
	}

	hour, err := strconv.Atoi(ss[0])
	if err != nil {
		log.Printf("invalid tag value %s: %+v\n", tagValue, err)
		return false
	}

	minute, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Printf("invalid tag value %s: %+v\n", tagValue, err)
		return false
	}

	schedule := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, now.Second(), now.Nanosecond(), jstLoc)
	return schedule.Before(now) && schedule.After(fewMinutesAgo)
}

func getTagValue(i *ec2.Instance, tag string) string {
	for _, t := range i.Tags {
		if *t.Key == tag {
			return *t.Value
		}
	}
	return ""
}
