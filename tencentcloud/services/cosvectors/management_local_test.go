package cosvectors

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

func cosVectorGraphConfig(bucket, index, effect string, dimension int, metadata string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_vector_bucket" "test" {
  vector_bucket_name = %q
}
resource "tencentcloud_cos_vector_index" "test" {
  vector_bucket_name = tencentcloud_cos_vector_bucket.test.vector_bucket_name
  index_name = %q
  data_type = "float32"
  dimension = %d
  distance_metric = "cosine"
  %s
}
resource "tencentcloud_cos_vector_bucket_policy" "test" {
  vector_bucket_qcs = tencentcloud_cos_vector_bucket.test.vector_bucket_qcs
  policy = jsonencode({
    version = "2.0"
    statement = [{
      effect = %q
      principal = { qcs = ["qcs::cam::uin/100000000001:uin/100000000002"] }
      action = ["name/cos:QueryVectors", "name/cos:GetVectors"]
      resource = ["${tencentcloud_cos_vector_bucket.test.vector_bucket_qcs}/*"]
    }]
  })
}
`, bucket, index, dimension, metadata, effect)
}

func TestCosVectorManagementLocalGraph(t *testing.T) {
	_, mock, factories := cosVectorLocalSetup(t)
	bucket := "tf-graph-1250000000"
	config := cosVectorGraphConfig(bucket, "index-test", "allow", 3, "")
	metadata := `metadata_configuration { non_filterable_metadata_keys = ["text", "source"] }`
	expectedPolicy := strings.Replace(cosVectorTestPolicy("qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/"+bucket, "allow"), `["name/cos:QueryVectors"]`, `["name/cos:QueryVectors","name/cos:GetVectors"]`, 1)
	checkImportedPolicy := func(states []*terraform.InstanceState) error {
		if len(states) != 1 || !cosVectorPoliciesEqual(states[0].Attributes["policy"], expectedPolicy) {
			return fmt.Errorf("imported policy differs from the original authorization")
		}
		return nil
	}
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_index.test", "dimension", "3"),
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_index.test", "id", bucket+"/index-test"),
			resource.TestCheckResourceAttrSet("tencentcloud_cos_vector_bucket_policy.test", "policy"))},
		{Config: config, PlanOnly: true},
		{Config: strings.ReplaceAll(config, `["name/cos:QueryVectors", "name/cos:GetVectors"]`, `["name/cos:GetVectors", "name/cos:QueryVectors"]`), PlanOnly: true},
		{ResourceName: "tencentcloud_cos_vector_index.test", ImportState: true, ImportStateVerify: true, ImportStateId: bucket + "/index-test"},
		{ResourceName: "tencentcloud_cos_vector_index.test", ImportState: true, ImportStateVerify: true, ImportStateId: "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/" + bucket + "/index/index-test"},
		{ResourceName: "tencentcloud_cos_vector_bucket_policy.test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"policy"}, ImportStateCheck: checkImportedPolicy, ImportStateId: "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/" + bucket},
		{ResourceName: "tencentcloud_cos_vector_bucket_policy.test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"policy"}, ImportStateCheck: checkImportedPolicy, ImportStateId: bucket},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 3, "")},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 3, ""), PlanOnly: true},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, metadata)},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, `metadata_configuration { non_filterable_metadata_keys = ["source", "text"] }`), PlanOnly: true},
		{ResourceName: "tencentcloud_cos_vector_index.test", ImportState: true, ImportStateVerify: true, ImportStateId: bucket + "/index-test"},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, ""), PlanOnly: true},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, metadata), PreConfig: func() { mock.mu.Lock(); delete(mock.indexes, bucket+"/index-test"); mock.mu.Unlock() }, PlanOnly: true, ExpectNonEmptyPlan: true},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, metadata)},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, metadata), PreConfig: func() { mock.mu.Lock(); delete(mock.policies, bucket); mock.mu.Unlock() }, PlanOnly: true, ExpectNonEmptyPlan: true},
		{Config: cosVectorGraphConfig(bucket, "index-test", "deny", 4, metadata)},
		{Config: cosVectorGraphConfig("tf-graph-renamed-1250000000", "index-renamed", "deny", 4, metadata)},
		{Config: cosVectorGraphConfig("tf-graph-renamed-1250000000", "index-renamed", "deny", 4, metadata), PlanOnly: true},
	}})
	if mock.count() != 0 {
		t.Fatal("resource graph destroy left buckets")
	}
	if mock.callCount("PutVectorBucketPolicy") != 4 {
		t.Fatalf("expected create, update, drift recreate and new bucket policy, got %d", mock.callCount("PutVectorBucketPolicy"))
	}
}

func TestCosVectorManagementLocalCleanupPending(t *testing.T) {
	_, mock, factories := cosVectorLocalSetup(t)
	config := cosVectorGraphConfig("tf-gc-local-1250000000", "index-test", "allow", 3, "")
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config},
		{Config: config, Destroy: true, PreConfig: func() {
			mock.mu.Lock()
			mock.cleanupPending = 1
			mock.mu.Unlock()
		}},
	}})
	if mock.count() != 0 || mock.callCount("DeleteVectorBucket") != 2 || mock.callCount("ListIndexes") != 1 || mock.callCount("DeleteIndex") != 1 {
		t.Fatal("Terraform destroy did not wait for GC or unexpectedly emptied the bucket")
	}
}

func TestCosVectorManagementLocalInvalidConfig(t *testing.T) {
	_, _, factories := cosVectorLocalSetup(t)
	for _, tc := range []struct{ name, config, error string }{
		{"indexName", cosVectorGraphConfig("tf-invalid-1250000000", "Upper_Index", "allow", 3, ""), "Invalid COS vector index name"},
		{"dimension", cosVectorGraphConfig("tf-invalid-1250000000", "valid-index", "allow", 0, ""), "Invalid COS vector dimension"},
		{"metadata", cosVectorGraphConfig("tf-invalid-1250000000", "valid-index", "allow", 3, `metadata_configuration { non_filterable_metadata_keys = ["invalid-key"] }`), "Invalid COS vector metadata"},
		{"awsPolicy", `resource "tencentcloud_cos_vector_bucket_policy" "test" {
 vector_bucket_qcs = "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-invalid-1250000000"
 policy = jsonencode({ Version = "2012-10-17", Statement = [] })
}`, "Invalid COS vector bucket policy"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{{Config: tc.config, PlanOnly: true, ExpectError: regexp.MustCompile(tc.error)}}})
		})
	}
}

func TestCosVectorManagementLocalForceDestroy(t *testing.T) {
	service, mock, factories := cosVectorLocalSetup(t)
	bucket := "tf-force-local-1250000000"
	config := cosVectorLocalConfig(bucket, "force_destroy = true")
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config, Check: resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "force_destroy", "true")},
		{Config: config, PlanOnly: true, PreConfig: func() {
			for _, name := range []string{"outside-one", "outside-two"} {
				if _, err := service.CreateIndex(context.Background(), &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: name, DataType: "float32", Dimension: 3, DistanceMetric: "cosine"}); err != nil {
					t.Fatal(err)
				}
			}
			mock.mu.Lock()
			mock.listPageSize = 1
			mock.mu.Unlock()
		}},
	}})
	if mock.count() != 0 || mock.callCount("DeleteIndex") != 2 {
		t.Fatal("force destroy did not remove unmanaged indexes")
	}
}
