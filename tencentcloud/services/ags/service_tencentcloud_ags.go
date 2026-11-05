package ags

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// SDKClient deliberately exposes only Tool operations. A Resource cannot start,
// stop, resume or connect to instances through this interface.
type SDKClient interface {
	CreateSandboxToolWithContext(context.Context, *sdk.CreateSandboxToolRequest) (*sdk.CreateSandboxToolResponse, error)
	DescribeSandboxToolListWithContext(context.Context, *sdk.DescribeSandboxToolListRequest) (*sdk.DescribeSandboxToolListResponse, error)
	UpdateSandboxToolWithContext(context.Context, *sdk.UpdateSandboxToolRequest) (*sdk.UpdateSandboxToolResponse, error)
	DeleteSandboxToolWithContext(context.Context, *sdk.DeleteSandboxToolRequest) (*sdk.DeleteSandboxToolResponse, error)
}

var (
	ErrConfirmedAbsent  = errors.New("AGS confirmed that the sandbox Tool does not exist")
	ErrAmbiguousAbsence = errors.New("AGS returned an empty Tool list; absence and permission filtering cannot be distinguished; state is retained")
)

// Service uses bounded read retries and never silently retries a mutation.
type Service struct {
	client       SDKClient
	initialDelay time.Duration
	maxDelay     time.Duration
}

func NewService(client SDKClient) *Service {
	if concrete, ok := client.(*sdk.Client); ok {
		client = &presenceClient{Client: concrete}
	}
	return &Service{client: client, initialDelay: 2 * time.Second, maxDelay: 30 * time.Second}
}

func pointer[T any](v T) *T { return &v }
func value[T any](v *T) (zero T) {
	if v != nil {
		return *v
	}
	return zero
}

// APIError excludes API message text and request bodies, which can contain
// caller configuration. Request IDs remain available for service investigation.
func APIError(action string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", action, err)
	}
	var api *tcerr.TencentCloudSDKError
	if errors.As(err, &api) {
		return fmt.Errorf("%s failed: code=%s request_id=%s", action, api.GetCode(), api.GetRequestId())
	}
	return fmt.Errorf("%s failed (%T); the remote outcome may need read-only reconciliation", action, err)
}

func confirmedAbsent(err error) bool {
	var api *tcerr.TencentCloudSDKError
	return errors.As(err, &api) && api.GetCode() == "ResourceNotFound.SandboxTool"
}

func retryableRead(err error) bool {
	var api *tcerr.TencentCloudSDKError
	if !errors.As(err, &api) {
		return false
	}
	code := api.GetCode()
	return code == "RequestLimitExceeded" || strings.HasPrefix(code, "RequestLimitExceeded.") || code == "InternalError" || strings.HasPrefix(code, "InternalError.") || code == "ServiceUnavailable"
}

func sleepContext(ctx context.Context, delay time.Duration) error {
	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func (s *Service) describe(ctx context.Context, req *sdk.DescribeSandboxToolListRequest) (*sdk.DescribeSandboxToolListResponseParams, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("AGS client is not configured")
	}
	for attempt := 0; attempt < 3; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		resp, err := s.client.DescribeSandboxToolListWithContext(ctx, req)
		if err == nil {
			if resp == nil || resp.Response == nil {
				return nil, errors.New("DescribeSandboxToolList returned no response; state is retained")
			}
			return resp.Response, nil
		}
		if confirmedAbsent(err) && len(req.ToolIds) == 1 {
			return nil, ErrConfirmedAbsent
		}
		if !retryableRead(err) || attempt == 2 {
			return nil, APIError("DescribeSandboxToolList", err)
		}
		if err := sleepContext(ctx, time.Duration(attempt+1)*200*time.Millisecond); err != nil {
			return nil, err
		}
	}
	return nil, errors.New("read retry budget exhausted")
}

// Lookup returns an object or an explicit error; a filtered empty list is never
// a successful nil result. Callers must opt in to handling ErrConfirmedAbsent.
func (s *Service) Lookup(ctx context.Context, id string) (*sdk.SandboxTool, error) {
	if id == "" {
		return nil, errors.New("Tool ID is required for an exact read")
	}
	req := sdk.NewDescribeSandboxToolListRequest()
	req.ToolIds = []*string{pointer(id)}
	resp, err := s.describe(ctx, req)
	if err != nil {
		return nil, err
	}
	return pointTool(resp, id)
}

func pointTool(resp *sdk.DescribeSandboxToolListResponseParams, id string) (*sdk.SandboxTool, error) {
	if len(resp.SandboxToolSet) == 0 {
		return nil, ErrAmbiguousAbsence
	}
	if len(resp.SandboxToolSet) != 1 || resp.SandboxToolSet[0] == nil || value(resp.SandboxToolSet[0].ToolId) != id {
		return nil, errors.New("point query returned a nonunique or mismatched Tool ID; state is retained")
	}
	tool := resp.SandboxToolSet[0]
	if value(tool.ToolName) == "" || value(tool.ToolType) == "" || value(tool.Status) == "" || tool.NetworkConfiguration == nil {
		return nil, errors.New("point query omitted required Tool configuration; state is retained")
	}
	switch value(tool.Status) {
	case "CREATING", "ACTIVE", "DELETING", "FAILED":
	default:
		return nil, errors.New("point query returned an unknown Tool status; state is retained")
	}
	return tool, nil
}

