// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/spf13/cobra"
)

var svc *ec2.EC2

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "creates terraform for ec2 instances",
	Long: `
	terragen ec2 - will generate terraform for all ec2 instances
	terragen ec2 {{ec2.arn}} - will generate terraform for the ec2 arn specified
`,
	Run: func(cmd *cobra.Command, args []string) {
		ec2 := instance{}
		ec2.describeInstance()
		inVPC := ec2.IsInstanceInVPC()
		fmt.Println(inVPC)
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
	svc = ec2.New(session.New())
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type instance struct {
	ami                   string
	associatePublicIP     bool
	availabilityZone      string
	placementGroup        string
	instanceType          string
	keyName               string
	subnetID              string
	privateIP             string
	securityGroups        string
	vpcSecurityGroupIDs   string
	publicDNS             string
	publicIP              string
	privateDNS            string
	ebsOptimized          bool
	disableAPITermination bool
	monitoring            bool
	iamInstanceProfile    string
	tenancy               string
	tags                  string
	instanceDescription   *ec2.Instance
}

func (i *instance) describeInstance() {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String("i-02a3b23acf213dbf9"),
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	reservations := result.Reservations[0]
	i.instanceDescription = reservations.Instances[0]
}

func (i *instance) IsInstanceInVPC() bool {
	networkInterfaces := i.instanceDescription.NetworkInterfaces
	for _, nwInterface := range networkInterfaces {
		return nwInterface.VpcId != nil

	}
	return false
}
