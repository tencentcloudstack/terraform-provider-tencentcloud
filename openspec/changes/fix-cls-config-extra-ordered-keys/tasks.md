## 1. Resource Schema and State

- [x] 1.1 Change `extract_rule.keys` from `schema.TypeSet` to `schema.TypeList`.
- [x] 1.2 Preserve key order in Create and Update request expansion.
- [x] 1.3 Set schema version 1 and register the version 0 legacy schema.
- [x] 1.4 Ensure legacy schema construction does not mutate the current schema.

## 2. Verification

- [x] 2.1 Add a unit test for the ordered list schema.
- [x] 2.2 Exercise JSON and legacy flatmap migration through
  `GRPCProviderServer.UpgradeResourceState` and decode the resulting current
  schema state.
- [x] 2.3 Extend acceptance coverage for create order, reversed update order,
  post-refresh convergence, and import state.
- [ ] 2.4 Run the acceptance test in an isolated TencentCloud test account.
- [x] 2.5 Rebase onto the current upstream `master` and rerun build, format,
  lint, mux, and targeted tests.

## 3. Documentation and Contribution

- [x] 3.1 Update the resource example source and generated website docs to
  document `keys` as an ordered list.
- [x] 3.2 Record the schema, migration, compatibility, and rollback contract in
  this OpenSpec change.
- [x] 3.3 Create the upstream issue using the bug report template.
- [x] 3.4 Add `.changelog/4451.txt` after the upstream PR number is assigned.
- [x] 3.5 Open the upstream PR with the issue reference and verification status.
