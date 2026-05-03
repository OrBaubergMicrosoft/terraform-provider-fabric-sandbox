// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabanomalydetector "github.com/microsoft/fabric-sdk-go/fabric/anomalydetector"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsAnomalyDetector struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsAnomalyDetector) CreateDefinition(data fabanomalydetector.CreateAnomalyDetectorRequest) *fabanomalydetector.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsAnomalyDetector) TransformDefinition(entity *fabanomalydetector.Definition) fabanomalydetector.ItemsClientGetAnomalyDetectorDefinitionResponse {
	return fabanomalydetector.ItemsClientGetAnomalyDetectorDefinitionResponse{
		DefinitionResponse: fabanomalydetector.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsAnomalyDetector) UpdateDefinition(_ *fabanomalydetector.Definition, data fabanomalydetector.UpdateAnomalyDetectorDefinitionRequest) *fabanomalydetector.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsAnomalyDetector) CreateWithParentID(parentID string, data fabanomalydetector.CreateAnomalyDetectorRequest) fabanomalydetector.AnomalyDetector {
	entity := NewRandomAnomalyDetectorWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsAnomalyDetector) Filter(entities []fabanomalydetector.AnomalyDetector, parentID string) []fabanomalydetector.AnomalyDetector {
	ret := make([]fabanomalydetector.AnomalyDetector, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsAnomalyDetector) GetID(entity fabanomalydetector.AnomalyDetector) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsAnomalyDetector) TransformCreate(entity fabanomalydetector.AnomalyDetector) fabanomalydetector.ItemsClientCreateAnomalyDetectorResponse {
	return fabanomalydetector.ItemsClientCreateAnomalyDetectorResponse{
		AnomalyDetector: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsAnomalyDetector) TransformGet(entity fabanomalydetector.AnomalyDetector) fabanomalydetector.ItemsClientGetAnomalyDetectorResponse {
	return fabanomalydetector.ItemsClientGetAnomalyDetectorResponse{
		AnomalyDetector: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsAnomalyDetector) TransformList(entities []fabanomalydetector.AnomalyDetector) fabanomalydetector.ItemsClientListAnomalyDetectorsResponse {
	return fabanomalydetector.ItemsClientListAnomalyDetectorsResponse{
		AnomalyDetectors: fabanomalydetector.AnomalyDetectors{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsAnomalyDetector) TransformUpdate(entity fabanomalydetector.AnomalyDetector) fabanomalydetector.ItemsClientUpdateAnomalyDetectorResponse {
	return fabanomalydetector.ItemsClientUpdateAnomalyDetectorResponse{
		AnomalyDetector: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsAnomalyDetector) Update(base fabanomalydetector.AnomalyDetector, data fabanomalydetector.UpdateAnomalyDetectorRequest) fabanomalydetector.AnomalyDetector {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsAnomalyDetector) Validate(newEntity fabanomalydetector.AnomalyDetector, existing []fabanomalydetector.AnomalyDetector) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureAnomalyDetector(server *fakeServer) fabanomalydetector.AnomalyDetector {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabanomalydetector.AnomalyDetector,
			fabanomalydetector.ItemsClientGetAnomalyDetectorResponse,
			fabanomalydetector.ItemsClientUpdateAnomalyDetectorResponse,
			fabanomalydetector.ItemsClientCreateAnomalyDetectorResponse,
			fabanomalydetector.ItemsClientListAnomalyDetectorsResponse,
			fabanomalydetector.CreateAnomalyDetectorRequest,
			fabanomalydetector.UpdateAnomalyDetectorRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			fabanomalydetector.Definition,
			fabanomalydetector.CreateAnomalyDetectorRequest,
			fabanomalydetector.UpdateAnomalyDetectorDefinitionRequest,
			fabanomalydetector.ItemsClientGetAnomalyDetectorDefinitionResponse,
			fabanomalydetector.ItemsClientUpdateAnomalyDetectorDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsAnomalyDetector{}
	var definitionOperations concreteDefinitionOperations = &operationsAnomalyDetector{}
	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.AnomalyDetector.ItemsServer.GetAnomalyDetector,
		&server.ServerFactory.AnomalyDetector.ItemsServer.UpdateAnomalyDetector,
		&server.ServerFactory.AnomalyDetector.ItemsServer.BeginCreateAnomalyDetector,
		&server.ServerFactory.AnomalyDetector.ItemsServer.NewListAnomalyDetectorsPager,
		&server.ServerFactory.AnomalyDetector.ItemsServer.DeleteAnomalyDetector)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.AnomalyDetector.ItemsServer.BeginCreateAnomalyDetector,
		&server.ServerFactory.AnomalyDetector.ItemsServer.BeginGetAnomalyDetectorDefinition,
		&server.ServerFactory.AnomalyDetector.ItemsServer.BeginUpdateAnomalyDetectorDefinition)

	return fabanomalydetector.AnomalyDetector{}
}

func NewRandomAnomalyDetector() fabanomalydetector.AnomalyDetector {
	return fabanomalydetector.AnomalyDetector{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabanomalydetector.ItemTypeAnomalyDetector),
	}
}

func NewRandomAnomalyDetectorWithWorkspace(workspaceID string) fabanomalydetector.AnomalyDetector {
	result := NewRandomAnomalyDetector()
	result.WorkspaceID = &workspaceID

	return result
}
func NewRandomAnomalyDetectorDefinition() fabanomalydetector.Definition {
	defPart := fabanomalydetector.DefinitionPart{
		PayloadType: to.Ptr(fabanomalydetector.PayloadTypeInlineBase64),
		Path:        to.Ptr("Configurations.json"),
		Payload: to.Ptr(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==", // {"content":"Hello World"} in base64
		),
	}

	defParts := make([]fabanomalydetector.DefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return fabanomalydetector.Definition{
		Format: to.Ptr(fabanomalydetector.DefinitionFormatAnomalyDetectorV1),
		Parts:  defParts,
	}
}
