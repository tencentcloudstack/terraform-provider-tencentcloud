# ccn-instances-accept-attach-order-type Specification

## Purpose
Add the `order_type` optional sub-field to the `instances` block of the `tencentcloud_ccn_instances_accept_attach` resource, allowing users to specify the billing account (`PayByCcnOwner` or `PayByInstanceOwner`) for accepted CCN attachment instances.
## Requirements
### Requirement: Order type on CCN instances accept attach
The `tencentcloud_ccn_instances_accept_attach` resource SHALL support an optional `order_type` parameter (TypeString) as a sub-field of the `instances` block. When specified, the provider SHALL pass the value to the `AcceptAttachCcnInstances` API by setting `CcnInstance.OrderType` for each instance in the request. Valid values are `PayByCcnOwner` (CCN owner pays) and `PayByInstanceOwner` (instance owner pays). When the user does not specify `order_type`, the provider SHALL NOT set `OrderType` on the `CcnInstance` struct, allowing the API to apply its default behavior.

#### Scenario: Create CCN instances accept attach with order_type specified
- **WHEN** a user specifies `order_type = "PayByCcnOwner"` within an `instances` block in the `tencentcloud_ccn_instances_accept_attach` resource configuration
- **THEN** the provider SHALL set `CcnInstance.OrderType = "PayByCcnOwner"` for that instance in the `AcceptAttachCcnInstances` API request

#### Scenario: Create CCN instances accept attach without order_type
- **WHEN** a user does NOT specify `order_type` within an `instances` block in the `tencentcloud_ccn_instances_accept_attach` resource configuration
- **THEN** the provider SHALL NOT set `OrderType` on the `CcnInstance` struct for that instance in the `AcceptAttachCcnInstances` API request

#### Scenario: Order type is immutable after creation
- **WHEN** a user changes `order_type` within an `instances` block after the resource has been created
- **THEN** the provider SHALL trigger recreation of the resource (ForceNew behavior inherited from the `instances` block)

#### Scenario: Order type with empty value is ignored
- **WHEN** a user specifies `order_type = ""` within an `instances` block in the `tencentcloud_ccn_instances_accept_attach` resource configuration
- **THEN** the provider SHALL NOT set `OrderType` on the `CcnInstance` struct for that instance in the `AcceptAttachCcnInstances` API request
