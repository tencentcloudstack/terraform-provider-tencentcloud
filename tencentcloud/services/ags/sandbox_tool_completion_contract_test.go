package ags

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func completionPrior(t *testing.T) types.Object {
	t.Helper()
	tool := remoteCode()
	tool.Description = pointer("old")
	tool.Tags = append(tool.Tags, &sdk.Tag{Key: pointer("revision"), Value: pointer("old")})
	state, err := flattenTool(context.Background(), types.ObjectNull(inputAttributeTypes()), tool, true)
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func completionUpdatePlan(t *testing.T, prior types.Object) types.Object {
	t.Helper()
	fields := prior.Attributes()
	fields["description"] = types.StringValue("new")
	fields["tags"] = types.MapValueMust(types.StringType, map[string]attr.Value{"revision": types.StringValue("new")})
	result, ds := types.ObjectValue(inputAttributeTypes(), fields)
	if ds.HasError() {
		t.Fatal(ds)
	}
	return result
}

func TestCT002ExpiredIdempotencyEvidenceNeverTriggersAnotherMutation(t *testing.T) {
	// These virtual gateway windows intentionally expire while the first
	// response is unavailable. Replaying after expiry would no longer return
	// a cached success. The Provider must stop, retaining name/region recovery
	// context, rather than extending the gateway's finite guarantee itself.
	for _, window := range []string{"processing-lock-ttl", "success-record-ttl"} {
		t.Run(window, func(t *testing.T) {
			virtualNow := 0
			expiresAt := 1
			requests, accepted := 0, 0
			var firstToken string
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				if action != "CreateSandboxTool" {
					t.Fatalf("idless create performed automatic adoption or cleanup: %s", action)
				}
				requests++
				token, ok := input["ClientToken"].(string)
				if !ok || token == "" {
					t.Fatal("missing idempotency token on the actual SDK request")
				}
				if requests == 1 {
					accepted++
					firstToken = token
					virtualNow = expiresAt + 1
					return nil, io.ErrUnexpectedEOF
				}
				if virtualNow <= expiresAt {
					t.Fatal("fixture did not expire its gateway evidence")
				}
				if token != firstToken {
					t.Error("Provider replaced its token after an ambiguous acceptance")
				}
				// A Tool name remains unique independently of the expired token.
				// This response must never be reached automatically.
				return errorResult("ResourceInUse.SandboxToolName"), nil
			})
			config := modelFixture(t, nil)
			r := &SandboxToolResource{testService: fastService(client)}
			response := resource.CreateResponse{State: stateFor(t, types.ObjectNull(inputAttributeTypes()))}
			r.Create(context.Background(), resource.CreateRequest{Config: configFor(t, config), Plan: planFor(t, computedPlan(t, config))}, &response)
			if !response.Diagnostics.HasError() || !response.State.Raw.IsNull() || requests != 1 || accepted != 1 || virtualNow <= expiresAt {
				t.Fatalf("expired %s was treated as permanent dedupe: requests=%d accepted=%d diagnostics=%v", window, requests, accepted, response.Diagnostics)
			}
			if !strings.Contains(fmt.Sprint(response.Diagnostics), "owned-code") || !strings.Contains(fmt.Sprint(response.Diagnostics), "import") {
				t.Fatal("idless failure lost its explicit owner-reconciliation instruction")
			}
		})
	}
}

