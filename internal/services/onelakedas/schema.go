// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(_ context.Context, isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: map[string]superschema.Attribute{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Fabric Item ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Data Access Role.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"kind": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The kind of the Data Access Role.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleDataAccessRoleKindValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"decision_rules": superschema.SuperListNestedAttributeOf[decisionRuleModel]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "The array of permissions that make up the Data Access Role.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: decisionRuleSchema(),
			},
			"members": superschema.SuperSingleNestedAttributeOf[membersModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The members of the Data Access Role.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: membersSchema(),
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: dsTimeout,
			},
		},
	}
}

func decisionRuleSchema() superschema.Attributes {
	return superschema.Attributes{
		"effect": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The effect of the decision rule.",
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleEffectValues(), true)...),
				},
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"permission": superschema.SuperListNestedAttributeOf[permissionScopeModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The permission scopes for the decision rule.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Required: true,
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: permissionScopeSchema(),
		},
		"constraints": superschema.SuperSingleNestedAttributeOf[constraintsModel]{
			Common: &schemaR.SingleNestedAttribute{
				MarkdownDescription: "The constraints applied to the decision rule.",
			},
			Resource: &schemaR.SingleNestedAttribute{
				Optional: true,
			},
			DataSource: &schemaD.SingleNestedAttribute{
				Computed: true,
			},
			Attributes: constraintsSchema(),
		},
	}
}

func permissionScopeSchema() superschema.Attributes {
	return superschema.Attributes{
		"attribute_name": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The name of the attribute being evaluated for access permissions.",
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleAttributeNameValues(), true)...),
				},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"attribute_value_included_in": superschema.SuperListAttribute{
			Common: &schemaR.ListAttribute{
				MarkdownDescription: "The list of values for the attribute to define the scope and level of access.",
				CustomType: supertypes.ListTypeOf[types.String]{ListType: basetypes.ListType{ElemType: types.StringType}},
				ElementType: types.StringType,
			},
			Resource: &schemaR.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			DataSource: &schemaD.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func constraintsSchema() superschema.Attributes {
	return superschema.Attributes{
		"columns": superschema.SuperListNestedAttributeOf[columnConstraintModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The array of column constraints.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Optional: true,
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: columnConstraintSchema(),
		},
		"rows": superschema.SuperListNestedAttributeOf[rowConstraintModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The array of row constraints.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Optional: true,
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: rowConstraintSchema(),
		},
	}
}

func columnConstraintSchema() superschema.Attributes {
	return superschema.Attributes{
		"column_action": superschema.SuperListAttribute{
			Common: &schemaR.ListAttribute{
				MarkdownDescription: "The actions applied to the column names.",
				CustomType: supertypes.ListTypeOf[types.String]{ListType: basetypes.ListType{ElemType: types.StringType}},
				ElementType: types.StringType,
			},
			Resource: &schemaR.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			DataSource: &schemaD.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		"column_effect": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The effect given to the column names.",
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleColumnEffectValues(), true)...),
				},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"column_names": superschema.SuperListAttribute{
			Common: &schemaR.ListAttribute{
				MarkdownDescription: "An array of column names.",
				CustomType: supertypes.ListTypeOf[types.String]{ListType: basetypes.ListType{ElemType: types.StringType}},
				ElementType: types.StringType,
			},
			Resource: &schemaR.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			DataSource: &schemaD.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		"table_path": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "A relative file path specifying which table the constraint applies to.",
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}

func rowConstraintSchema() superschema.Attributes {
	return superschema.Attributes{
		"table_path": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "A relative file path specifying which table the row constraint applies to.",
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"value": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "A T-SQL expression used to evaluate which rows the role members can see.",
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}

func membersSchema() superschema.Attributes {
	return superschema.Attributes{
		"fabric_item_members": superschema.SuperListNestedAttributeOf[fabricItemMemberModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The list of Fabric item members.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Optional: true,
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: fabricItemMemberSchema(),
		},
		"microsoft_entra_members": superschema.SuperListNestedAttributeOf[microsoftEntraMemberModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The list of Microsoft Entra ID members.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Optional: true,
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: microsoftEntraMemberSchema(),
		},
	}
}

func fabricItemMemberSchema() superschema.Attributes {
	return superschema.Attributes{
		"item_access": superschema.SuperListAttribute{
			Common: &schemaR.ListAttribute{
				MarkdownDescription: "The access permissions for the Fabric item member.",
				CustomType: supertypes.ListTypeOf[types.String]{ListType: basetypes.ListType{ElemType: types.StringType}},
				ElementType: types.StringType,
			},
			Resource: &schemaR.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			DataSource: &schemaD.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		"source_path": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The path to the Fabric item.",
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}

func microsoftEntraMemberSchema() superschema.Attributes {
	return superschema.Attributes{
		"object_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The object ID.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"tenant_id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The tenant ID.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"object_type": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The type of Microsoft Entra ID object.",
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleObjectTypeValues(), true)...),
				},
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}
