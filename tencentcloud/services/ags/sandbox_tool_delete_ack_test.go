package ags

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const deleteAckRequestID = "f1538c42-c84f-4990-8c28-eaa113e28bc1"

// Capture actual public SDK HTTP serialization without using a network socket.
// Each response explicitly controls whether a public request receipt exists.
func deleteAckClient(t *testing.T, handle func(string, map[string]any) (map[string]any, error)) *sdk.Client {
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
		result, err := handle(req.Header.Get("X-TC-Action"), input)
		if err != nil {
			return nil, err
		}
		encoded, err := json.Marshal(map[string]any{"Response": result})
		if err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(encoded)))}, nil
	}))
	return client
}

func deleteAckEmpty() map[string]any {
	return map[string]any{"RequestId": deleteAckRequestID, "SandboxToolSet": []any{}, "TotalCount": 0}
}

func TestDeleteAckPolicyAcknowledgedEmptyCompletes(t *testing.T) {
	deletes, reads := 0, 0
	client := deleteAckClient(t, func(action string, input map[string]any) (map[string]any, error) {
		switch action {
		case "DeleteSandboxTool":
			deletes++
			if input["ToolId"] != "sdt-owned" {
				t.Fatal("wrong deletion identity")
			}
			return map[string]any{"RequestId": deleteAckRequestID}, nil
		case "DescribeSandboxToolList":
			reads++
			ids, ok := input["ToolIds"].([]any)
			if !ok || len(ids) != 1 || ids[0] != "sdt-owned" || input["Filters"] != nil {
				t.Fatal("confirmation must point query the exact deleted ID without status filters")
			}
			if reads == 1 {
				result := pointResult(testTool("sdt-owned", "owned", "DELETING"))
				result["RequestId"] = deleteAckRequestID
				return result, nil
			}
			return deleteAckEmpty(), nil
		default:
			t.Fatalf("unexpected action %s", action)
			return nil, nil
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := fastService(client).DeleteAndWait(ctx, "sdt-owned"); err != nil || deletes != 1 || reads != 2 {
		t.Fatalf("acknowledged deletion failed: err=%v deletes=%d reads=%d", err, deletes, reads)
	}
}

func TestDeleteAckPolicyRejectsUnprovenAndMalformedResults(t *testing.T) {
	for _, mode := range []string{"lost-delete", "delete-denied", "missing-ack", "invalid-ack", "read-denied", "read-transport", "missing-total", "nonzero-total", "missing-list", "missing-read-receipt", "wrong-id"} {
		t.Run(mode, func(t *testing.T) {
			deletes := 0
			client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
				if action == "DeleteSandboxTool" {
					deletes++
					if mode == "lost-delete" {
						return nil, errors.New("response lost")
					}
					if mode == "delete-denied" {
						return errorResult("UnauthorizedOperation"), nil
					}
					if mode == "missing-ack" {
						return map[string]any{}, nil
					}
					if mode == "invalid-ack" {
						return map[string]any{"RequestId": "not-a-cloud-request-id"}, nil
					}
					return map[string]any{"RequestId": deleteAckRequestID}, nil
				}
				if action != "DescribeSandboxToolList" {
					t.Fatalf("unexpected action %s", action)
				}
				if mode == "read-denied" {
					return errorResult("UnauthorizedOperation"), nil
				}
				if mode == "read-transport" {
					return nil, errors.New("read response lost")
				}
				body := deleteAckEmpty()
				switch mode {
				case "missing-total":
					delete(body, "TotalCount")
				case "nonzero-total":
					body["TotalCount"] = 1
				case "missing-list":
					delete(body, "SandboxToolSet")
				case "missing-read-receipt":
					delete(body, "RequestId")
				case "wrong-id":
					body = pointResult(testTool("sdt-other", "other", "DELETING"))
					body["RequestId"] = deleteAckRequestID
				}
				return body, nil
			})
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := fastService(client).DeleteAndWait(ctx, "sdt-owned"); err == nil || deletes != 1 {
				t.Fatalf("unproven deletion accepted or replayed: err=%v deletes=%d", err, deletes)
			}
		})
	}
}

func TestDeleteAckPolicyDoesNotChangeOrdinaryReads(t *testing.T) {
	client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
		if action != "DescribeSandboxToolList" {
			t.Fatalf("read performed mutation %s", action)
		}
		return deleteAckEmpty(), nil
	})
	svc := fastService(client)
	if _, err := svc.Lookup(context.Background(), "sdt-owned"); !errors.Is(err, ErrAmbiguousAbsence) {
		t.Fatalf("ordinary empty read changed classification: %v", err)
	}
	if err := svc.WaitDeleted(context.Background(), "sdt-owned"); !errors.Is(err, ErrAmbiguousAbsence) {
		t.Fatalf("wait without a Delete receipt accepted empty: %v", err)
	}
}

