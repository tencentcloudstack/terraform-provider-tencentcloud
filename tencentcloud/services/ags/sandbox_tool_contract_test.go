package ags

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkjson "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/json"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testClient(t *testing.T, handle func(string, map[string]any) (map[string]any, error)) *sdk.Client {
	t.Helper()
	client, err := sdk.NewClient(common.NewCredential("offline-id", "offline-key"), "ap-guangzhou", profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}
	client.WithHttpTransport(testTransport(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var input map[string]any
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		response, err := handle(req.Header.Get("X-TC-Action"), input)
		if err != nil {
			return nil, err
		}
		response["RequestId"] = "offline-contract-request"
		encoded, err := json.Marshal(map[string]any{"Response": response})
		if err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(encoded)))}, nil
	}))
	return client
}

func testTool(id, name, status string) map[string]any {
	return map[string]any{"ToolId": id, "ToolName": name, "ToolType": "code-interpreter", "Status": status, "Description": "old", "DefaultTimeoutSeconds": 300, "Persistent": false, "NetworkConfiguration": map[string]any{"NetworkMode": "PUBLIC"}, "Tags": []any{}}
}
func pointResult(tool map[string]any) map[string]any {
	return map[string]any{"SandboxToolSet": []any{tool}, "TotalCount": 1}
}
func errorResult(code string) map[string]any {
	return map[string]any{"Error": map[string]any{"Code": code, "Message": "contract test"}}
}
func fastService(client SDKClient) *Service {
	service := NewService(client)
	service.initialDelay, service.maxDelay = 100*time.Microsecond, time.Millisecond
	return service
}

func TestCT001ActualSDKPreservesExplicitEmptyValues(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		t.Run(strconv.FormatBool(explicit), func(t *testing.T) {
			calls := 0
			client := testClient(t, func(action string, body map[string]any) (map[string]any, error) {
				calls++
				if action != "UpdateSandboxTool" {
					t.Fatalf("unexpected action %s", action)
				}
				for _, key := range []string{"Description", "Tags", "CustomConfiguration"} {
					_, present := body[key]
					if present != explicit {
						t.Errorf("%s presence=%v, want %v", key, present, explicit)
					}
				}
				if explicit {
					if body["Description"] != "" {
						t.Error("description did not reach HTTP as an empty string")
					}
					if tags, ok := body["Tags"].([]any); !ok || len(tags) != 0 {
						t.Error("tags did not reach HTTP as []")
					}
					custom, ok := body["CustomConfiguration"].(map[string]any)
					if !ok {
						t.Fatal("missing custom configuration")
					}
					for _, key := range []string{"Command", "Args", "Env", "Ports"} {
						if values, ok := custom[key].([]any); !ok || len(values) != 0 {
							t.Errorf("%s did not reach HTTP as []", key)
						}
					}
				}
				return map[string]any{}, nil
			})
			req := sdk.NewUpdateSandboxToolRequest()
			req.ToolId = pointer("sdt-offline")
			if explicit {
				req.Description = pointer("")
				req.Tags = []*sdk.Tag{}
				req.CustomConfiguration = &sdk.CustomConfiguration{Command: []*string{}, Args: []*string{}, Env: []*sdk.EnvVar{}, Ports: []*sdk.PortConfiguration{}}
			}
			if err := NewService(client).Update(context.Background(), req); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("HTTP calls=%d, want 1", calls)
			}
		})
	}
}

func TestCT001ProviderGlobalOmitEmptyCannotDropAGSClears(t *testing.T) {
	previous := sdkjson.OmitBehaviour
	sdkjson.OmitBehaviour = sdkjson.OmitEmpty // tencentcloud/provider.go init sets this globally.
	t.Cleanup(func() { sdkjson.OmitBehaviour = previous })
	TestCT001ActualSDKPreservesExplicitEmptyValues(t)
}

