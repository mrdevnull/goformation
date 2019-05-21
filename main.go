package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

func unmarshal() {

	template, err := goformation.Open("template.yaml")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	function, err := template.GetAWSEC2VPCWithName("myVPC")

	if err != nil {
		log.Fatalf("There was an error processing the function: %s", err)
	}

	log.Printf("Found a %s\n\n", function.AWSCloudFormationType())

	yaml, err := template.YAML()
	if err != nil {
		fmt.Printf("Failed to generate YAML: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(yaml))
	}

}

func marshal() {
	// Create a new CloudFormation template
	template := cloudformation.NewTemplate()

	// Create an Amazon SNS topic, with a unique name based off the current timestamp
	template.Resources["MyTopic"] = &resources.AWSSNSTopic{
		TopicName: "my-topic-" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	// Create a subscription, connected to our topic, that forwards notifications to an email address
	template.Resources["MyTopicSubscription"] = &resources.AWSSNSSubscription{
		TopicArn: cloudformation.Ref("MyTopic"),
		Protocol: "email",
		Endpoint: "some.email@example.com",
	}

	template.Resources["HoneyPotVPC"] = &resources.AWSEC2VPC{
		CidrBlock:          "10.0.0.0/16",
		EnableDnsHostnames: false,
		EnableDnsSupport:   false,
	}

	template.Resources["HoneyPotSubnet"] = &resources.AWSEC2Subnet{
		AssignIpv6AddressOnCreation: false,
		AvailabilityZone:            "",
		CidrBlock:                   "10.0.10.0/24",
		Ipv6CidrBlock:               "",
		MapPublicIpOnLaunch:         false,
		Tags:                        nil,
		VpcId:                       cloudformation.Ref("HoneyPotVPC"),
	}

	template.Resources["BastianSubnet"] = &resources.AWSEC2Subnet{
		AssignIpv6AddressOnCreation: false,
		AvailabilityZone:            "",
		CidrBlock:                   "10.0.20.0/24",
		Ipv6CidrBlock:               "",
		MapPublicIpOnLaunch:         false,
		Tags:                        nil,
		VpcId:                       cloudformation.Ref("HoneyPotVPC"),
	}

	var egress = []resources.AWSEC2SecurityGroup_Egress{
		{
			CidrIp:                     "",
			CidrIpv6:                   "",
			Description:                "",
			DestinationPrefixListId:    "",
			DestinationSecurityGroupId: "",
			FromPort:                   0,
			IpProtocol:                 "",
			ToPort:                     0,
		},
	}
	ingress := []resources.AWSEC2SecurityGroup_Ingress{
		{
			CidrIp:                     "",
			CidrIpv6:                   "",
			Description:                "",
			FromPort:                   0,
			IpProtocol:                 "",
			SourcePrefixListId:         "",
			SourceSecurityGroupId:      "",
			SourceSecurityGroupName:    "",
			SourceSecurityGroupOwnerId: "",
			ToPort:                     0,
		},
	}

	template.Resources["HoneyPotSecurityGroup"] = &resources.AWSEC2SecurityGroup{
		GroupDescription:     "",
		GroupName:            "",
		SecurityGroupEgress:  egress,
		SecurityGroupIngress: ingress,
		Tags:                 nil,
		VpcId:                "",
	}

	template.Resources["BastionSecurityGroup"] = &resources.AWSEC2SecurityGroup{
		GroupDescription:     "",
		GroupName:            "",
		SecurityGroupEgress:  nil,
		SecurityGroupIngress: nil,
		Tags:                 nil,
		VpcId:                cloudformation.Ref("HoneyPotVPC"),
	}

	// Let's see the JSON AWS CloudFormation template
	j, err := template.JSON()
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(j))
	}

	// and also the YAML AWS CloudFormation template
	y, err := template.YAML()
	if err != nil {
		fmt.Printf("Failed to generate YAML: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(y))
	}
}

func main() {

	marshal()

	unmarshal()

}
