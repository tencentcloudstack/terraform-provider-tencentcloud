package cosvectors

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

func (m *cosVectorBucketMock) hasIndexes(bucket string) bool {
	for _, index := range m.indexes {
		if index.VectorBucketName == bucket {
			return true
		}
	}
	return false
}

func (m *cosVectorBucketMock) serveVectorManagement(w http.ResponseWriter, r *http.Request, action string) bool {
	switch action {
	case "CreateIndex", "GetIndex", "DeleteIndex", "ListIndexes", "PutVectorBucketPolicy", "GetVectorBucketPolicy", "DeleteVectorBucketPolicy":
	default:
		return false
	}
	if m.indexes == nil {
		m.indexes = map[string]*cosSDK.IndexInfo{}
	}
	if m.policies == nil {
		m.policies = map[string]string{}
	}
	var req struct {
		cosSDK.CreateIndexOptions
		Policy    string `json:"policy"`
		NextToken string `json:"nextToken"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req) != nil || validateCosVectorBucketName(req.VectorBucketName) != nil {
		cosVectorMockError(w, 400)
		return true
	}
	m.events = append(m.events, action+":"+req.VectorBucketName+"/"+req.IndexName)
	if m.buckets[req.VectorBucketName] == nil {
		cosVectorMockError(w, 404)
		return true
	}
	key := req.VectorBucketName + "/" + req.IndexName
	qcs := m.buckets[req.VectorBucketName].VectorBucketQcs + "/index/" + req.IndexName
	switch action {
	case "CreateIndex":
		if m.indexes[key] != nil {
			cosVectorMockError(w, 409)
			return true
		}
		if validateCosVectorIndexOptions(&req.CreateIndexOptions) != nil {
			cosVectorMockError(w, 400)
			return true
		}
		m.indexes[key] = &cosSDK.IndexInfo{IndexName: req.IndexName, VectorBucketName: req.VectorBucketName, IndexQcs: qcs, DataType: req.DataType, Dimension: req.Dimension, DistanceMetric: req.DistanceMetric, CreationTime: int64(1700000000 + m.calls[action]), MetadataConfiguration: req.MetadataConfiguration}
		if m.badIndexCreate {
			_, _ = w.Write([]byte(`{"indexQcs":`))
			return true
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"indexQcs": qcs})
	case "GetIndex":
		if m.indexes[key] == nil {
			cosVectorMockError(w, 404)
			return true
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"index": m.indexes[key]})
	case "DeleteIndex":
		if m.indexes[key] == nil {
			cosVectorMockError(w, 404)
			return true
		}
		delete(m.indexes, key)
		w.WriteHeader(200)
	case "ListIndexes":
		names := []string{}
		for key, index := range m.indexes {
			if index.VectorBucketName == req.VectorBucketName {
				names = append(names, key)
			}
		}
		sort.Strings(names)
		start := 0
		if req.NextToken != "" {
			var err error
			start, err = strconv.Atoi(req.NextToken)
			if err != nil || start < 0 || start > len(names) {
				cosVectorMockError(w, 400)
				return true
			}
		}
		end := len(names)
		if m.listPageSize > 0 && start+m.listPageSize < end {
			end = start + m.listPageSize
		}
		indexes := []cosSDK.IndexBrief{}
		for _, key := range names[start:end] {
			index := m.indexes[key]
			indexes = append(indexes, cosSDK.IndexBrief{IndexName: index.IndexName, IndexQcs: index.IndexQcs, VectorBucketName: index.VectorBucketName, CreationTime: index.CreationTime})
		}
		if m.badListIdentity && len(indexes) > 0 {
			indexes[0].VectorBucketName = "wrong-account-123"
		}
		token := ""
		if end < len(names) {
			token = strconv.Itoa(end)
		}
		if m.badListToken {
			token = "0"
		}
		_ = json.NewEncoder(w).Encode(cosSDK.ListIndexesResult{Indexes: indexes, NextToken: token})
	case "PutVectorBucketPolicy":
		if _, err := normalizeCosVectorPolicy(req.Policy); err != nil {
			cosVectorMockError(w, 400)
			return true
		}
		var policy any
		_ = json.Unmarshal([]byte(req.Policy), &policy)
		text, _ := json.MarshalIndent(policy, "", "  ")
		m.policies[req.VectorBucketName] = string(text)
		w.WriteHeader(200)
	case "GetVectorBucketPolicy":
		policy, ok := m.policies[req.VectorBucketName]
		if !ok {
			cosVectorMockError(w, 404)
			return true
		}
		if m.badPolicyRead {
			_, _ = w.Write([]byte(`{"policy":`))
			return true
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"policy": policy})
	case "DeleteVectorBucketPolicy":
		if _, ok := m.policies[req.VectorBucketName]; !ok {
			cosVectorMockError(w, 404)
			return true
		}
		delete(m.policies, req.VectorBucketName)
		w.WriteHeader(200)
	}
	return true
}

func cosVectorTestPolicy(bucketQCS, effect string) string {
	policy := map[string]any{"version": "2.0", "statement": []any{map[string]any{"effect": effect, "principal": map[string]any{"qcs": []string{"qcs::cam::uin/100000000001:uin/100000000002"}}, "action": []string{"name/cos:QueryVectors"}, "resource": []string{bucketQCS + "/*"}}}}
	data, _ := json.Marshal(policy)
	return string(data)
}

func TestCosVectorManagementCRUD(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	ctx := context.Background()
	bucket := "tf-management-1250000000"
	qcs, err := s.CreateBucket(ctx, bucket, nil)
	if err != nil {
		t.Fatal(err)
	}
	opt := &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: "test-index", DataType: "float32", Dimension: 3, DistanceMetric: "cosine", MetadataConfiguration: &cosSDK.MetadataConfiguration{NonFilterableMetadataKeys: []string{"text"}}}
	indexQCS, err := s.CreateIndex(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}
	if b, i, err := parseCosVectorIndexImportID(indexQCS, s.region); err != nil || b != bucket || i != opt.IndexName {
		t.Fatal("invalid index QCS")
	}
	if _, err := s.GetIndex(ctx, bucket, opt.IndexName); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteBucket(ctx, bucket); err == nil {
		t.Fatal("nonempty bucket deleted")
	}
	policy := cosVectorTestPolicy(qcs, "allow")
	if err := s.PutBucketPolicy(ctx, bucket, policy); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetBucketPolicy(ctx, bucket)
	if err != nil || !cosVectorPoliciesEqual(policy, got) {
		t.Fatal("policy JSON format mismatch")
	}
	if err := s.PutBucketPolicy(ctx, bucket, cosVectorTestPolicy(qcs, "deny")); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteBucketPolicy(ctx, bucket); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetBucketPolicy(ctx, bucket); !isCosVectorBucketNotFound(err) {
		t.Fatal("missing policy not recognized")
	}
	if err := s.DeleteIndex(ctx, bucket, opt.IndexName); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetIndex(ctx, bucket, opt.IndexName); !isCosVectorBucketNotFound(err) {
		t.Fatal("missing index not recognized")
	}
	if err := s.DeleteBucket(ctx, bucket); err != nil || m.count() != 0 {
		t.Fatal("cleanup failed")
	}
}

func TestCosVectorForceDestroySafety(t *testing.T) {
	for _, scenario := range []string{"disabled", "enabled", "bad-token", "wrong-identity", "forbidden-list", "forbidden-delete"} {
		t.Run(scenario, func(t *testing.T) {
			s, m := newCosVectorBucketTestService(t)
			ctx := context.Background()
			bucket := "tf-force-1250000000"
			if _, err := s.CreateBucket(ctx, bucket, nil); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"index-one", "index-two", "index-three"} {
				if _, err := s.CreateIndex(ctx, &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: name, DataType: "float32", Dimension: 3, DistanceMetric: "cosine"}); err != nil {
					t.Fatal(err)
				}
			}
			m.listPageSize = 1
			m.badListToken = scenario == "bad-token"
			m.badListIdentity = scenario == "wrong-identity"
			if scenario == "forbidden-list" {
				m.setFailure("ListIndexes", 403)
			}
			if scenario == "forbidden-delete" {
				m.setFailure("DeleteIndex", 403)
			}
			r := &CosVectorBucketResource{service: s}
			data := CosVectorBucketResourceModel{Region: types.StringValue(s.region), ForceDestroy: types.BoolValue(scenario != "disabled")}
			info, err := s.GetBucket(ctx, bucket)
			if err != nil {
				t.Fatal(err)
			}
			data.setBucket(info)
			state := cosVectorTestState(t, r, data)
			var resp resource.DeleteResponse
			r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
			if scenario == "enabled" {
				if resp.Diagnostics.HasError() || m.count() != 0 || m.callCount("DeleteIndex") != 3 || m.callCount("ListIndexes") != 3 {
					t.Fatal("paginated force destroy failed")
				}
			} else {
				if !resp.Diagnostics.HasError() || m.count() != 1 {
					t.Fatal("unsafe delete succeeded")
				}
				if scenario != "forbidden-delete" && m.callCount("DeleteIndex") != 0 {
					t.Fatal("invalid list caused partial deletion")
				}
			}
		})
	}
}

func TestCosVectorManagementReadErrorsKeepState(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	ctx := context.Background()
	bucket := "tf-errors-1250000000"
	qcs, err := s.CreateBucket(ctx, bucket, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateIndex(ctx, &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: "test-index", DataType: "float32", Dimension: 3, DistanceMetric: "cosine"}); err != nil {
		t.Fatal(err)
	}
	policy := cosVectorTestPolicy(qcs, "allow")
	if err := s.PutBucketPolicy(ctx, bucket, policy); err != nil {
		t.Fatal(err)
	}
	index, err := s.GetIndex(ctx, bucket, "test-index")
	if err != nil {
		t.Fatal(err)
	}
	indexResource := &CosVectorIndexResource{service: s}
	indexData := CosVectorIndexResourceModel{Region: types.StringValue(s.region)}
	indexData.setIndex(index)
	var indexSchema resource.SchemaResponse
	indexResource.Schema(ctx, resource.SchemaRequest{}, &indexSchema)
	indexState := tfsdk.State{Schema: indexSchema.Schema}
	if diags := indexState.Set(ctx, &indexData); diags.HasError() {
		t.Fatal(diags)
	}
	policyResource := &CosVectorBucketPolicyResource{service: s}
	policyData := CosVectorBucketPolicyResourceModel{Id: types.StringValue(qcs), VectorBucketQcs: types.StringValue(qcs), VectorBucketName: types.StringValue(bucket), Region: types.StringValue(s.region), Policy: cosVectorPolicyString(policy)}
	var policySchema resource.SchemaResponse
	policyResource.Schema(ctx, resource.SchemaRequest{}, &policySchema)
	policyState := tfsdk.State{Schema: policySchema.Schema}
	if diags := policyState.Set(ctx, &policyData); diags.HasError() {
		t.Fatal(diags)
	}
	for _, status := range []int{403, 500} {
		m.setFailure("GetIndex", status)
		indexResp := resource.ReadResponse{State: indexState}
		indexResource.Read(ctx, resource.ReadRequest{State: indexState}, &indexResp)
		if !indexResp.Diagnostics.HasError() || !indexResp.State.Raw.Equal(indexState.Raw) {
			t.Fatal("index error lost state")
		}
		m.setFailure("GetVectorBucketPolicy", status)
		policyResp := resource.ReadResponse{State: policyState}
		policyResource.Read(ctx, resource.ReadRequest{State: policyState}, &policyResp)
		if !policyResp.Diagnostics.HasError() || !policyResp.State.Raw.Equal(policyState.Raw) {
			t.Fatal("policy error lost state")
		}
	}
	if _, err := policyResource.target(strings.Replace(qcs, s.region, "ap-shanghai", 1)); err == nil {
		t.Fatal("cross-region policy target accepted")
	}
	if _, _, err := parseCosVectorIndexImportID(strings.Replace(index.IndexQcs, "uid/1250000000", "uid/123", 1), s.region); err == nil {
		t.Fatal("mismatched account index import accepted")
	}
}

func TestCosVectorPolicyNumberPrecision(t *testing.T) {
	policy := func(number string) string {
		return `{"version":"2.0","statement":[{"condition":{"number":` + number + `}}]}`
	}
	if cosVectorPoliciesEqual(policy("9007199254740992"), policy("9007199254740993")) {
		t.Fatal("large numeric condition was rounded")
	}
	for _, pair := range [][2]string{{"1", "1.0"}, {"1e2", "100.00"}, {"-0.0", "0"}, {"0.0010", "1e-3"}} {
		if !cosVectorPoliciesEqual(policy(pair[0]), policy(pair[1])) {
			t.Fatal("equivalent numeric syntax causes diff")
		}
	}
}

func TestCosVectorIndexEmptyMetadata(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	ctx := context.Background()
	bucket := "tf-empty-metadata-1250000000"
	if _, err := s.CreateBucket(ctx, bucket, nil); err != nil {
		t.Fatal(err)
	}
	opt := &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: "empty-metadata", DataType: "float32", Dimension: 3, DistanceMetric: "cosine"}
	if _, err := s.CreateIndex(ctx, opt); err != nil {
		t.Fatal(err)
	}
	m.mu.Lock()
	m.indexes[bucket+"/empty-metadata"].MetadataConfiguration = &cosSDK.MetadataConfiguration{}
	m.mu.Unlock()
	info, err := s.GetIndex(ctx, bucket, opt.IndexName)
	if err != nil || info.MetadataConfiguration != nil {
		t.Fatal("empty returned metadata did not normalize to null")
	}
	opt.MetadataConfiguration = &cosSDK.MetadataConfiguration{}
	if err := validateCosVectorIndexOptions(opt); err == nil {
		t.Fatal("explicit empty creation metadata accepted")
	}
}

func TestCosVectorPolicySemanticEquality(t *testing.T) {
	left := `{"version":"2.0","statement":[{"effect":"allow","action":["name/cos:GetVectors","name/cos:QueryVectors"],"resource":["a","b"],"principal":{"qcs":["x","y"]}}]}`
	right := `{ "statement": [{"principal":{"qcs":["y","x"]},"resource":["b","a"],"action":["name/cos:QueryVectors","name/cos:GetVectors"],"effect":"allow"}], "version":"2.0" }`
	if !cosVectorPoliciesEqual(left, right) {
		t.Fatal("format or set order causes diff")
	}
	if cosVectorPoliciesEqual(left, strings.Replace(right, "allow", "deny", 1)) {
		t.Fatal("policy meaning change ignored")
	}
	if _, err := normalizeCosVectorPolicy(`{"Version":"2012-10-17","Statement":[]}`); err == nil {
		t.Fatal("AWS syntax accepted as COS policy")
	}
}
