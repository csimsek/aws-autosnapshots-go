package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"time"
)

var descriptionPrefix string = "master-data-snapshot-"
var volumeId string = "***"

func deleteOldSnapshots(svc *ec2.EC2, days int) {
	output, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{Filters:[]*ec2.Filter{
		{
			Name:aws.String("description"),
			Values:[]*string{
				aws.String(descriptionPrefix + "-*"),
			},
		},
	}})
	if err != nil {
		fmt.Printf("%T(%v)", err, err)
		os.Exit(-1)
	}
	time := time.Now().AddDate(0, 0, -1 * days)
	for _, snapshot := range output.Snapshots {
		if snapshot.StartTime.Before(time) {
			svc.DeleteSnapshot(&ec2.DeleteSnapshotInput{SnapshotId:snapshot.SnapshotId})
		}
	}
}

func createSnapshot(svc *ec2.EC2, volumeId string) {
	description := descriptionPrefix + time.Now().Local().String()
	svc.CreateSnapshot(&ec2.CreateSnapshotInput{
		VolumeId:&volumeId,
		Description:&description,
	})
}

func main() {
	svc := ec2.New(&aws.Config{Region: aws.String("us-west-2")})
	createSnapshot(svc, volumeId)
	deleteOldSnapshots(svc, 7)
}
