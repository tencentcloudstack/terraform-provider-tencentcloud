package ags

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestManagedRefreshCompleteEmptyRemovesResource(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	client := deleteAckClient(t, func(action string, input map[string]any) (map[string]any, error) {
		calls++
		if action != "DescribeSandboxToolList" {
			t.Fatalf("refresh mutated resources: %s", action)
		}
		ids, ok := input["ToolIds"].([]any)
		if !ok || len(ids) != 1 || ids[0] != "sdt-offline" || input["Filters"] != nil {
			t.Fatal("refresh did not use the exact unfiltered state ID")
		}
		return deleteAckEmpty(), nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	resp := resource.ReadResponse{}
	r.Read(ctx, resource.ReadRequest{State: stateFor(t, prior)}, &resp)
	if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() || calls != 1 {
		t.Fatalf("complete empty refresh did not remove resource: %v, calls=%d", resp.Diagnostics, calls)
	}
}

func TestManagedRefreshErrorsAndIncompleteResultsKeepFullState(t *testing.T) {
	ctx := context.Background()
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remoteCode(), true)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"denied", "expired", "transport", "missing-list", "missing-total", "nonzero-total", "missing-receipt", "invalid-receipt", "wrong-id"} {
		t.Run(mode, func(t *testing.T) {
			client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
				if action != "DescribeSandboxToolList" {
					t.Fatalf("refresh mutated resources: %s", action)
				}
				switch mode {
				case "denied":
					return errorResult("UnauthorizedOperation"), nil
				case "expired":
					return errorResult("AuthFailure.TokenFailure"), nil
				case "transport":
					return nil, errors.New("read response lost")
				}
				body := deleteAckEmpty()
				switch mode {
				case "missing-list":
					delete(body, "SandboxToolSet")
				case "missing-total":
					delete(body, "TotalCount")
				case "nonzero-total":
					body["TotalCount"] = 1
				case "missing-receipt":
					delete(body, "RequestId")
				case "invalid-receipt":
					body["RequestId"] = "not-a-public-request-id"
				case "wrong-id":
					body = pointResult(testTool("sdt-another", "another", "ACTIVE"))
					body["RequestId"] = deleteAckRequestID
				}
				return body, nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			req := resource.ReadRequest{State: stateFor(t, prior)}
			resp := resource.ReadResponse{}
			r.Read(ctx, req, &resp)
			if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(req.State.Raw) {
				t.Fatalf("failed or incomplete refresh changed full state: %v", resp.Diagnostics)
			}
		})
	}
}

func TestManagedRefreshDoesNotChangeStrictLookupOrImport(t *testing.T) {
	client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
		if action != "DescribeSandboxToolList" {
			t.Fatalf("read/import mutated resources: %s", action)
		}
		return deleteAckEmpty(), nil
	})
	svc := fastService(client)
	if _, err := svc.ReadManaged(context.Background(), "sdt-owned"); !errors.Is(err, ErrConfirmedAbsent) {
		t.Fatalf("managed read did not normalize complete empty: %v", err)
	}
	if _, err := svc.Lookup(context.Background(), "sdt-owned"); !errors.Is(err, ErrAmbiguousAbsence) {
		t.Fatalf("strict lookup changed: %v", err)
	}
	r := &SandboxToolResource{testService: svc}
	resp := resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "sdt-owned"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("nonexistent import was accepted")
	}
}

func TestManagedRefreshRetainsProvisionalCreatingIdentity(t *testing.T) {
	ctx := context.Background()
	remote := remoteCode()
	remote.Status = pointer("CREATING")
	prior, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), remote, true)
	if err != nil {
		t.Fatal(err)
	}
	client := deleteAckClient(t, func(action string, _ map[string]any) (map[string]any, error) {
		if action != "DescribeSandboxToolList" {
			t.Fatalf("provisional read mutated resources: %s", action)
		}
		return deleteAckEmpty(), nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	req := resource.ReadRequest{State: stateFor(t, prior)}
	resp := resource.ReadResponse{}
	r.Read(ctx, req, &resp)
	if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(req.State.Raw) {
		t.Fatal("a provisional accepted ID was treated as an externally deleted managed resource")
	}
}
