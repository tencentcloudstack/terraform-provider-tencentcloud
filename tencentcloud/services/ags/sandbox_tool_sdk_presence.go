package ags

import (
	"context"
	"encoding/json"
	"errors"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

// presenceClient preserves explicit resets independently of the legacy
// Provider's process-wide SDK OmitEmpty setting. Only official typed fields are
// encoded. The public SDK still performs endpoint selection, signing and I/O;
// the payload is finalized before signing and globals are never changed.
type presenceClient struct{ *sdk.Client }

type presenceRequest struct {
	tchttp.Request
	payload []byte
}

func (r *presenceRequest) MarshalJSON() ([]byte, error) { return r.payload, nil }

func presencePayload(request any, tags []*sdk.Tag, custom *sdk.CustomConfiguration) ([]byte, error) {
	encoded, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return nil, err
	}
	if tags != nil {
		fields["Tags"], err = json.Marshal(tags)
		if err != nil {
			return nil, err
		}
	}
	if custom != nil {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(fields["CustomConfiguration"], &nested); err != nil {
			return nil, err
		}
		// A nil slice omits an update; a non-nil empty slice resets it.
		for _, item := range []struct {
			name    string
			present bool
			value   any
		}{
			{"Command", custom.Command != nil, custom.Command},
			{"Args", custom.Args != nil, custom.Args},
			{"Env", custom.Env != nil, custom.Env},
			{"Ports", custom.Ports != nil, custom.Ports},
		} {
			if item.present {
				nested[item.name], err = json.Marshal(item.value)
				if err != nil {
					return nil, err
				}
			}
		}
		fields["CustomConfiguration"], err = json.Marshal(nested)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(fields)
}

func (c *presenceClient) CreateSandboxToolWithContext(ctx context.Context, req *sdk.CreateSandboxToolRequest) (*sdk.CreateSandboxToolResponse, error) {
	if req == nil {
		return nil, errors.New("CreateSandboxTool requires a typed request")
	}
	c.InitBaseRequest(&req.BaseRequest, "ags", sdk.APIVersion, "CreateSandboxTool")
	if c.GetCredential() == nil {
		return nil, errors.New("CreateSandboxTool requires credentials")
	}
	encoded, err := presencePayload(req, req.Tags, req.CustomConfiguration)
	if err != nil {
		return nil, err
	}
	req.SetContext(ctx)
	resp := sdk.NewCreateSandboxToolResponse()
	err = c.Send(&presenceRequest{Request: req, payload: encoded}, resp)
	return resp, err
}

func (c *presenceClient) UpdateSandboxToolWithContext(ctx context.Context, req *sdk.UpdateSandboxToolRequest) (*sdk.UpdateSandboxToolResponse, error) {
	if req == nil {
		return nil, errors.New("UpdateSandboxTool requires a typed request")
	}
	c.InitBaseRequest(&req.BaseRequest, "ags", sdk.APIVersion, "UpdateSandboxTool")
	if c.GetCredential() == nil {
		return nil, errors.New("UpdateSandboxTool requires credentials")
	}
	encoded, err := presencePayload(req, req.Tags, req.CustomConfiguration)
	if err != nil {
		return nil, err
	}
	req.SetContext(ctx)
	resp := sdk.NewUpdateSandboxToolResponse()
	err = c.Send(&presenceRequest{Request: req, payload: encoded}, resp)
	return resp, err
}
