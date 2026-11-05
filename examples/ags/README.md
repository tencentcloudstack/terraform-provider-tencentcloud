# AGS Terraform V1 candidate

This directory is an unpublished integration. Examples pin `0.0.1-local` and
require an isolated local Provider mirror; they do not claim an available
Registry release. The tested Terraform version is 1.17.0. The Provider uses
public Cloud AGS SDK v1.3.202 and the Provider's common SDK v1.3.206 and API version 2025-09-20.

`code`, `browser` and `custom` call `modules/sandbox-environment`. `existing`
performs a Data Source query. `combined` demonstrates existing VPC, COS, CFS,
CAM and CLS references. The module creates only one Tool; callers own every
dependency and all runtime instances. Credentials belong in the Provider's
supported environment/profile/role configuration, never in these files.

The custom example takes an explicit existing `role_arn` for registry/data
access. The observed personal-registry API requires it even for Tool creation;
the module does not create a role or grant permissions. Use verified image
startup/default settings for runtime acceptance; an image visible in a list
only supplies candidate metadata.

The candidate supports only browser, code-interpreter and custom. A successful
Tool create does not prove runtime usability. Image compatibility, regional
CPU/memory support, registry access and data-access role permissions require
separate verification. No dependency is created or granted access implicitly.

## Local candidate verification

From the repository root, source `examples/ags/e2e/env.sh`. The supplied
`e2e/install_candidate.py` and `e2e/validate_examples.py` accept explicit paths
to the built Provider and Terraform executables. The former writes a local
mirror and CLI config, the latter copies the examples into a fresh private
work directory, initializes/validates them, checks the actual mux schema and
plans one code Tool without applying it. Keep state, saved plans and raw API
observations outside the repository.

The E2E README describes executable evidence entry points and current coverage.
Ordinary Go tests and example validation do not create cloud resources. Real
cloud runs require an explicit opt-in and create only ledger-owned resources.

## Import an existing Tool

1. Select the owning account and region. Ensure the Tool is not in another
   Terraform state. Read its full configuration with the public Cloud API.
2. Write equivalent HCL, including supported network, custom configuration,
   mounts, role and logs. Keep it out of a Module during initial review if this
   makes the address easier to inspect. Unsupported or incomplete settings
   cause import to fail instead of being silently dropped.
3. Run `terraform import tencentcloud_ags_sandbox_tool.code sdt-xxxxxxxx` for
   the direct resource example, or use the actual module resource address
   `module.environment.tencentcloud_ags_sandbox_tool.this`.
4. Run `terraform plan -detailed-exitcode` twice. Both must exit 0 before the
   object is considered safely adopted. A replacement plan means the HCL and
   immutable remote settings differ. Operation `timeouts` are local settings
   and cannot be discovered from AGS; adding a custom timeouts block after
   import requires a local state update and is not a remote configuration
   drift or an API mutation.

## Recover a failed create or update

- A create that received an ID keeps that identity when readiness fails.
  Read the exact ID in the same account/region and compare actual status and
  **every managed field** against HCL. An ACTIVE status alone is insufficient.
  Only after that verification may a tainted but healthy resource be recovered
  using `terraform untaint <actual-resource-address>`. Plan twice afterward.
- If no ID was returned, reconcile the run-owned name and API audit evidence
  before another create. A finite-lived ClientToken is not permanent dedupe.
  Do not automatically adopt a same-name shared object. Import the proven
  run-owned identity after removing any conflicting state ownership.
- Updates wait for the requested fields, including empty resets. A failure
  preserves ID and trustworthy readback, including partial success. Inspect
  the remaining plan and apply only after the underlying cause is understood;
  the Provider does not roll back, restart instances or recreate automatically.
- Read failures, denied access and ambiguous empty lists retain state. FAILED
  Tools remain readable for diagnostics and can be planned for destroy.

## Delete and current release gates

Release running **and paused** instances through their runtime owner before
destroying the Tool. The Provider never starts, stops, resumes or obtains an
instance token. Fixed-field changes replace the Tool; same-name
create-before-destroy can conflict. Prefer a new Tool name, application ID
handoff, old instance release, then old Tool deletion for a gradual migration.

In the observed environment an empty description is accepted but not cleared;
the Provider reports nonconvergence. This is recorded as a server defect and
accepted as nonblocking for the current local validation, not as a passing reset.

For deletion, a successful Delete with a valid public request receipt starts a
bounded confirmation on the same client and exact ID. An explicit NotFound or
a complete successful empty point query completes that deletion. A lost or
missing acknowledgement, denied/incomplete read, cancellation or timeout keeps
state. For an already managed Resource, refresh now treats a complete successful
empty exact-ID query as absence, allowing the next plan/apply to recreate it.
This follows the existing VPC/CFS/TCR convention and cannot distinguish a
successful permission-filtered empty response; use the intended account/region.
Explicit API/network errors and incomplete responses retain state. A provisional
CREATING identity and import of a missing Tool remain conservative. Delete
receipts never carry over into later operations or other resources. Independent owner cleanup is a
separate recovery measure and does not authorize generic state removal.
Full V1 acceptance remains incomplete.

User tags can clear; the immutable `qcs:tag:createdBy` system tag remains visible
only in the effective actual configuration. Other unexpected keys remain drift.
Environment variables and configuration are stored in state: use only
non-sensitive values and protect state/plan artifacts.
