package cmd

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// ListInstances lists all EC2 instance infomation
func ListInstances(c *cli.Context) {
	ec2svc, err := createEC2Service(region(c), profile(c))
	exitIfError(err)

	out, err := ec2svc.DescribeInstances(nil)
	exitIfError(err)

	header := []string{
		"Name",
		"InstanceID",
		"InstanceType",
		"AZ",
		"PrivateIP",
		"PublicIP",
		"Status",
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(dataFromReservations(out.Reservations))
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.Render()
}

func dataFromReservations(reservations []*ec2.Reservation) [][]string {
	var data [][]string
	for _, r := range reservations {
		for _, i := range r.Instances {
			var name string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					name = *t.Value
				}
			}

			var privateIP, publicIP string
			if i.PublicIpAddress == nil {
				publicIP = "-"
			} else {
				publicIP = *i.PublicIpAddress
			}
			if i.PrivateIpAddress == nil {
				privateIP = "-"
			} else {
				privateIP = *i.PrivateIpAddress
			}

			data = append(data, []string{
				name,
				*i.InstanceId,
				*i.InstanceType,
				*i.Placement.AvailabilityZone,
				privateIP,
				publicIP,
				*i.State.Name,
			})
		}
	}
	return data
}