func (s *Service) FindByName(ctx context.Context, name string) (*sdk.SandboxTool, error) {
	if name == "" {
		return nil, errors.New("an exact Tool name is required")
	}
	const limit = int64(100)
	seen := map[string]bool{}
	var matches []*sdk.SandboxTool
	var total *int64
	for page := int64(0); page < 1000; page++ {
		req := sdk.NewDescribeSandboxToolListRequest()
		req.Limit, req.Offset = pointer(limit), pointer(page*limit)
		req.Filters = []*sdk.Filter{{Name: pointer("tool-name"), Values: []*string{pointer(name)}}}
		resp, err := s.describe(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp.TotalCount == nil || *resp.TotalCount < 0 || len(resp.SandboxToolSet) > int(limit) {
			return nil, errors.New("list response has invalid pagination metadata")
		}
		if total != nil && *total != *resp.TotalCount {
			return nil, errors.New("Tool list changed while reading; retry the exact-name query")
		}
		total = resp.TotalCount
		newIDs := 0
		for _, tool := range resp.SandboxToolSet {
			if tool == nil || value(tool.ToolId) == "" {
				return nil, errors.New("Tool list contains an invalid item")
			}
			id := value(tool.ToolId)
			if seen[id] {
				continue
			}
			seen[id], newIDs = true, newIDs+1
			if value(tool.ToolName) == name {
				matches = append(matches, tool)
			}
		}
		if len(matches) > 1 {
			return nil, errors.New("more than one visible Tool has the exact requested name; use tool_id")
		}
		end := page*limit + int64(len(resp.SandboxToolSet))
		if len(resp.SandboxToolSet) == 0 && end < *total || len(resp.SandboxToolSet) > 0 && newIDs == 0 {
			return nil, errors.New("Tool pagination made no progress; refusing a partial result")
		}
		if end >= *total {
			if len(matches) != 1 {
				return nil, errors.New("no visible exact Tool match; verify account, region, name and permissions")
			}
			return s.Lookup(ctx, value(matches[0].ToolId))
		}
		if len(resp.SandboxToolSet) < int(limit) {
			return nil, errors.New("Tool list ended before its reported total; refusing a partial result")
		}
	}
	return nil, errors.New("exact-name query exceeded its bounded page budget; use tool_id")
}

func NewClientToken() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func (s *Service) Create(ctx context.Context, req *sdk.CreateSandboxToolRequest) (string, error) {
	if req.ClientToken == nil || *req.ClientToken == "" {
		return "", errors.New("a stable token is required before creating a Tool")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	resp, err := s.client.CreateSandboxToolWithContext(ctx, req)
	// Preserve any returned identity even when the SDK also reports an error.
	var id string
	if resp != nil && resp.Response != nil {
		id = value(resp.Response.ToolId)
	}
	if err != nil {
		return id, APIError("CreateSandboxTool", err)
	}
	if id == "" {
		return "", errors.New("CreateSandboxTool returned no ID; do not create again with a new token; reconcile the exact owned name and import")
	}
	return id, nil
}

func (s *Service) Update(ctx context.Context, req *sdk.UpdateSandboxToolRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err := s.client.UpdateSandboxToolWithContext(ctx, req)
	return APIError("UpdateSandboxTool", err)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.requestDelete(ctx, id)
	return err
}

// WaitReady checks both status and configuration. An ACTIVE response containing
// the previous configuration is not an update acknowledgement.
func (s *Service) WaitReady(ctx context.Context, id string, creating bool, matches func(*sdk.SandboxTool) (bool, error)) (*sdk.SandboxTool, error) {
	delay := s.initialDelay
	var last *sdk.SandboxTool
	for {
		tool, err := s.Lookup(ctx, id)
		if err != nil {
			if !creating || !errors.Is(err, ErrAmbiguousAbsence) && !errors.Is(err, ErrConfirmedAbsent) {
				return last, err
			}
		} else {
			last = tool
			if value(tool.Status) == "FAILED" {
				return last, fmt.Errorf("Tool %s reached FAILED: %s", id, value(tool.StatusReason))
			}
			if value(tool.Status) == "DELETING" {
				return last, errors.New("Tool is being deleted while awaiting readiness")
			}
			if value(tool.Status) == "ACTIVE" {
				ok, e := matches(tool)
				if e != nil {
					return last, e
				}
				if ok {
					return last, nil
				}
			}
		}
		if err := sleepContext(ctx, delay); err != nil {
			return last, fmt.Errorf("Tool %s did not reach the requested readable configuration: %w", id, err)
		}
		delay *= 2
		if delay > s.maxDelay {
			delay = s.maxDelay
		}
		if delay > time.Millisecond {
			if n, e := rand.Int(rand.Reader, big.NewInt(int64(delay/5)+1)); e == nil {
				delay -= time.Duration(n.Int64())
			}
		}
	}
}

func (s *Service) WaitDeleted(ctx context.Context, id string) error {
	for {
		_, err := s.Lookup(ctx, id)
		if errors.Is(err, ErrConfirmedAbsent) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := sleepContext(ctx, s.initialDelay); err != nil {
			return err
		}
	}
}

const creatorSystemTag = "qcs:tag:createdBy"

// UserTags excludes exactly the documented immutable creator system tag. Other
// unexpected keys remain observable drift; there is no broad prefix suppression.
// https://cloud.tencent.com/document/product/651/98750
func UserTags(tags []*sdk.Tag) (map[string]string, error) {
	out := map[string]string{}
	for _, tag := range tags {
		if tag == nil || tag.Key == nil || tag.Value == nil {
			return nil, errors.New("Tool response contains an incomplete tag")
		}
		key := *tag.Key
		if key == creatorSystemTag {
			continue
		}
		if _, exists := out[key]; exists {
			return nil, errors.New("Tool response contains duplicate tag keys")
		}
		out[key] = *tag.Value
	}
	return out, nil
}

func expandTags(tags map[string]string) []*sdk.Tag {
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]*sdk.Tag, 0, len(keys))
	for _, key := range keys {
		out = append(out, &sdk.Tag{Key: pointer(key), Value: pointer(tags[key])})
	}
	return out
}
