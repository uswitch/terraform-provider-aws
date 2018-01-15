package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAwsVpcEndpointSubnetAssociation_basic(t *testing.T) {
	var vpce ec2.VpcEndpoint

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcEndpointSubnetAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcEndpointSubnetAssociationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointSubnetAssociationExists(
						"aws_vpc_endpoint_subnet_association.a", &vpce),
				),
			},
		},
	})
}

func testAccCheckVpcEndpointSubnetAssociationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ec2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpc_endpoint_subnet_association" {
			continue
		}

		// Try to find the resource
		resp, err := conn.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
			VpcEndpointIds: aws.StringSlice([]string{rs.Primary.Attributes["vpc_endpoint_id"]}),
		})
		if err != nil {
			// Verify the error is what we want
			ec2err, ok := err.(awserr.Error)
			if !ok {
				return err
			}
			if ec2err.Code() != "InvalidVpcEndpointId.NotFound" {
				return err
			}
			return nil
		}

		vpce := resp.VpcEndpoints[0]
		if len(vpce.SubnetIds) > 0 {
			return fmt.Errorf(
				"Vpc endpoint %s has subnets", *vpce.VpcEndpointId)
		}
	}

	return nil
}

func testAccCheckVpcEndpointSubnetAssociationExists(n string, vpce *ec2.VpcEndpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		resp, err := conn.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
			VpcEndpointIds: aws.StringSlice([]string{rs.Primary.Attributes["vpc_endpoint_id"]}),
		})
		if err != nil {
			return err
		}
		if len(resp.VpcEndpoints) == 0 {
			return fmt.Errorf("Vpc endpoint not found")
		}

		*vpce = *resp.VpcEndpoints[0]

		if len(vpce.SubnetIds) == 0 {
			return fmt.Errorf("no subnet associations")
		}

		for _, id := range vpce.SubnetIds {
			if aws.StringValue(id) == rs.Primary.Attributes["subnet_id"] {
				return nil
			}
		}

		return fmt.Errorf("subnet association not found")
	}
}

const testAccVpcEndpointSubnetAssociationConfig_basic = `
provider "aws" {
  region = "us-west-2"
}

resource "aws_vpc" "foo" {
  cidr_block = "10.0.0.0/16"
}

data "aws_security_group" "default" {
  vpc_id = "${aws_vpc.foo.id}"
  name = "default"
}

resource "aws_vpc_endpoint" "ec2" {
  vpc_id = "${aws_vpc.foo.id}"
  vpc_endpoint_type = "Interface"
  service_name = "com.amazonaws.us-west-2.ec2"
  security_group_ids = ["${data.aws_security_group.default.id}"]
  private_dns_enabled = false
}

resource "aws_subnet" "sn" {
  vpc_id = "${aws_vpc.foo.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.0.0/17"
}

resource "aws_vpc_endpoint_subnet_association" "a" {
  vpc_endpoint_id = "${aws_vpc_endpoint.ec2.id}"
  subnet_id  = "${aws_subnet.sn.id}"
}
`
