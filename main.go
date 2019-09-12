package main

import (
	"fmt"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
	"github.com/davecgh/go-spew/spew"
	"log"
)

func unmarshal() {

	template, err := goformation.Open("template.json")
	spew.Dump(template)
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	function, err := template.GetAWSEC2VPCWithName("myVPC")

	if err != nil {
		log.Fatalf("There was an error processing the function: %s", err)
	}

	log.Printf("Found a %s\n\n", function.AWSCloudFormationType())

	json, err := template.JSON()
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(json))
	}

}

func marshal() {
	template := cloudformation.NewTemplate()

	template.Resources["HoneyPotVPC"] = &resources.AWSEC2VPC{
		CidrBlock:          "10.0.0.0/16",
		EnableDnsHostnames: false,
		EnableDnsSupport:   false,
	}

	spew.Dump(template)

	template.Resources["LoadBalancerSubnet"] = &resources.AWSEC2Subnet{
		AssignIpv6AddressOnCreation: false,
		AvailabilityZone:            "",
		CidrBlock:                   "10.0.10.0/24",
		Ipv6CidrBlock:               "",
		MapPublicIpOnLaunch:         false,
		Tags:                        nil,
		VpcId:                       cloudformation.Ref("HoneyPotVPC"),
	}

	template.Resources["HoneyPotSubnet"] = &resources.AWSEC2Subnet{
		AssignIpv6AddressOnCreation: false,
		AvailabilityZone:            "",
		CidrBlock:                   "10.0.20.0/24",
		Ipv6CidrBlock:               "",
		MapPublicIpOnLaunch:         false,
		Tags:                        nil,
		VpcId:                       cloudformation.Ref("HoneyPotVPC"),
	}

	template.Resources["BastianSubnet"] = &resources.AWSEC2Subnet{
		AssignIpv6AddressOnCreation: false,
		AvailabilityZone:            "",
		CidrBlock:                   "10.0.30.0/24",
		Ipv6CidrBlock:               "",
		MapPublicIpOnLaunch:         false,
		Tags:                        nil,
		VpcId:                       cloudformation.Ref("HoneyPotVPC"),
	}

	loadBalancerEgress := []resources.AWSEC2SecurityGroup_Egress{
		{
			CidrIp:                     "0.0.0.0/0",
			CidrIpv6:                   "",
			Description:                "",
			DestinationPrefixListId:    "",
			DestinationSecurityGroupId: "",
			FromPort:                   0,
			IpProtocol:                 "",
			ToPort:                     0,
		},
	}
	loafBalancerIngress := []resources.AWSEC2SecurityGroup_Ingress{
		{
			CidrIp:                     "0.0.0.0/0",
			CidrIpv6:                   "",
			Description:                "",
			FromPort:                   22,
			IpProtocol:                 "",
			SourcePrefixListId:         "",
			SourceSecurityGroupId:      "",
			SourceSecurityGroupName:    "",
			SourceSecurityGroupOwnerId: "",
			ToPort:                     2222,
		},
	}

	template.Resources["LoadBalancerSecurityGroup"] = &resources.AWSEC2SecurityGroup{
		GroupDescription:     "",
		GroupName:            "",
		SecurityGroupEgress:  loadBalancerEgress,
		SecurityGroupIngress: loafBalancerIngress,
		Tags:                 nil,
		VpcId:                cloudformation.Ref("HoneyPotVPC"),
	}

	honeyPotEgress := []resources.AWSEC2SecurityGroup_Egress{
		{
			CidrIp:                     "0.0.0.0/0",
			CidrIpv6:                   "",
			Description:                "",
			DestinationPrefixListId:    "",
			DestinationSecurityGroupId: "",
			FromPort:                   0,
			IpProtocol:                 "",
			ToPort:                     0,
		},
	}
	honeyPotIngress := []resources.AWSEC2SecurityGroup_Ingress{
		{
			CidrIp:                     "",
			CidrIpv6:                   "",
			Description:                "",
			FromPort:                   2222,
			IpProtocol:                 "",
			SourcePrefixListId:         "",
			SourceSecurityGroupId:      cloudformation.Ref("LoadBalancerSecurityGroup"),
			SourceSecurityGroupName:    cloudformation.Ref("LoadBalancerSecurityGroup"),
			SourceSecurityGroupOwnerId: cloudformation.Ref("LoadBalancerSecurityGroup"),
			ToPort:                     2222,
		},
		{
			CidrIp:                     "",
			CidrIpv6:                   "",
			Description:                "",
			FromPort:                   22,
			IpProtocol:                 "",
			SourcePrefixListId:         "",
			SourceSecurityGroupId:      cloudformation.Ref("BastionSecurityGroup"),
			SourceSecurityGroupName:    cloudformation.Ref("BastionSecurityGroup"),
			SourceSecurityGroupOwnerId: cloudformation.Ref("BastionSecurityGroup"),
			ToPort:                     22,
		},
	}

	template.Resources["HoneyPotSecurityGroup"] = &resources.AWSEC2SecurityGroup{
		GroupDescription:     "",
		GroupName:            "",
		SecurityGroupEgress:  honeyPotEgress,
		SecurityGroupIngress: honeyPotIngress,
		Tags:                 nil,
		VpcId:                cloudformation.Ref("HoneyPotVPC"),
	}

	bastionEgress := []resources.AWSEC2SecurityGroup_Egress{
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
	bastionIngress := []resources.AWSEC2SecurityGroup_Ingress{
		{
			CidrIp:                     "",
			CidrIpv6:                   "",
			Description:                "",
			FromPort:                   22,
			IpProtocol:                 "",
			SourcePrefixListId:         "",
			SourceSecurityGroupId:      "",
			SourceSecurityGroupName:    "",
			SourceSecurityGroupOwnerId: "",
			ToPort:                     22,
		},
	}

	template.Resources["BastionSecurityGroup"] = &resources.AWSEC2SecurityGroup{
		GroupDescription:     "",
		GroupName:            "",
		SecurityGroupEgress:  bastionEgress,
		SecurityGroupIngress: bastionIngress,
		Tags:                 nil,
		VpcId:                cloudformation.Ref("HoneyPotVPC"),
	}

	template.Resources["InternetGateway"] = &resources.AWSEC2InternetGateway{
		Tags: nil,
	}

	template.Resources["InternetGatewayAttachment"] = &resources.AWSEC2VPCGatewayAttachment{
		InternetGatewayId: cloudformation.Ref("InternetGateway"),
		VpcId:             cloudformation.Ref("HoneyPotVPC"),
		VpnGatewayId:      "",
	}

	var eip = resources.AWSEC2EIP{
		Domain:         "vpc",
		InstanceId:     "",
		PublicIpv4Pool: "",
	}

	eip.SetDependsOn([]string{"InternetGatewayAttachment"})

	template.Resources["NATEIP"] = &eip

	template.Resources["NATGateway"] = &resources.AWSEC2NatGateway{
		AllocationId: cloudformation.GetAtt("NATEIP", "AllocationId"),
		SubnetId:     cloudformation.Ref("HoneyPotSubnet"),
		Tags:         nil,
	}

	template.Resources["RouteTable"] = &resources.AWSEC2RouteTable{
		Tags:  nil,
		VpcId: cloudformation.Ref("HoneyPotVPC"),
	}

	template.Resources["Route"] = &resources.AWSEC2Route{
		DestinationCidrBlock:        "0.0.0.0/0",
		DestinationIpv6CidrBlock:    "",
		EgressOnlyInternetGatewayId: "",
		GatewayId:                   "",
		InstanceId:                  "",
		NatGatewayId:                cloudformation.Ref("NATGateway"),
		NetworkInterfaceId:          "",
		RouteTableId:                cloudformation.Ref("RouteTable"),
		VpcPeeringConnectionId:      cloudformation.ImportValue(cloudformation.Sub(("${NetworkStackNameParameter}-SecurityGroupID"))),
	}

	template.Outputs["BastianSubnet"] = map[string]interface{}{
		"Description": "Bastion Subnet",
		"Value": cloudformation.Ref("BastionSubnet"),
		"Export": map[string]interface{}{
			"Name": cloudformation.Sub("${AWS::StackName}-SubnetID"),
		},
	}

	y, err := template.YAML()
	if err != nil {
		fmt.Printf("Failed to generate YAML: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(y))
	}

	y, err = template.JSON()
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(y))
	}
}

func main() {

	marshal()

	//unmarshal()

}
