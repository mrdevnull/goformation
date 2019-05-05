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

	function, error := template.GetAWSEC2VPCWithName("myVPC")

	if error != nil {
		log.Fatalf("There was an error processing the function: %s", err)
	}

	log.Printf("Found a %s\n\n", function.AWSCloudFormationType())

	yaml, error := template.YAML()
	if error != nil {
		fmt.Printf("Failed to generate YAML: %s\n", error)
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
