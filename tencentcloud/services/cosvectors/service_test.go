package cosvectors

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

type cosVectorBucketMock struct {
	mu                 sync.Mutex
	buckets            map[string]*cosSDK.VectorBucketInfo
	failures           map[string]int
	calls              map[string]int
	nonEmpty           bool
	cleanupPending     int
	invisibleReads     int
	pendingReads       map[string]int
	malformedCreate    bool
	lostCreateResponse bool
	invalidRead        bool
	malformedNotFound  bool
	omitEncryption     bool
	transientFailures  map[string][]int
	indexes            map[string]*cosSDK.IndexInfo
	policies           map[string]string
	listPageSize       int
	badListToken       bool
	badListIdentity    bool
	badIndexCreate     bool
	badPolicyRead      bool
	events             []string
}

func newCosVectorBucketTestService(t *testing.T) (*CosVectorService, *cosVectorBucketMock) {
	t.Helper()
	mock := &cosVectorBucketMock{buckets: map[string]*cosSDK.VectorBucketInfo{}, failures: map[string]int{}, calls: map[string]int{}, pendingReads: map[string]int{}, transientFailures: map[string][]int{}}
	server := httptest.NewServer(mock)
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	shared := &connectivity.TencentCloudClient{Region: "ap-guangzhou", Credential: common.NewTokenCredential("local-id", "local-key", "local-token")}
	service, err := NewCosVectorService(shared)
	if err != nil {
		t.Fatal(err)
	}
	// 测试只替换网络传输，继续使用正式 SDK 的 V5 签名和 JSON 编解码。
	service.client = cosSDK.NewClient(&cosSDK.BaseURL{VectorURL: u}, &http.Client{
		Transport:     &cosSDK.AuthorizationTransport{SecretID: "local-id", SecretKey: "local-key", SessionToken: "local-token", Transport: &connectivity.CosVectorResponseTransport{Transport: http.DefaultTransport}},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
	})
	service.client.Conf.RetryOpt.Count = 1
	service.client.UserAgent = "Terraform-local-test"
	return service, mock
}