func TestCT002CreateDoesNotBlindlyRetryOrReplaceToken(t *testing.T) {
	calls := 0
	token, err := NewClientToken()
	if err != nil {
		t.Fatal(err)
	}
	client := testClient(t, func(action string, body map[string]any) (map[string]any, error) {
		calls++
		if action != "CreateSandboxTool" || body["ClientToken"] != token {
			t.Fatal("mutation or token changed")
		}
		return nil, io.ErrUnexpectedEOF
	})
	req := sdk.NewCreateSandboxToolRequest()
	req.ToolName, req.ToolType, req.ClientToken = pointer("owned-offline"), pointer("code-interpreter"), pointer(token)
	id, err := NewService(client).Create(context.Background(), req)
	if err == nil || id != "" || calls != 1 {
		t.Fatalf("ambiguous create: id=%q err=%v calls=%d", id, err, calls)
	}
	if value(req.ClientToken) != token {
		t.Fatal("token changed after an uncertain result")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := NewService(client).Create(ctx, req); !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatal("cancelled create sent another request")
	}
	missing := sdk.NewCreateSandboxToolRequest()
	if _, err := NewService(client).Create(context.Background(), missing); err == nil || calls != 1 {
		t.Fatal("create without stable token was sent")
	}
}

func TestCT003OnlyExplicitToolNotFoundConfirmsAbsence(t *testing.T) {
	cases := []struct {
		name, code string
		empty      bool
		want       error
	}{
		{"confirmed", "ResourceNotFound.SandboxTool", false, ErrConfirmedAbsent},
		{"filtered-empty", "", true, ErrAmbiguousAbsence},
		{"denied", "UnauthorizedOperation", false, nil},
		{"expired", "AuthFailure.TokenFailure", false, nil},
		{"other-resource", "ResourceNotFound.RegistryRecord", false, nil},
		{"throttled", "RequestLimitExceeded", false, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				calls++
				if action != "DescribeSandboxToolList" {
					t.Fatalf("read invoked a mutation: %s", action)
				}
				ids, ok := input["ToolIds"].([]any)
				if !ok || len(ids) != 1 || ids[0] != "sdt-owned" {
					t.Fatal("point query lost exact identity")
				}
				if tc.empty {
					return map[string]any{"SandboxToolSet": []any{}, "TotalCount": 0}, nil
				}
				return errorResult(tc.code), nil
			})
			tool, err := NewService(client).Lookup(context.Background(), "sdt-owned")
			if tool != nil || err == nil {
				t.Fatalf("unexpected success: %v", tool)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
			if tc.want == nil && (errors.Is(err, ErrConfirmedAbsent) || errors.Is(err, ErrAmbiguousAbsence)) {
				t.Fatalf("misclassified API error: %v", err)
			}
			wantCalls := 1
			if tc.code == "RequestLimitExceeded" {
				wantCalls = 3
			}
			if calls != wantCalls {
				t.Fatalf("calls=%d, want %d", calls, wantCalls)
			}
		})
	}
	for _, invalid := range []map[string]any{testTool("wrong-id", "owned", "ACTIVE"), testTool("sdt-owned", "owned", "NEW_UNKNOWN_STATUS"), {"ToolId": "sdt-owned"}} {
		client := testClient(t, func(string, map[string]any) (map[string]any, error) { return pointResult(invalid), nil })
		if _, err := NewService(client).Lookup(context.Background(), "sdt-owned"); err == nil {
			t.Fatal("accepted incomplete, mismatched or unknown response")
		}
	}
}

func TestCT004ExactNamePaginationCannotReturnPartialOrAmbiguousMatch(t *testing.T) {
	for _, mode := range []string{"later-page", "duplicate-page", "duplicate-name", "changed-total"} {
		t.Run(mode, func(t *testing.T) {
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				if action != "DescribeSandboxToolList" {
					t.Fatal("query performed a write")
				}
				if _, point := input["ToolIds"]; point {
					return pointResult(testTool("sdt-target", "target", "ACTIVE")), nil
				}
				offset := int(input["Offset"].(float64))
				if mode == "duplicate-name" {
					return map[string]any{"TotalCount": 2, "SandboxToolSet": []any{testTool("one", "target", "ACTIVE"), testTool("two", "target", "ACTIVE")}}, nil
				}
				if offset > 0 && mode == "later-page" {
					return map[string]any{"TotalCount": 101, "SandboxToolSet": []any{testTool("sdt-target", "target", "ACTIVE")}}, nil
				}
				items := make([]any, 100)
				for i := range items {
					items[i] = testTool("sdt-other-"+strconv.Itoa(i), "target-extra", "ACTIVE")
				}
				total := 201
				if mode == "later-page" {
					total = 101
				}
				if mode == "changed-total" && offset > 0 {
					total = 202
				}
				return map[string]any{"TotalCount": total, "SandboxToolSet": items}, nil
			})
			tool, err := NewService(client).FindByName(context.Background(), "target")
			if mode == "later-page" {
				if err != nil || value(tool.ToolId) != "sdt-target" {
					t.Fatalf("later page lookup failed: %v", err)
				}
			} else if err == nil || tool != nil {
				t.Fatalf("accepted partial or ambiguous %s result", mode)
			}
		})
	}
}

