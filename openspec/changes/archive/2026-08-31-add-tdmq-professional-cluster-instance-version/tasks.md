## 1. Schema Definition

- [x] 1.1 Add `instance_version` field to resource schema in `resource_tc_tdmq_professional_cluster.go` as `Optional: true, Computed: true, Type: schema.TypeString`

## 2. Create Method

- [x] 2.1 Set `request.InstanceVersion` in `resourceTencentCloudTdmqProfessionalClusterCreate` when `instance_version` is provided

## 3. Read Method

- [x] 3.1 Read `professionalCluster.InstanceVersion` and set `instance_version` in `resourceTencentCloudTdmqProfessionalClusterRead` (with nil check)

## 4. Documentation

- [x] 4.1 Update `resource_tc_tdmq_professional_cluster.md` to include `instance_version` in example usage

## 5. Unit Tests

- [x] 5.1 Add unit test cases for `instance_version` parameter in `resource_tc_tdmq_professional_cluster_test.go`