func TestDeleteAckPolicyReceiptDoesNotEscapeOperation(t *testing.T) {
	client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
		if action == "DeleteSandboxTool" {
			return map[string]any{"RequestId": deleteAckRequestID}, nil
		}
		if action != "DescribeSandboxToolList" {
			t.Fatalf("unexpected action %s", action)
		}
		return deleteAckEmpty(), nil
	})
	svc := fastService(client)
	if err := svc.DeleteAndWait(context.Background(), "sdt-owned"); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"sdt-owned", "sdt-another"} {
		if _, err := svc.Lookup(context.Background(), id); !errors.Is(err, ErrAmbiguousAbsence) {
			t.Fatalf("receipt leaked into ordinary Lookup for %s: %v", id, err)
		}
		if err := svc.WaitDeleted(context.Background(), id); !errors.Is(err, ErrAmbiguousAbsence) {
			t.Fatalf("receipt leaked into a later unacknowledged wait for %s: %v", id, err)
		}
	}
}

func TestDeleteAckPolicyCancellationAndDeadline(t *testing.T) {
	for _, mode := range []string{"already-cancelled", "cancel-after-read", "deadline"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			if mode == "already-cancelled" {
				cancel()
			}
			deletes, reads := 0, 0
			client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
				if action == "DeleteSandboxTool" {
					deletes++
					return map[string]any{"RequestId": deleteAckRequestID}, nil
				}
				if action != "DescribeSandboxToolList" {
					t.Fatalf("unexpected action %s", action)
				}
				reads++
				if mode == "cancel-after-read" {
					cancel()
				}
				result := pointResult(testTool("sdt-owned", "owned", "DELETING"))
				result["RequestId"] = deleteAckRequestID
				return result, nil
			})
			err := fastService(client).DeleteAndWait(ctx, "sdt-owned")
			if mode == "already-cancelled" {
				if !errors.Is(err, context.Canceled) || deletes != 0 || reads != 0 {
					t.Fatalf("cancelled deletion issued requests: %v %d/%d", err, deletes, reads)
				}
			} else if err == nil || deletes != 1 || reads < 1 || mode == "cancel-after-read" && reads != 1 {
				t.Fatalf("poll did not honor operation boundary: %v %d/%d", err, deletes, reads)
			}
		})
	}
}

func TestDeleteAckPolicyResourceOnlyClearsAcknowledgedState(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	for _, acknowledged := range []bool{true, false} {
		client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
			if action == "DeleteSandboxTool" {
				if !acknowledged {
					return nil, errors.New("Delete acknowledgement lost")
				}
				return map[string]any{"RequestId": deleteAckRequestID}, nil
			}
			if action != "DescribeSandboxToolList" {
				t.Fatalf("implicit instance action %s", action)
			}
			return deleteAckEmpty(), nil
		})
		r := &SandboxToolResource{testService: fastService(client)}
		req := resource.DeleteRequest{State: stateFor(t, prior)}
		resp := resource.DeleteResponse{}
		r.Delete(ctx, req, &resp)
		if acknowledged {
			if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() {
				t.Fatalf("acknowledged delete retained state: %v", resp.Diagnostics)
			}
		} else if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(req.State.Raw) {
			t.Fatal("lost acknowledgement discarded state")
		}
	}
}
