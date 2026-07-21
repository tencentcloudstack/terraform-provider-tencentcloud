## ADDED Requirements

### Requirement: OSType parameter for CLS machine group resource
The `tencentcloud_cls_machine_group` resource SHALL support an `ostype` parameter of type `TypeInt` that specifies the operating system type of the machine group. The valid values are 0 (Linux) and 1 (Windows). This parameter is Optional.

#### Scenario: Create machine group with OSType
- **WHEN** a user creates a `tencentcloud_cls_machine_group` resource with `ostype = 1` (Windows)
- **THEN** the `CreateMachineGroup` API is called with `OSType = 1`
- **AND** the resource is created successfully

#### Scenario: Create machine group without OSType
- **WHEN** a user creates a `tencentcloud_cls_machine_group` resource without specifying `ostype`
- **THEN** the `CreateMachineGroup` API is called without setting `OSType`
- **AND** the cloud API uses its default value (Linux)
- **AND** the resource is created successfully

#### Scenario: Read machine group with OSType
- **WHEN** Terraform reads an existing `tencentcloud_cls_machine_group` resource that has `OSType` set in the cloud API response
- **THEN** the `ostype` value is set in the Terraform state

#### Scenario: Read machine group without OSType (legacy)
- **WHEN** Terraform reads an existing `tencentcloud_cls_machine_group` resource that does NOT have `OSType` in the cloud API response (nil)
- **THEN** the `ostype` field is NOT set in the Terraform state
- **AND** no panic or error occurs

#### Scenario: Change OSType returns error
- **WHEN** a user changes `ostype` from 0 to 1 on an existing resource
- **THEN** Terraform SHALL return an error: "argument `ostype` cannot be changed"
- **AND** the resource is NOT modified

### Requirement: MetaTags parameter for CLS machine group resource
The `tencentcloud_cls_machine_group` resource SHALL support a `meta_tags` parameter of type `TypeMap` that specifies metadata key-value pairs for the machine group.

#### Scenario: Create machine group with MetaTags
- **WHEN** a user creates a `tencentcloud_cls_machine_group` resource with `meta_tags = { "env" = "production" }`
- **THEN** the `CreateMachineGroup` API is called with `MetaTags` containing the key-value pair
- **AND** the resource is created successfully

#### Scenario: Read machine group with MetaTags
- **WHEN** Terraform reads an existing `tencentcloud_cls_machine_group` resource that has `MetaTags` set in the cloud API response
- **THEN** the `meta_tags` map is set in the Terraform state with the correct key-value pairs

#### Scenario: Update MetaTags
- **WHEN** a user changes `meta_tags` on an existing resource
- **THEN** the `ModifyMachineGroup` API is called with the updated `MetaTags`
- **AND** the resource is updated successfully

#### Scenario: Read machine group without MetaTags (legacy)
- **WHEN** Terraform reads an existing `tencentcloud_cls_machine_group` resource that does NOT have `MetaTags` in the cloud API response (nil)
- **THEN** the `meta_tags` field is NOT set in the Terraform state
- **AND** no panic or error occurs