func (m *cosVectorBucketMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-cos-request-id", "mock-request-id")
	if r.Method != http.MethodPost || r.URL.RawQuery != "" || r.Header.Get("Content-Type") != "application/json" || !strings.HasPrefix(r.Header.Get("Authorization"), "q-sign-algorithm=sha1&") || r.Header.Get("x-cos-security-token") != "local-token" || r.Header.Get("User-Agent") != "Terraform-local-test" {
		cosVectorMockError(w, http.StatusForbidden)
		return
	}
	action := strings.TrimPrefix(r.URL.Path, "/")
	m.calls[action]++
	if statuses := m.transientFailures[action]; len(statuses) > 0 {
		m.transientFailures[action] = statuses[1:]
		cosVectorMockError(w, statuses[0])
		return
	}
	if status := m.failures[action]; status != 0 {
		cosVectorMockError(w, status)
		return
	}
	if m.serveVectorManagement(w, r, action) {
		return
	}
	var request struct {
		Name       string                         `json:"vectorBucketName"`
		Encryption *cosSDK.VectorEncryptionConfig `json:"encryptionConfiguration"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil || validateCosVectorBucketName(request.Name) != nil {
		cosVectorMockError(w, 400)
		return
	}
	switch action {
	case "CreateVectorBucket":
		if m.buckets[request.Name] != nil {
			cosVectorMockError(w, 409)
			return
		}
		appid := request.Name[strings.LastIndex(request.Name, "-")+1:]
		qcs := "qcs::cosvector:ap-guangzhou:uid/" + appid + ":bucket/" + request.Name
		m.buckets[request.Name] = &cosSDK.VectorBucketInfo{VectorBucketName: request.Name, VectorBucketQcs: qcs, CreationTime: int64(1700000000 + m.calls[action]), EncryptionConfiguration: request.Encryption}
		m.pendingReads[request.Name] = m.invisibleReads
		if m.lostCreateResponse {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				_ = conn.Close()
			}
			return
		}
		if m.malformedCreate {
			_, _ = w.Write([]byte(`{"vectorBucketQcs":`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"vectorBucketQcs": qcs})
	case "GetVectorBucket":
		bucket := m.buckets[request.Name]
		if bucket == nil {
			if m.malformedNotFound {
				w.Header().Set("x-cos-error-code", "NotFoundException")
				w.WriteHeader(404)
				_, _ = w.Write([]byte(`{"message":`))
				return
			}
			cosVectorMockError(w, 404)
			return
		}
		if m.pendingReads[request.Name] > 0 {
			m.pendingReads[request.Name]--
			cosVectorMockError(w, 404)
			return
		}
		if m.invalidRead {
			_, _ = w.Write([]byte(`{"vectorBucket":{}}`))
			return
		}
		if m.omitEncryption {
			copy := *bucket
			copy.EncryptionConfiguration = nil
			_ = json.NewEncoder(w).Encode(map[string]any{"vectorBucket": &copy})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"vectorBucket": bucket})
	case "DeleteVectorBucket":
		if m.buckets[request.Name] == nil {
			cosVectorMockError(w, 404)
			return
		}
		if m.nonEmpty || m.hasIndexes(request.Name) {
			cosVectorMockError(w, 409)
			return
		}
		if m.cleanupPending > 0 {
			m.cleanupPending--
			w.Header().Set("x-cos-error-code", "ConflictException")
			w.WriteHeader(409)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "bucket not empty"})
			return
		}
		delete(m.buckets, request.Name)
		w.WriteHeader(200)
	default:
		cosVectorMockError(w, 404)
	}
}

func cosVectorMockError(w http.ResponseWriter, status int) {
	code := map[int]string{400: "ValidationException", 403: "AccessDeniedException", 404: "NotFoundException", 409: "ConflictException", 429: "ThrottlingException", 500: "InternalServerException", 503: "ServiceUnavailableException"}[status]
	w.Header().Set("x-cos-error-code", code)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "local-key must not appear in diagnostics"})
}

func (m *cosVectorBucketMock) setFailure(action string, status int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failures[action] = status
}
func (m *cosVectorBucketMock) drop(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.buckets, name)
}
func (m *cosVectorBucketMock) count() int { m.mu.Lock(); defer m.mu.Unlock(); return len(m.buckets) }
func (m *cosVectorBucketMock) callCount(action string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls[action]
}

func TestCosVectorServiceCRUD(t *testing.T) {
	service, mock := newCosVectorBucketTestService(t)
	ctx := context.Background()
	name := "tf-vector-1250000000"
	if _, err := service.CreateBucket(ctx, name, &cosSDK.VectorEncryptionConfig{SseType: "AES256"}); err != nil {
		t.Fatal(err)
	}
	bucket, err := service.GetBucket(ctx, name)
	if err != nil || bucket.EncryptionConfiguration.SseType != "AES256" {
		t.Fatalf("unexpected Get: %v", err)
	}
	if _, err := service.CreateBucket(ctx, name, nil); err == nil || isCosVectorBucketCreateAccepted(err) {
		t.Fatal("conflict was adopted")
	}
	for _, status := range []int{403, 409, 500} {
		mock.setFailure("GetVectorBucket", status)
		if _, err := service.GetBucket(ctx, name); err == nil || isCosVectorBucketNotFound(err) || strings.Contains(err.Error(), "local-key") {
			t.Fatal("error classification or redaction failed")
		}
	}
	mock.setFailure("GetVectorBucket", 0)
	if err := service.DeleteBucket(ctx, name); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetBucket(ctx, name); !isCosVectorBucketNotFound(err) {
		t.Fatal("missing bucket was not recognized")
	}
	if err := service.DeleteBucket(ctx, name); !isCosVectorBucketNotFound(err) {
		t.Fatal("delete missing bucket was not recognized")
	}
	if mock.callCount("CreateVectorBucket") != 2 || mock.count() != 0 {
		t.Fatal("unexpected retries or leftover resource")
	}
}

func TestCosVectorCreateUnknownOutcomeDoesNotRetry(t *testing.T) {
	service, mock := newCosVectorBucketTestService(t)
	mock.lostCreateResponse = true
	name := "tf-unknown-1250000000"
	_, err := service.CreateBucket(context.Background(), name, nil)
	if err == nil || isCosVectorBucketCreateAccepted(err) {
		t.Fatal("unconfirmed creation was treated as successful")
	}
	if mock.callCount("CreateVectorBucket") != 1 || mock.count() != 1 {
		t.Fatal("creation was retried or mock did not create the bucket")
	}
	if _, err := service.GetBucket(context.Background(), name); err != nil {
		t.Fatal("existing bucket could not be queried for import")
	}
}

type cosVectorResponseTestTransport func(*http.Request) (*http.Response, error)

func (f cosVectorResponseTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type cosVectorInterruptedBody struct{}

func (cosVectorInterruptedBody) Read(_ []byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (cosVectorInterruptedBody) Close() error               { return nil }

func TestCosVectorWaitBodyReadFailureRecovers(t *testing.T) {
	service, _ := newCosVectorBucketTestService(t)
	name := "tf-body-retry-1250000000"
	calls := 0
	transport := &connectivity.CosVectorResponseTransport{Transport: cosVectorResponseTestTransport(func(req *http.Request) (*http.Response, error) {
		calls++
		body := io.ReadCloser(cosVectorInterruptedBody{})
		if calls > 1 {
			data, _ := json.Marshal(map[string]any{"vectorBucket": &cosSDK.VectorBucketInfo{VectorBucketName: name, VectorBucketQcs: "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/" + name}})
			body = io.NopCloser(strings.NewReader(string(data)))
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: body}, nil
	})}
	service.client = cosSDK.NewClient(service.client.BaseURL, &http.Client{Transport: transport})
	service.client.Conf.RetryOpt.Count = 1
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := service.WaitBucket(ctx, name, true); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal("body read failure was not retried only by the waiter")
	}
}

func TestCosVectorWaitTransientRecovery(t *testing.T) {
	service, mock := newCosVectorBucketTestService(t)
	name := "tf-wait-1250000000"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := service.CreateBucket(ctx, name, nil); err != nil {
		t.Fatal(err)
	}
	mock.mu.Lock()
	mock.transientFailures["GetVectorBucket"] = []int{500, 429, 503}
	mock.mu.Unlock()
	if _, err := service.WaitBucket(ctx, name, true); err != nil {
		t.Fatal(err)
	}
	if mock.callCount("CreateVectorBucket") != 1 || mock.callCount("GetVectorBucket") != 4 {
		t.Fatal("wait retried create or did not recover")
	}
	if err := service.DeleteBucket(ctx, name); err != nil {
		t.Fatal(err)
	}
	mock.mu.Lock()
	mock.transientFailures["GetVectorBucket"] = []int{503}
	mock.mu.Unlock()
	if _, err := service.WaitBucket(ctx, name, false); err != nil {
		t.Fatal(err)
	}
	if mock.callCount("DeleteVectorBucket") != 1 {
		t.Fatal("wait retried delete")
	}
}

func TestCosVectorWaitDeadlineAndPermanentErrors(t *testing.T) {
	service, mock := newCosVectorBucketTestService(t)
	name := "tf-missing-1250000000"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := service.WaitBucket(ctx, name, true); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline lost: %v", err)
	}
	for _, status := range []int{403, 409} {
		mock.setFailure("GetVectorBucket", status)
		calls := mock.callCount("GetVectorBucket")
		if _, err := service.WaitBucket(context.Background(), name, true); err == nil {
			t.Fatal("permanent error ignored")
		}
		if mock.callCount("GetVectorBucket") != calls+1 {
			t.Fatal("permanent error was retried")
		}
	}
	mock.setFailure("GetVectorBucket", 0)
	mock.mu.Lock()
	mock.malformedNotFound = true
	mock.mu.Unlock()
	calls := mock.callCount("GetVectorBucket")
	if _, err := service.WaitBucket(context.Background(), name, true); err == nil || isCosVectorBucketNotFound(err) || isCosVectorBucketRetryable(err) {
		t.Fatal("malformed response accepted or retried")
	}
	if mock.callCount("GetVectorBucket") != calls+1 {
		t.Fatal("malformed response was retried")
	}
	ctx, canceled := context.WithCancel(context.Background())
	canceled()
	calls = mock.callCount("GetVectorBucket")
	if _, err := service.WaitBucket(ctx, name, true); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation lost")
	}
	if mock.callCount("GetVectorBucket") != calls {
		t.Fatal("canceled wait sent a request")
	}
}

func TestCosVectorCanceledOperationsDoNotSend(t *testing.T) {
	service, mock := newCosVectorBucketTestService(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	name := "tf-canceled-1250000000"
	if _, err := service.CreateBucket(ctx, name, nil); !errors.Is(err, context.Canceled) {
		t.Fatal("create cancellation lost")
	}
	if _, err := service.GetBucket(ctx, name); !errors.Is(err, context.Canceled) {
		t.Fatal("get cancellation lost")
	}
	if err := service.DeleteBucket(ctx, name); !errors.Is(err, context.Canceled) {
		t.Fatal("delete cancellation lost")
	}
	if mock.callCount("CreateVectorBucket")+mock.callCount("GetVectorBucket")+mock.callCount("DeleteVectorBucket") != 0 {
		t.Fatal("canceled operation sent requests")
	}
	for _, tc := range []struct {
		cause     error
		retryable bool
	}{{context.Canceled, false}, {context.DeadlineExceeded, true}, {io.EOF, true}, {errors.New("local-key sensitive detail"), false}} {
		err := service.apiError("GetVectorBucket", nil, &cosSDK.RetryError{Errs: []error{tc.cause}})
		if isCosVectorBucketRetryable(err) != tc.retryable || strings.Contains(err.Error(), "local-key") {
			t.Fatal("unsafe retry classification or diagnostic")
		}
	}
}

func TestCosVectorImportID(t *testing.T) {
	qcs := "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-1250000000"
	for _, id := range []string{"tf-vector-1250000000", qcs} {
		if name, err := parseCosVectorBucketImportID(id, "ap-guangzhou"); err != nil || name != "tf-vector-1250000000" {
			t.Fatalf("valid import rejected: %v", err)
		}
	}
	for _, id := range []string{"missing-appid", "Upper-1250000000", "arn:aws:s3vectors:region:account:bucket/test", strings.Replace(qcs, "ap-guangzhou", "ap-shanghai", 1), strings.Replace(qcs, "uid/1250000000", "uid/123", 1), qcs + "/index/test"} {
		if _, err := parseCosVectorBucketImportID(id, "ap-guangzhou"); err == nil {
			t.Fatal("invalid import accepted")
		}
	}
}

func TestCosVectorErrorClassificationAndRedaction(t *testing.T) {
	service, _ := newCosVectorBucketTestService(t)
	for _, tc := range []struct {
		status int
		code   string
		absent bool
	}{{404, "NotFoundException", true}, {404, "Unknown", false}, {403, "AccessDeniedException", false}, {409, "ConflictException", false}, {500, "NotFoundException", false}} {
		resp := &cosSDK.Response{Response: &http.Response{StatusCode: tc.status, Header: make(http.Header)}}
		resp.Header.Set("x-cos-request-id", "prefix-local-token-suffix")
		err := service.apiError("GetVectorBucket", resp, &cosSDK.VectorErrorResponse{Code: tc.code, Response: resp.Response})
		if isCosVectorBucketNotFound(err) != tc.absent || strings.Contains(err.Error(), "local-token") {
			t.Fatal("unsafe error classification or response header leak")
		}
	}
}
