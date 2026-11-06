package ags

import (
	"context"
	"errors"
	"regexp"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

var cloudRequestID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// requestDelete returns the public request receipt only after a successful
// response to this exact deletion. Missing/malformed receipts cannot authorize
// empty-list confirmation. Mutations are never replayed here.
func (s *Service) requestDelete(ctx context.Context, id string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if s == nil || s.client == nil || id == "" {
		return "", errors.New("an AGS client and exact Tool ID are required for deletion")
	}
	req := sdk.NewDeleteSandboxToolRequest()
	req.ToolId = pointer(id)
	resp, err := s.client.DeleteSandboxToolWithContext(ctx, req)
	if confirmedAbsent(err) {
		return "", ErrConfirmedAbsent
	}
	if err != nil {
		return "", APIError("DeleteSandboxTool", err)
	}
	if resp == nil || resp.Response == nil || !cloudRequestID.MatchString(value(resp.Response.RequestId)) {
		return "", nil
	}
	return value(resp.Response.RequestId), nil
}

// DeleteAndWait keeps a deletion and its confirmation on the same client and
// exact ID. Only this path accepts a complete successful empty point query
// after a public Delete acknowledgement. Lookup/import stay strict; Resource
// refresh has its separate, explicitly approved empty-result policy.
// The caller supplies the operation deadline; no acknowledgement is retained
// across operations, cancellation, credential changes or process restarts.
func (s *Service) DeleteAndWait(ctx context.Context, id string) error {
	receipt, err := s.requestDelete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrConfirmedAbsent) {
			return nil
		}
		return err
	}
	if receipt == "" {
		// An explicit NotFound can still complete a deletion without an ACK,
		// but a successful empty list alone cannot substitute for a receipt.
		return s.WaitDeleted(ctx, id)
	}
	for {
		req := sdk.NewDescribeSandboxToolListRequest()
		req.ToolIds = []*string{pointer(id)}
		resp, err := s.describe(ctx, req)
		if errors.Is(err, ErrConfirmedAbsent) {
			return nil
		}
		if err != nil {
			return err
		}
		if !cloudRequestID.MatchString(value(resp.RequestId)) || resp.TotalCount == nil || resp.SandboxToolSet == nil || *resp.TotalCount != int64(len(resp.SandboxToolSet)) {
			return errors.New("delete confirmation omitted a request receipt or complete point-query metadata; state is retained")
		}
		if len(resp.SandboxToolSet) == 0 {
			return nil
		}
		if _, err := pointTool(resp, id); err != nil {
			return err
		}
		if err := sleepContext(ctx, s.initialDelay); err != nil {
			return err
		}
	}
}