func TestCT002ProcessExitBeforeCreateResponseRequiresOwnerImport(t *testing.T) {
	const childFlag = "AGS_CT_CREATE_CRASH_CHILD"
	if os.Getenv(childFlag) == "1" {
		client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
			if action != "CreateSandboxTool" || input["ToolName"] != "owned-code" || input["ClientToken"] == "" {
				t.Fatal("unexpected pre-crash request")
			}
			// This is after actual SDK request encoding and gateway acceptance,
			// but before any response/ID reaches Resource.Create or Terraform.
			if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"accepted": true, "tool_id": "sdt-crashowned", "token_present": input["ClientToken"] != nil}); err != nil {
				t.Fatal(err)
			}
			var release [1]byte
			if _, err := io.ReadFull(os.Stdin, release[:]); err != nil {
				return nil, err
			}
			t.Fatal("crash helper unexpectedly received a release instead of exiting")
			return nil, errors.New("unreachable response")
		})
		config := modelFixture(t, nil)
		r := &SandboxToolResource{testService: fastService(client)}
		response := resource.CreateResponse{}
		r.Create(context.Background(), resource.CreateRequest{Config: configFor(t, config), Plan: planFor(t, computedPlan(t, config))}, &response)
		t.Fatal("Create returned before the parent terminated the process")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCT002ProcessExitBeforeCreateResponseRequiresOwnerImport$", "-test.v")
	child.Env = append(os.Environ(), childFlag+"=1")
	stdout, err := child.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := child.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer stdin.Close()
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	accepted := make(chan map[string]any, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var event map[string]any
			if json.Unmarshal(scanner.Bytes(), &event) == nil && event["accepted"] == true {
				accepted <- event
				return
			}
		}
	}()
	var receipt map[string]any
	select {
	case receipt = <-accepted:
	case <-ctx.Done():
		_ = child.Process.Kill()
		_ = child.Wait()
		t.Fatal("child never reached the pre-response acceptance boundary")
	}
	if err := child.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	if err := child.Wait(); err == nil {
		t.Fatal("pre-response crash did not terminate the SDK caller")
	}
	if receipt["token_present"] != true || receipt["tool_id"] != "sdt-crashowned" {
		t.Fatal("owner did not retain independent acceptance evidence")
	}
	reads := 0
	client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
		if action != "DescribeSandboxToolList" {
			t.Fatalf("restart recovery wrote or silently adopted by name: %s", action)
		}
		reads++
		ids, ok := input["ToolIds"].([]any)
		if !ok || len(ids) != 1 || ids[0] != receipt["tool_id"] {
			t.Fatal("recovery did not use the independently reconciled exact identity")
		}
		return pointResult(testTool("sdt-crashowned", "owned-code", "ACTIVE")), nil
	})
	r := &SandboxToolResource{testService: fastService(client)}
	imported := resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: receipt["tool_id"].(string)}, &imported)
	if imported.Diagnostics.HasError() || reads != 1 {
		t.Fatalf("explicit owner import failed: %v", imported.Diagnostics)
	}
	var actual types.Object
	if ds := imported.State.Get(context.Background(), &actual); ds.HasError() {
		t.Fatal(ds)
	}
	if value(stringInput(actual.Attributes(), "id")) != receipt["tool_id"] {
		t.Fatal("explicit import did not restore the accepted identity")
	}
}

func TestCT003ServiceErrorMatrixKeepsExactResourceState(t *testing.T) {
	prior := completionPrior(t)
	for _, code := range []string{"InternalError", "InternalError.DatabaseError", "ServiceUnavailable", "RequestLimitExceeded.GlobalRegionUinLimitExceeded"} {
		t.Run(code, func(t *testing.T) {
			calls := 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				calls++
				if action != "DescribeSandboxToolList" {
					t.Fatalf("transient read failure triggered mutation %s", action)
				}
				return errorResult(code), nil
			})
			r := &SandboxToolResource{testService: fastService(client)}
			request := resource.ReadRequest{State: stateFor(t, prior)}
			response := resource.ReadResponse{}
			r.Read(context.Background(), request, &response)
			if !response.Diagnostics.HasError() || !response.State.Raw.Equal(request.State.Raw) || calls != 3 {
				t.Fatalf("read failure was not bounded/state-preserving: calls=%d diagnostics=%v", calls, response.Diagnostics)
			}
		})
	}
	t.Run("cancellation-during-backoff", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := 0
		client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
			calls++
			if action != "DescribeSandboxToolList" {
				t.Fatal("cancelled Read sent a mutation")
			}
			cancel()
			return errorResult("InternalError"), nil
		})
		r := &SandboxToolResource{testService: fastService(client)}
		request := resource.ReadRequest{State: stateFor(t, prior)}
		response := resource.ReadResponse{}
		r.Read(ctx, request, &response)
		if !response.Diagnostics.HasError() || calls != 1 || !response.State.Raw.Equal(request.State.Raw) {
			t.Fatal("cancellation retried or lost state", response.Diagnostics)
		}
	})
}

