---
layout: "aws"
page_title: "AWS: aws_vpc_endpoint"
sidebar_current: "docs-aws-resource-vpc-endpoint"
description: |-
  Provides a VPC Endpoint resource.
---

# aws_vpc_endpoint

Provides a VPC Endpoint resource.

~> **NOTE on VPC Endpoints and VPC Endpoint Associations:** Terraform provides both standalone VPC Endpoint Associations for
[Route Tables](vpc_endpoint_route_table_association.html) - (an association between a VPC endpoint and a single `route_table_id`) and
[Subnets](vpc_endpoint_subnet_association.html) - (an association between a VPC endpoint and a single `subnet_id`) and
a VPC Endpoint resource with `route_table_ids` and `subnet_ids` attributes.
Do not use the same resource ID in both a VPC Endpoint resource and a VPC Endpoint Association resource.
Doing so will cause a conflict of associations and will overwrite the association.

## Example Usage

Basic usage:

```hcl
resource "aws_vpc_endpoint" "s3" {
  vpc_id       = "${aws_vpc.main.id}"
  service_name = "com.amazonaws.us-west-2.s3"
}
```

Interface type usage:

```hcl
resource "aws_vpc_endpoint" "ec2" {
  vpc_id            = "${aws_vpc.main.id}"
  service_name      = "com.amazonaws.us-west-2.ec2"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    "${aws_security_group.sg1.id}"
  ]

  private_dns_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The ID of the VPC in which the endpoint will be used.
* `vpc_endpoint_type` - (Optional) The VPC endpoint type, `Gateway` or `Interface`. Defaults to `Gateway`.
* `service_name` - (Required) The service name, in the form `com.amazonaws.region.service` for AWS services.
* `auto_accept` - (Optional) Accept the VPC endpoint (the VPC endpoint and service need to be in the same AWS account).
* `policy` - (Optional) A policy to attach to the endpoint that controls access to the service. Applicable for endpoints of type `Gateway`.
Defaults to full access.
* `route_table_ids` - (Optional) One or more route table IDs. Applicable for endpoints of type `Gateway`.
* `subnet_ids` - (Optional) The ID of one or more subnets in which to create a network interface for the endpoint. Applicable for endpoints of type `Interface`.
* `security_group_ids` - (Optional) The ID of one or more security groups to associate with the network interface. Required for endpoints of type `Interface`.
* `private_dns_enabled` - (Optional) Whether or not to associate a private hosted zone with the specified VPC. Applicable for endpoints of type `Interface`.
Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC endpoint.
* `state` - The state of the VPC endpoint.
* `prefix_list_id` - The prefix list ID of the exposed AWS service. Applicable for endpoints of type `Gateway`.
* `cidr_blocks` - The list of CIDR blocks for the exposed AWS service. Applicable for endpoints of type `Gateway`.
* `network_interface_ids` - One or more network interfaces for the VPC Endpoint. Applicable for endpoints of type `Interface`.
* `dns_entry` - The DNS entries for the VPC Endpoint. Applicable for endpoints of type `Interface`. DNS blocks are documented below.

DNS blocks (for `dns_entry`) support the following attributes:

* `dns_name` - The DNS name.
* `hosted_zone_id` - The ID of the private hosted zone.

## Import

VPC Endpoints can be imported using the `vpc endpoint id`, e.g.

```
$ terraform import aws_vpc_endpoint.endpoint1 vpce-3ecf2a57
```
