# vpc-private-nat-gateway-translation-nat-rule-typed-lists Specification

## ADDED Requirements

### Requirement: Schema MUST expose three typed `TypeList` fields per bucket

`tencentcloud_vpc_private_nat_gateway_translation_nat_rule` SHALL declare three new `TypeList` schema fields, each with `Optional: true, Computed: true`, with bucket-tailored inner schema:

* `local_network_layer_rules`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional)
* `local_transport_layer_rules`: `translation_ip` (Required), `description` (Optional). MUST NOT expose `original_ip`.
* `peer_network_layer_rules`: `translation_ip` (Required), `original_ip` (Required), `description` (Optional)

Inner `translation_ip` and `original_ip` MUST NOT carry `ForceNew`.

#### Scenario: `local_transport_layer_rules` does not expose `original_ip`

- **WHEN** a developer inspects the inner schema of `local_transport_layer_rules`
- **THEN** the inner attribute keys are exactly `translation_ip` and `description`; `original_ip` is absent.

### Requirement: `translation_nat_rules` MUST remain Deprecated

The legacy `translation_nat_rules` field SHALL retain its existing `TypeSet` shape and inner attributes unchanged, with a `Deprecated:` annotation pointing to the three new fields. The legacy field continues to plan and apply as before.

#### Scenario: Legacy field still works

- **GIVEN** an HCL config that uses only `translation_nat_rules`
- **WHEN** the user runs `terraform plan` and `terraform apply`
- **THEN** the resource behaves exactly as before; only a deprecation warning is emitted.

### Requirement: Pairwise `ConflictsWith` between legacy and each new field

Each of the four nested-block fields SHALL declare `ConflictsWith` only against fields on the other side of the legacy/new boundary:

* `translation_nat_rules` ConflictsWith `[local_network_layer_rules, local_transport_layer_rules, peer_network_layer_rules]`
* Each of the three new fields ConflictsWith `[translation_nat_rules]` only

The three new fields MUST NOT be mutually exclusive among themselves.

#### Scenario: Mixing legacy and any new field fails at plan time

- **GIVEN** HCL that sets both `translation_nat_rules` and `local_network_layer_rules`
- **WHEN** the user runs `terraform plan`
- **THEN** Terraform aborts with a `ConflictsWith` validation error.

#### Scenario: Mixing two of the new fields is allowed

- **GIVEN** HCL that sets both `local_network_layer_rules` and `peer_network_layer_rules`
- **WHEN** the user runs `terraform plan`
- **THEN** plan succeeds.

### Requirement: Read MUST preserve the API's authoritative order

For each typed list field, Read SHALL emit rules in the order returned by `DescribePrivateNatGatewayTranslationNatRule` and MUST NOT reorder, deduplicate, or hide any rule of the matching `(direction, type)` bucket. Specifically, Read SHALL:

1. Walk `respData` in API order.
2. For each item whose `(TranslationDirection, TranslationType)` matches the bucket, append a flattened map to that bucket's list.
3. Call `d.Set(field, list)` for each typed list independently.

The provider MUST NOT use the user's HCL order to reshuffle Read output. The user is responsible for honoring the append-only contract: new rules are added at the tail, existing rules MUST NOT be inserted in the middle, reordered, or removed from the middle.

#### Scenario: Read returns rules in API order regardless of HCL order

- **GIVEN** HCL declares four `local_transport_layer_rules` items in order `[A, B, C, D]` (by `translation_ip`), and the API returns them in order `[B, D, C, A]`
- **WHEN** Read runs
- **THEN** state stores `local_transport_layer_rules` as `[B, D, C, A]` (API authoritative order). The next `terraform plan` will report a diff against the user's HCL order until the user fixes HCL to match the API order.

#### Scenario: External rule additions appear at the API-order position

- **GIVEN** the user manages two rules via Terraform and a third party adds an extra rule out-of-band that the API places at index 1
- **WHEN** Read runs
- **THEN** state stores three rules at indices `[0, 1, 2]` matching the API's response order; the next `terraform plan` highlights the third-party addition.

### Requirement: Create MUST aggregate inputs in canonical bucket order

When the user provides any of the new fields, Create SHALL concatenate rules in this exact canonical order before sending: `local_network_layer_rules` → `local_transport_layer_rules` → `peer_network_layer_rules`. Each item MUST be tagged with the bucket's `(direction, type)`. The aggregated rules SHALL be batched by `MAX_CREATE_RULES_LEN = 20` and each batch wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Mixed typed-list create produces correctly tagged inputs