func TestCT005ActiveMustAlsoMatchRequestedConfiguration(t *testing.T) {
	for _, converges := range []bool{true, false} {
		t.Run(strconv.FormatBool(converges), func(t *testing.T) {
			calls := 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				calls++
				if action != "DescribeSandboxToolList" {
					t.Fatal("waiter performed a write")
				}
				tool := testTool("sdt-owned", "owned", "ACTIVE")
				if converges && calls >= 2 {
					tool["Description"] = "new"
				}
				return pointResult(tool), nil
			})
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			tool, err := fastService(client).WaitReady(ctx, "sdt-owned", false, func(v *sdk.SandboxTool) (bool, error) { return value(v.Description) == "new", nil })
			if converges {
				if err != nil || calls < 2 || value(tool.Description) != "new" {
					t.Fatalf("convergence failed: %v", err)
				}
			} else if err == nil || tool == nil || value(tool.Description) != "old" {
				t.Fatal("old ACTIVE configuration reported as success or recovery state lost")
			}
		})
	}
}

func TestCT006TerminalStatesCancellationAndDeletionProtection(t *testing.T) {
	t.Run("failed-ready", func(t *testing.T) {
		client := testClient(t, func(string, map[string]any) (map[string]any, error) {
			tool := testTool("sdt-owned", "owned", "FAILED")
			tool["StatusReason"] = "image preparation failed"
			return pointResult(tool), nil
		})
		tool, err := fastService(client).WaitReady(context.Background(), "sdt-owned", true, func(*sdk.SandboxTool) (bool, error) { return true, nil })
		if err == nil || tool == nil || value(tool.Status) != "FAILED" || !strings.Contains(err.Error(), "image preparation failed") {
			t.Fatal("terminal failure was ignored")
		}
	})
	t.Run("deleted-or-ambiguous", func(t *testing.T) {
		for _, confirmed := range []bool{true, false} {
			client := testClient(t, func(string, map[string]any) (map[string]any, error) {
				if confirmed {
					return errorResult("ResourceNotFound.SandboxTool"), nil
				}
				return map[string]any{"SandboxToolSet": []any{}, "TotalCount": 0}, nil
			})
			err := fastService(client).WaitDeleted(context.Background(), "sdt-owned")
			if confirmed && err != nil || !confirmed && !errors.Is(err, ErrAmbiguousAbsence) {
				t.Fatalf("wrong delete classification: %v", err)
			}
		}
	})
	t.Run("dependency-no-instance-mutation", func(t *testing.T) {
		calls := 0
		client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
			calls++
			if action != "DeleteSandboxTool" {
				t.Fatalf("implicit instance action %s", action)
			}
			return errorResult("UnsupportedOperation.SandboxToolHasActiveInstances"), nil
		})
		if err := fastService(client).Delete(context.Background(), "sdt-owned"); err == nil || calls != 1 {
			t.Fatal("dependency protection hidden or retried")
		}
	})
	t.Run("cancelled-read", func(t *testing.T) {
		calls := 0
		client := testClient(t, func(string, map[string]any) (map[string]any, error) {
			calls++
			return nil, errors.New("unexpected request")
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := NewService(client).Lookup(ctx, "sdt-owned"); !errors.Is(err, context.Canceled) || calls != 0 {
			t.Fatal("cancelled read still sent requests")
		}
	})
}

func TestUserTagsExcludesOnlyDocumentedCreatorTag(t *testing.T) {
	got, err := UserTags([]*sdk.Tag{{Key: pointer(creatorSystemTag), Value: pointer("system")}, {Key: pointer("qcs:other-observable"), Value: pointer("unexpected")}, {Key: pointer("owner"), Value: pointer("test")}})
	if err != nil || len(got) != 2 || got["owner"] != "test" || got["qcs:other-observable"] != "unexpected" {
		t.Fatalf("broad tag suppression: %#v, %v", got, err)
	}
	if _, err := UserTags([]*sdk.Tag{{Key: pointer("owner"), Value: pointer("one")}, {Key: pointer("owner"), Value: pointer("two")}}); err == nil {
		t.Fatal("duplicate remote tag keys accepted")
	}
}
