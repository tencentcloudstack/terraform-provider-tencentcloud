package cosvectors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

var cosVectorBucketNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*-[0-9]+$`)

type CosVectorService struct {
	client  *cosSDK.Client
	region  string
	secrets []string
}

func NewCosVectorService(client *connectivity.TencentCloudClient) (*CosVectorService, error) {
	if client == nil {
		return nil, errors.New("TencentCloud client is not configured")
	}
	sdk, err := client.UseTencentCosVectorClient()
	if err != nil {
		return nil, err
	}
	return &CosVectorService{client: sdk, region: client.Region, secrets: []string{client.Credential.SecretId, client.Credential.SecretKey, client.Credential.Token}}, nil
}

type cosVectorBucketError struct {
	action, code, requestID  string
	status                   int
	notFound, acceptedCreate bool
	retryable                bool
	bucketNotEmpty           bool
	cause                    error
}

func (e *cosVectorBucketError) Error() string {
	return fmt.Sprintf("%s failed: HTTP %d, code=%s, request_id=%s", e.action, e.status, e.code, e.requestID)
}

func (e *cosVectorBucketError) Unwrap() error { return e.cause }

func isCosVectorBucketRetryable(err error) bool {
	var apiErr *cosVectorBucketError
	return errors.As(err, &apiErr) && apiErr.retryable
}

func isCosVectorBucketNotFound(err error) bool {
	var apiErr *cosVectorBucketError
	return errors.As(err, &apiErr) && apiErr.notFound
}

func isCosVectorBucketCreateAccepted(err error) bool {
	var apiErr *cosVectorBucketError
	return errors.As(err, &apiErr) && apiErr.acceptedCreate
}

func validateCosVectorBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 || !cosVectorBucketNamePattern.MatchString(name) {
		return errors.New("vector_bucket_name must be 3-63 lowercase letters, digits or hyphens, ending in -APPID")
	}
	return nil
}

func parseCosVectorBucketImportID(id, region string) (string, error) {
	if !strings.HasPrefix(id, "qcs:") {
		return id, validateCosVectorBucketName(id)
	}
	prefix := "qcs::cosvector:" + region + ":uid/"
	if !strings.HasPrefix(id, prefix) {
		return "", errors.New("bucket QCS must match the COS Vectors service and provider region")
	}
	parts := strings.Split(strings.TrimPrefix(id, prefix), ":bucket/")
	if len(parts) != 2 {
		return "", errors.New("invalid bucket QCS")
	}
	name := parts[1]
	if err := validateCosVectorBucketName(name); err != nil {
		return "", err
	}
	if parts[0] != name[strings.LastIndex(name, "-")+1:] {
		return "", errors.New("bucket QCS account must match bucket APPID")
	}
	return name, nil
}

func (s *CosVectorService) CreateBucket(ctx context.Context, name string, encryption *cosSDK.VectorEncryptionConfig) (string, error) {
	if err := validateCosVectorBucketName(name); err != nil {
		return "", err
	}
	if encryption != nil && encryption.SseType != "AES256" {
		return "", errors.New("COS Vectors only supports AES256 encryption")
	}
	if err := ratelimit.CheckContext(ctx, "CreateVectorBucket"); err != nil {
		return "", s.apiError("CreateVectorBucket", nil, err)
	}
	result, resp, err := s.client.Vector.CreateVectorBucket(ctx, &cosSDK.CreateVectorBucketOptions{VectorBucketName: name, EncryptionConfiguration: encryption})
	if err != nil || !cosVectorResponseOK(resp) {
		return "", s.apiError("CreateVectorBucket", resp, err)
	}
	if result == nil || !s.validBucketQCS(result.VectorBucketQcs, name) {
		return "", s.apiError("CreateVectorBucket", resp, errors.New("invalid create response"))
	}
	return result.VectorBucketQcs, nil
}

func (s *CosVectorService) GetBucket(ctx context.Context, name string) (*cosSDK.VectorBucketInfo, error) {
	if err := validateCosVectorBucketName(name); err != nil {
		return nil, err
	}
	if err := ratelimit.CheckContext(ctx, "GetVectorBucket"); err != nil {
		return nil, s.apiError("GetVectorBucket", nil, err)
	}
	result, resp, err := s.client.Vector.GetVectorBucket(ctx, &cosSDK.GetVectorBucketOptions{VectorBucketName: name})
	if err != nil || !cosVectorResponseOK(resp) {
		return nil, s.apiError("GetVectorBucket", resp, err)
	}
	if result == nil || result.VectorBucket == nil || result.VectorBucket.VectorBucketName != name || !s.validBucketQCS(result.VectorBucket.VectorBucketQcs, name) {
		return nil, errors.New("GetVectorBucket returned missing or mismatched bucket identity")
	}
	bucket := result.VectorBucket
	if bucket.EncryptionConfiguration != nil && bucket.EncryptionConfiguration.SseType != "AES256" {
		return nil, errors.New("this resource does not support the returned encryption configuration")
	}
	return bucket, nil
}

func (s *CosVectorService) DeleteBucket(ctx context.Context, name string) error {
	if err := validateCosVectorBucketName(name); err != nil {
		return err
	}
	if err := ratelimit.CheckContext(ctx, "DeleteVectorBucket"); err != nil {
		return s.apiError("DeleteVectorBucket", nil, err)
	}
	resp, err := s.client.Vector.DeleteVectorBucket(ctx, &cosSDK.DeleteVectorBucketOptions{VectorBucketName: name})
	if err != nil || !cosVectorResponseOK(resp) {
		return s.apiError("DeleteVectorBucket", resp, err)
	}
	return nil
}

func (s *CosVectorService) validBucketQCS(qcs, name string) bool {
	got, err := parseCosVectorBucketImportID(qcs, s.region)
	return strings.HasPrefix(qcs, "qcs:") && err == nil && got == name
}

func cosVectorResponseOK(resp *cosSDK.Response) bool {
	return resp != nil && resp.Response != nil && resp.StatusCode == http.StatusOK && resp.Header.Get(connectivity.CosVectorInvalidResponseHeader) == ""
}

func (s *CosVectorService) apiError(action string, resp *cosSDK.Response, err error) error {
	result := &cosVectorBucketError{action: action, code: "NetworkOrInvalidResponse", requestID: ""}
	if resp != nil && resp.Response != nil {
		result.status = resp.StatusCode
		result.requestID = s.safeLabel(resp.Header.Get("x-cos-request-id"))
	}
	if apiErr, ok := cosSDK.IsVectorError(err); ok {
		switch apiErr.Code {
		case "NotFoundException", "AccessDeniedException", "ConflictException", "ValidationException", "InternalServerException", "ThrottlingException", "TooManyRequestsException", "ServiceQuotaExceededException", "ServiceUnavailableException":
			result.code = s.safeLabel(apiErr.Code)
		default:
			result.code = "Unknown"
		}
		valid := resp != nil && resp.Header.Get(connectivity.CosVectorInvalidResponseHeader) == ""
		result.notFound = valid && result.status == http.StatusNotFound && apiErr.Code == "NotFoundException"
		result.bucketNotEmpty = valid && action == "DeleteVectorBucket" && result.status == http.StatusConflict && apiErr.Code == "ConflictException" && apiErr.Message == "bucket not empty"
		result.retryable = valid && ((result.status == 429 && (apiErr.Code == "ThrottlingException" || apiErr.Code == "TooManyRequestsException")) ||
			(result.status == 500 && apiErr.Code == "InternalServerException") || (result.status == 503 && apiErr.Code == "ServiceUnavailableException"))
	} else {
		var netErr net.Error
		switch {
		case errors.Is(err, context.Canceled):
			result.code, result.cause = "Canceled", context.Canceled
		case errors.Is(err, context.DeadlineExceeded):
			result.code, result.cause, result.retryable = "Timeout", context.DeadlineExceeded, true
		case errors.As(err, &netErr) && netErr.Timeout():
			result.code, result.retryable = "Timeout", true
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, syscall.ECONNRESET), errors.Is(err, syscall.ECONNREFUSED), errors.Is(err, syscall.EPIPE):
			result.code, result.retryable = "ConnectionFailure", true
		}
	}
	if resp != nil && resp.Header.Get(connectivity.CosVectorInvalidResponseHeader) != "" {
		result.code, result.notFound, result.retryable = "InvalidResponse", false, false
		switch resp.Header.Get(connectivity.CosVectorInvalidResponseHeader) {
		case "ConnectionFailure":
			result.code, result.retryable = "ConnectionFailure", true
		case "Timeout":
			result.code, result.cause, result.retryable = "Timeout", context.DeadlineExceeded, true
		case "Canceled":
			result.code, result.cause = "Canceled", context.Canceled
		}
	}
	result.acceptedCreate = (action == "CreateVectorBucket" || action == "CreateIndex" || action == "PutVectorBucketPolicy") && result.status == http.StatusOK
	return result
}

func (s *CosVectorService) safeLabel(value string) string {
	for _, secret := range s.secrets {
		if secret != "" && strings.Contains(value, secret) {
			return "[redacted]"
		}
	}
	if len(value) > 128 {
		value = value[:128]
	}
	return strings.Map(func(r rune) rune {
		if r >= 32 && r <= 126 {
			return r
		}
		return '?'
	}, value)
}

// WaitBucket 只轮询读取，不重试创建、接管冲突或清空非空桶。
func (s *CosVectorService) WaitBucket(ctx context.Context, name string, exists bool) (*cosSDK.VectorBucketInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	delay := 250 * time.Millisecond
	var lastErr error
	for {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("waiting for COS vector bucket visibility (last error: %v): %w", lastErr, err)
		}
		bucket, err := s.GetBucket(ctx, name)
		if exists && err == nil {
			return bucket, nil
		}
		if !exists && isCosVectorBucketNotFound(err) {
			return nil, nil
		}
		if err != nil && !isCosVectorBucketNotFound(err) && !isCosVectorBucketRetryable(err) {
			return nil, err
		}
		lastErr = err
		timer := time.NewTimer(delay/2 + time.Duration(rand.Int63n(int64(delay/2))))
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("waiting for COS vector bucket visibility (last error: %v): %w", lastErr, ctx.Err())
		case <-timer.C:
		}
		if delay < 2*time.Second {
			delay *= 2
		}
	}
}
