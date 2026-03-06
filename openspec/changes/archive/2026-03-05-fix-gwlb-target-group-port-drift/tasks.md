# Implementation Tasks

## 1. Schema Update
- [x] 1.1 Add `Computed: true` to the `port` field schema definition in `resource_tc_gwlb_target_group.go`
- [x] 1.2 Verify the field remains `Optional: true` (not Required)

## 2. Testing
- [x] 2.1 Test resource creation without specifying `port` field
- [x] 2.2 Verify no drift is detected when `port` is not specified
- [x] 2.3 Test resource creation with explicit `port` value
- [x] 2.4 Verify explicit `port` value is respected and no drift occurs

## 3. Documentation
- [x] 3.1 Add changelog entry for the fix (`.changelog/3841.txt`)
