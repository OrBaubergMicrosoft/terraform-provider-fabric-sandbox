// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedas_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDataAccessSecurityResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id": "00000000-0000-0000-0000-000000000000",
					"name":    "TestRole",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{
									"attribute_name":              "Path",
									"attribute_value_included_in": []string{"Tables/Test"},
								},
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - invalid UUID - workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      "00000000-0000-0000-0000-000000000000",
					"name":         "TestRole",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{
									"attribute_name":              "Path",
									"attribute_value_included_in": []string{"Tables/Test"},
								},
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - item_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"item_id":      "invalid uuid",
					"name":         "TestRole",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{
									"attribute_name":              "Path",
									"attribute_value_included_in": []string{"Tables/Test"},
								},
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("WorkspaceID/ItemID/RoleName"),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("WorkspaceID/ItemID/RoleName"),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "invalid/00000000-0000-0000-0000-000000000000/TestRole",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "00000000-0000-0000-0000-000000000000/invalid/TestRole",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := newRandomDataAccessRoleBase()

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRole()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRole(entity)
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRole()

	objectID := *entity.Members.MicrosoftEntraMembers[0].ObjectID
	tenantID := *entity.Members.MicrosoftEntraMembers[0].TenantID

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         "TestRole",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{
									"attribute_name":              "Path",
									"attribute_value_included_in": []string{"Tables/TestTable"},
								},
							},
						},
					},
					"members": map[string]any{
						"microsoft_entra_members": []map[string]any{
							{
								"object_id":   objectID,
								"tenant_id":   tenantID,
								"object_type": "User",
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", "TestRole"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "kind", string(fabcore.DataAccessRoleKindPolicy)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "decision_rules.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "1"),
			),
		},
		// Update (same config since it's upsert)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         "TestRole",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{
									"attribute_name":              "Path",
									"attribute_value_included_in": []string{"Tables/TestTable"},
								},
							},
						},
					},
					"members": map[string]any{
						"microsoft_entra_members": []map[string]any{
							{
								"object_id":   objectID,
								"tenant_id":   tenantID,
								"object_type": "User",
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", "TestRole"),
			),
		},
	}))
}
