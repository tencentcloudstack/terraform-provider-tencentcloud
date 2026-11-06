package ags

import (
	"context"
	"errors"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

// ReadManaged follows the repository's Resource refresh policy: a complete,
// successful empty point query removes an already-managed identity. Permission
// filtering that succeeds with the same empty response is indistinguishable;
// explicit API errors and incomplete responses never prove absence.
// Lookup remains conservative for import, Data Sources and operation waiters.
func (s *Service) ReadManaged(ctx context.Context, id string) (*sdk.SandboxTool, error) {
	if id == "" {
		return nil, errors.New("Tool ID is required for a managed Resource read")
	}
	req := sdk.NewDescribeSandboxToolListRequest()
	req.ToolIds = []*string{pointer(id)}
	resp, err := s.describe(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.SandboxToolSet) == 0 && resp.SandboxToolSet != nil && resp.TotalCount != nil && *resp.TotalCount == 0 && cloudRequestID.MatchString(value(resp.RequestId)) {
		return nil, ErrConfirmedAbsent
	}
	return pointTool(resp, id)
}
