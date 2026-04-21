# teo-check-cname-status Specification

## Purpose
TBD - created by archiving change add-teo-check-cname-status-operation. Update Purpose after archive.
## Requirements
### Requirement: Check CNAME status for TEO domains
The system SHALL provide a Terraform operation resource that allows users to check the CNAME configuration status for one or more TEO domains.

#### Scenario: Check CNAME status for multiple domains
- **WHEN** user provides valid `zone_id` and a list of `record_names`
- **THEN** system SHALL call the TEO CheckCnameStatus API with the provided parameters
- **THEN** system SHALL return a list of CNAME status information for each domain

#### Scenario: Check CNAME status for single domain
- **WHEN** user provides valid `zone_id` and a single `record_name` in the `record_names` list
- **THEN** system SHALL call the TEO CheckCnameStatus API
- **THEN** system SHALL return a list containing a single CNAME status entry

#### Scenario: Return CNAME status fields
- **WHEN** system receives a successful response from the TEO API
- **THEN** each CNAME status entry SHALL include:
  - `record_name`: The domain name being checked
  - `cname`: The CNAME address (may be null)
  - `status`: The CNAME status (active/moved)

#### Scenario: Handle empty record_names
- **WHEN** user provides valid `zone_id` but empty `record_names` list
- **THEN** system SHALL return an empty list of CNAME status entries

#### Scenario: Handle API errors
- **WHEN** the TEO API returns an error
- **THEN** system SHALL propagate the error to the user

#### Scenario: Validate required parameters
- **WHEN** user omits the required `zone_id` parameter
- **THEN** system SHALL return a validation error indicating that `zone_id` is required

#### Scenario: Validate required parameters
- **WHEN** user omits the required `record_names` parameter
- **THEN** system SHALL return a validation error indicating that `record_names` is required