func TestCT005PartialUpdateAndInterruptedReadPreserveTrustworthyState(t *testing.T) {
	for _, fault := range []string{"partial-then-denied", "partial-custom-env", "response-lost-partial-readable", "read-lost-before-any-observation"} {
		t.Run(fault, func(t *testing.T) {
			ctx := context.Background()
			prior := completionPrior(t)
			if fault == "partial-custom-env" {
				tool := remoteCode()
				tool.ToolType, tool.Description = pointer("custom"), pointer("old")
				tool.Tags = []*sdk.Tag{{Key: pointer("revision"), Value: pointer("old")}}
				tool.CustomConfiguration = &sdk.CustomConfigurationDetail{
					Image: pointer("offline.invalid/ags:1"), ImageRegistryType: pointer("enterprise"),
					Resources: &sdk.ResourceConfiguration{CPU: pointer("1"), Memory: pointer("1Gi")},
					Env:       []*sdk.EnvVar{{Name: pointer("REVISION"), Value: pointer("old")}},
				}
				var err error
				prior, err = flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), tool, true)
				if err != nil {
					t.Fatal(err)
				}
			}
			plan := completionUpdatePlan(t, prior)
			if fault == "partial-custom-env" {
				fields := plan.Attributes()
				block := fields["custom_configuration"].(types.List)
				custom := block.Elements()[0].(types.Object)
				attributes := custom.Attributes()
				attributes["env"] = types.MapValueMust(types.StringType, map[string]attr.Value{"REVISION": types.StringValue("new")})
				fields["custom_configuration"] = types.ListValueMust(block.ElementType(ctx), []attr.Value{types.ObjectValueMust(custom.AttributeTypes(ctx), attributes)})
				plan = types.ObjectValueMust(inputAttributeTypes(), fields)
			}
			updates, reads := 0, 0
			recovering := false
			backendDescription, backendRevision := "old", "old"
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				switch action {
				case "UpdateSandboxTool":
					updates++
					if updates > 1 && !recovering {
						t.Fatal("Provider blindly retried a partially accepted mutation")
					}
					backendRevision = "new"
					if fault == "partial-custom-env" && !recovering {
						custom, ok := input["CustomConfiguration"].(map[string]any)
						if !ok {
							t.Fatal("custom partial-update fixture did not receive a real env change")
						}
						env, ok := custom["Env"].([]any)
						if !ok || len(env) != 1 || env[0].(map[string]any)["Value"] != "new" {
							t.Fatal("custom environment did not reach the actual SDK body")
						}
					}
					if recovering {
						if _, present := input["Tags"]; present || len(input) != 2 || input["Description"] != "new" {
							t.Fatal("explicit recovery resent an already accepted field")
						}
						backendDescription = "new"
					}
					if fault == "response-lost-partial-readable" && !recovering {
						return nil, io.ErrUnexpectedEOF
					}
					return map[string]any{}, nil
				case "DescribeSandboxToolList":
					reads++
					if !recovering && (fault == "read-lost-before-any-observation" || (fault == "partial-then-denied" || fault == "partial-custom-env") && reads > 1) {
						return errorResult("UnauthorizedOperation"), nil
					}
					tool := testTool("sdt-offline", "owned-code", "ACTIVE")
					tool["Description"] = backendDescription
					tool["Tags"] = []any{map[string]any{"Key": "revision", "Value": backendRevision}}
					if fault == "partial-custom-env" {
						tool["ToolType"] = "custom"
						tool["CustomConfiguration"] = map[string]any{
							"Image": "offline.invalid/ags:1", "ImageRegistryType": "enterprise",
							"Resources": map[string]any{"CPU": "1", "Memory": "1Gi"},
							"Env":       []any{map[string]any{"Name": "REVISION", "Value": backendRevision}},
						}
					}
					return pointResult(tool), nil
				default:
					t.Fatalf("update recovery invoked forbidden Action %s", action)
					return nil, errors.New("forbidden action")
				}
			})
			r := &SandboxToolResource{testService: fastService(client)}
			response := resource.UpdateResponse{}
			r.Update(ctx, resource.UpdateRequest{State: stateFor(t, prior), Plan: planFor(t, plan)}, &response)
			if !response.Diagnostics.HasError() || updates != 1 {
				t.Fatal("partial/unobserved mutation was reported as completed", response.Diagnostics)
			}
			if fault == "read-lost-before-any-observation" {
				if !response.State.Raw.Equal(stateFor(t, prior).Raw) {
					t.Fatal("unobserved remote changes were fabricated in state")
				}
				recovering = true
				refreshed := resource.ReadResponse{}
				r.Read(ctx, resource.ReadRequest{State: response.State}, &refreshed)
				if refreshed.Diagnostics.HasError() {
					t.Fatal(refreshed.Diagnostics)
				}
				response.State = refreshed.State
			}
			var actual types.Object
			if ds := response.State.Get(ctx, &actual); ds.HasError() {
				t.Fatal(ds)
			}
			actualTags, err := stringMap(actual.Attributes()["tags"])
			if err != nil || value(stringInput(actual.Attributes(), "id")) != "sdt-offline" || value(stringInput(actual.Attributes(), "description")) != "old" || actualTags["revision"] != "new" {
				t.Fatal("actual mixed state was lost or replaced by the requested state")
			}
			patch, changed, err := updateRequest(actual, plan)
			if err != nil || !changed || patch.Tags != nil || patch.CustomConfiguration != nil || value(patch.Description) != "new" {
				t.Fatal("remaining difference did not preserve accepted partial success")
			}
			if fault == "partial-custom-env" {
				custom, err := oneBlock(actual.Attributes(), "custom_configuration", true)
				if err != nil {
					t.Fatal(err)
				}
				env, err := stringMap(custom["env"])
				if err != nil || env["REVISION"] != "new" {
					t.Fatal("accepted env update was lost from trustworthy partial state")
				}
			}
			recovering = true
			completed := resource.UpdateResponse{}
			r.Update(ctx, resource.UpdateRequest{State: response.State, Plan: planFor(t, plan)}, &completed)
			if completed.Diagnostics.HasError() || updates != 2 {
				t.Fatalf("explicit recovery failed: updates=%d diagnostics=%v", updates, completed.Diagnostics)
			}
		})
	}
}