- **GIVEN** HCL with one `local_network_layer_rules`, one `local_transport_layer_rules`, one `peer_network_layer_rules`
- **WHEN** Create runs
- **THEN** the request payload contains three rules, in this order: `LOCAL/NETWORK_LAYER`, `LOCAL/TRANSPORT_LAYER`, `PEER/NETWORK_LAYER`, each with the bucket-correct fields.

### Requirement: Update routes by positional (index-aligned) diff per typed list

For each of the three new fields, Update SHALL:

1. Skip the branch when `d.HasChange(field)` is false.
2. Cast `oldI, newI := d.GetChange(field)` to `[]interface{}` and walk both as `[]map[string]interface{}`.
3. Diff by index position (the typed lists are ordered lists; identity within the list is its slot index, not its IP value):
   - For `i` in `[0, min(len(old), len(new)))`: when `equalRuleMap(old[i], new[i], includeOriginalIp)` is false, emit a MODIFY pair `(old[i] → new[i])`. The `buildModifyDiff` helper SHALL populate both `OldOriginalIp/OriginalIp` and `OldTranslationIp/TranslationIp` so the SDK can rename slots whose key IP changes.
   - When `len(new) > len(old)`: items at indices `[len(old), len(new))` of `new` are emitted as CREATE in HCL order.
   - When `len(old) > len(new)`: items at indices `[len(new), len(old))` of `old` are emitted as DELETE in original order.
4. Issue calls in DELETE → MODIFY → CREATE order.

The legacy `translation_nat_rules` Update branch SHALL remain unchanged (`Set.Difference` add/remove).

#### Scenario: Editing only `description` produces a single MODIFY

- **GIVEN** state has one `local_network_layer_rules` rule with `description = "old"`
- **WHEN** the user changes `description` to `"new"`
- **THEN** Update issues exactly one `ModifyPrivateNatGatewayTranslationNatRule` call. No Create or Delete.

#### Scenario: Editing the non-key IP triggers MODIFY (three-layer)

- **GIVEN** state has one three-layer rule with `original_ip=1.1.1.1, translation_ip=2.2.2.2`
- **WHEN** the user changes `translation_ip` to `9.9.9.9`
- **THEN** Update issues exactly one `ModifyPrivateNatGatewayTranslationNatRule` with `OldTranslationIp=2.2.2.2, TranslationIp=9.9.9.9, OldOriginalIp=1.1.1.1, OriginalIp=1.1.1.1`.

#### Scenario: Editing the key IP at the same slot triggers MODIFY (three-layer rename)

- **GIVEN** state has one three-layer rule at index 0 with `original_ip=1.1.1.1`
- **WHEN** the user changes `original_ip` to `8.8.8.8` at the same slot
- **THEN** Update issues one `ModifyPrivateNatGatewayTranslationNatRule` with `OldOriginalIp=1.1.1.1, OriginalIp=8.8.8.8`. No Delete or Create.

#### Scenario: Appending a rule at the tail produces a single CREATE

- **GIVEN** state has rules `[A, B, C]`
- **WHEN** the user appends `D` to the tail of HCL, leaving the first three unchanged
- **THEN** Update issues exactly one `CreatePrivateNatGatewayTranslationNatRule` request containing `D`. No Modify or Delete.

### Requirement: MODIFY MUST issue one request per rule

Per the SDK constraint "only supports modifying a single translation rule", every modify pair SHALL produce its own `ModifyPrivateNatGatewayTranslationNatRule` request, each wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Two modify pairs produce two SDK calls

- **GIVEN** the user changes `description` on two different rules in the same plan
- **WHEN** Update runs
- **THEN** exactly two `ModifyPrivateNatGatewayTranslationNatRule` requests are issued.

### Requirement: Documentation MUST steer users to the typed lists

The resource markdown SHALL include two HCL example blocks, in this order:

1. **Recommended typed list usage** (uses `local_network_layer_rules`, `local_transport_layer_rules`, `peer_network_layer_rules`).
2. **Deprecated `translation_nat_rules` usage** (preceded by a `# DEPRECATED` note).

`make doc` regenerates `website/docs/r/vpc_private_nat_gateway_translation_nat_rule.html.markdown`. Hand-editing the website file is forbidden.

#### Scenario: Generated website doc shows the new fields and the deprecated marker

- **WHEN** `make doc` runs
- **THEN** the generated doc shows `**Deprecated**` only on `translation_nat_rules` and lists `local_network_layer_rules`, `local_transport_layer_rules`, `peer_network_layer_rules` as separate arguments without ForceNew on the inner IP fields.