func TestCT005DefaultsAreAcceptedOnlyForUndeclaredComputedInputs(t *testing.T) {
	custom := map[string]any{
		"image": "offline.invalid/ags:1", "image_registry_type": "enterprise",
		"resources": []map[string]any{{"cpu": "1", "memory": "1Gi"}},
		"command":   types.ListUnknown(types.StringType), "args": types.ListUnknown(types.StringType), "env": types.MapUnknown(types.StringType),
	}
	plan := computedPlan(t, modelFixture(t, map[string]any{"tool_type": "custom", "custom_configuration": []map[string]any{custom}}))
	tool := remoteCode()
	tool.ToolType = pointer("custom")
	tool.CustomConfiguration = &sdk.CustomConfigurationDetail{
		Image: pointer("offline.invalid/ags:1"), ImageRegistryType: pointer("enterprise"),
		Resources: &sdk.ResourceConfiguration{CPU: pointer("1000m"), Memory: pointer("1024Mi")},
		Command:   []*string{pointer("/image/default")}, Args: []*string{pointer("--default")}, Env: []*sdk.EnvVar{{Name: pointer("DEFAULT"), Value: pointer("1")}},
	}
	if ok, err := matchesPlan(context.Background(), plan, tool); err != nil || !ok {
		t.Fatalf("omitted computed input did not accept actual defaults: %v", err)
	}
	custom["env"] = map[string]string{}
	explicit := computedPlan(t, modelFixture(t, map[string]any{"tool_type": "custom", "custom_configuration": []map[string]any{custom}}))
	if ok, err := matchesPlan(context.Background(), explicit, tool); err != nil || ok {
		t.Fatal("an explicit reset was incorrectly suppressed as a remote default", err)
	}
}

func TestCT006CreatingTransitionsKeepTerminalIdentityAndAllowDestroy(t *testing.T) {
	for _, terminal := range []string{"ACTIVE", "FAILED", "UNKNOWN_FUTURE_STATUS"} {
		t.Run(terminal, func(t *testing.T) {
			creates, reads, deletes := 0, 0, 0
			deleting := false
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				switch action {
				case "CreateSandboxTool":
					creates++
					return map[string]any{"ToolId": "sdt-offline"}, nil
				case "DeleteSandboxTool":
					deletes++
					deleting = true
					return map[string]any{}, nil
				case "DescribeSandboxToolList":
					reads++
					if deleting {
						return errorResult("ResourceNotFound.SandboxTool"), nil
					}
					status := terminal
					if reads == 1 {
						status = "CREATING"
					}
					tool := testTool("sdt-offline", "owned-code", status)
					tool["StatusReason"] = "deterministic-" + terminal
					return pointResult(tool), nil
				default:
					t.Fatal("state machine invoked an instance or unrelated API", action)
					return nil, errors.New("forbidden action")
				}
			})
			r := &SandboxToolResource{testService: fastService(client)}
			config := modelFixture(t, nil)
			created := resource.CreateResponse{}
			r.Create(context.Background(), resource.CreateRequest{Config: configFor(t, config), Plan: planFor(t, computedPlan(t, config))}, &created)
			if creates != 1 || reads != 2 || created.Diagnostics.HasError() != (terminal != "ACTIVE") {
				t.Fatalf("wrong transition outcome: creates=%d reads=%d diagnostics=%v", creates, reads, created.Diagnostics)
			}
			var state types.Object
			if ds := created.State.Get(context.Background(), &state); ds.HasError() {
				t.Fatal(ds)
			}
			if value(stringInput(state.Attributes(), "id")) != "sdt-offline" {
				t.Fatal("transition failure lost its accepted identity")
			}
			if terminal == "FAILED" {
				refreshed := resource.ReadResponse{}
				r.Read(context.Background(), resource.ReadRequest{State: created.State}, &refreshed)
				if refreshed.Diagnostics.HasError() {
					t.Fatal("FAILED Read prevents a subsequent destroy", refreshed.Diagnostics)
				}
				created.State = refreshed.State
			}
			deleted := resource.DeleteResponse{}
			r.Delete(context.Background(), resource.DeleteRequest{State: created.State}, &deleted)
			if deleted.Diagnostics.HasError() || !deleted.State.Raw.IsNull() || deletes != 1 {
				t.Fatal("terminal-state cleanup failed", deleted.Diagnostics)
			}
		})
	}
}

func TestCT006DeletingAndPausedDependencyTransitions(t *testing.T) {
	for _, runtimeState := range []string{"RUNNING", "PAUSED"} {
		t.Run(runtimeState, func(t *testing.T) {
			prior := completionPrior(t)
			blocked, deleting := true, false
			deletes, reads := 0, 0
			client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
				switch action {
				case "DeleteSandboxTool":
					deletes++
					if blocked {
						return errorResult("UnsupportedOperation.SandboxToolHasActiveInstances"), nil
					}
					deleting = true
					return map[string]any{}, nil
				case "DescribeSandboxToolList":
					reads++
					if !deleting {
						t.Fatalf("Provider queried around the %s dependency denial", runtimeState)
					}
					if reads == 1 {
						return pointResult(testTool("sdt-offline", "owned-code", "DELETING")), nil
					}
					return errorResult("ResourceNotFound.SandboxTool"), nil
				default:
					t.Fatalf("Provider implicitly changed its %s instance: %s", runtimeState, action)
					return nil, errors.New("forbidden action")
				}
			})
			r := &SandboxToolResource{testService: fastService(client)}
			request := resource.DeleteRequest{State: stateFor(t, prior)}
			response := resource.DeleteResponse{}
			r.Delete(context.Background(), request, &response)
			if !response.Diagnostics.HasError() || deletes != 1 || reads != 0 || !response.State.Raw.Equal(request.State.Raw) {
				t.Fatal("dependency failure lost state or performed extra operations", response.Diagnostics)
			}
			// Only the fixture's separate runtime owner releases the instance.
			blocked = false
			response = resource.DeleteResponse{}
			r.Delete(context.Background(), request, &response)
			if response.Diagnostics.HasError() || deletes != 2 || reads != 2 || !response.State.Raw.IsNull() {
				t.Fatal("DELETING -> reliable NotFound did not finish", response.Diagnostics)
			}
		})
	}
	t.Run("cancel-after-trustworthy-read", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		reads := 0
		client := testClient(t, func(action string, input map[string]any) (map[string]any, error) {
			reads++
			if action != "DescribeSandboxToolList" {
				t.Fatal("wait cancellation mutated an instance or Tool")
			}
			return pointResult(testTool("sdt-offline", "owned-code", "ACTIVE")), nil
		})
		tool, err := fastService(client).WaitReady(ctx, "sdt-offline", false, func(*sdk.SandboxTool) (bool, error) {
			cancel()
			return false, nil
		})
		if !errors.Is(err, context.Canceled) || tool == nil || value(tool.ToolId) != "sdt-offline" || reads != 1 {
			t.Fatal("cancel lost last trustworthy read or sent background requests", err)
		}
	})
}
