// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package httproutes

import (
	"context"
	"errors"
	"net/http"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	manager "github.com/project-radius/radius/pkg/armrpc/asyncoperation/statusmanager"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/project-radius/radius/pkg/armrpc/servicecontext"
	"github.com/project-radius/radius/pkg/corerp/datamodel"
	"github.com/project-radius/radius/pkg/corerp/datamodel/converter"
	"github.com/project-radius/radius/pkg/radrp/rest"
	"github.com/project-radius/radius/pkg/ucp/store"
)

var _ ctrl.Controller = (*CreateOrUpdateHTTPRoute)(nil)

// CreateOrUpdateHTTPRoute is the controller implementation to create or update HTTPRoute resource.
type CreateOrUpdateHTTPRoute struct {
	ctrl.BaseController
}

// NewCreateOrUpdateTTPRoute creates a new CreateOrUpdateHTTPRoute.
func NewCreateOrUpdateHTTPRoute(ds store.StorageClient, sm manager.StatusManager) (ctrl.Controller, error) {
	return &CreateOrUpdateHTTPRoute{ctrl.NewBaseController(ds, sm)}, nil
}

// Run executes CreateOrUpdateHTTPRoute operation.
func (e *CreateOrUpdateHTTPRoute) Run(ctx context.Context, req *http.Request) (rest.Response, error) {
	serviceCtx := servicecontext.ARMRequestContextFromContext(ctx)
	newResource, err := e.Validate(ctx, req, serviceCtx.APIVersion)
	if err != nil {
		return nil, err
	}

	existingResource := &datamodel.HTTPRoute{}
	etag, err := e.GetResource(ctx, serviceCtx.ResourceID.String(), existingResource)
	if req.Method == http.MethodPatch && errors.Is(&store.ErrNotFound{}, err) {
		return rest.NewNotFoundResponse(serviceCtx.ResourceID), nil
	}
	if err != nil && !errors.Is(&store.ErrNotFound{}, err) {
		return nil, err
	}

	err = ctrl.ValidateETag(*serviceCtx, etag)
	if err != nil {
		return rest.NewPreconditionFailedResponse(serviceCtx.ResourceID.String(), err.Error()), nil
	}

	UpdateExistingResourceData(ctx, existingResource, newResource)

	nr, err := e.SaveResource(ctx, serviceCtx.ResourceID.String(), newResource, etag)
	if err != nil {
		return nil, err
	}

	versioned, err := converter.HTTPRouteDataModelToVersioned(newResource, serviceCtx.APIVersion)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"ETag": nr.ETag}

	return rest.NewOKResponseWithHeaders(versioned, headers), nil
}

// Validate extracts versioned resource from request and validate the properties.
func (e *CreateOrUpdateHTTPRoute) Validate(ctx context.Context, req *http.Request, apiVersion string) (*datamodel.HTTPRoute, error) {
	serviceCtx := servicecontext.ARMRequestContextFromContext(ctx)
	content, err := ctrl.ReadJSONBody(req)
	if err != nil {
		return nil, err
	}

	dm, err := converter.HTTPRouteDataModelFromVersioned(content, apiVersion)
	if err != nil {
		return nil, err
	}

	dm.ID = serviceCtx.ResourceID.String()
	dm.TrackedResource = ctrl.BuildTrackedResource(ctx)
	// TODO: Update the state.
	dm.Properties.ProvisioningState = v1.ProvisioningStateSucceeded

	return dm, err
}

// UpdateExistingResourceData updates the HTTPRoute resource before it is saved to the DB.
func UpdateExistingResourceData(ctx context.Context, er *datamodel.HTTPRoute, nr *datamodel.HTTPRoute) {
	sc := servicecontext.ARMRequestContextFromContext(ctx)
	nr.SystemData = ctrl.UpdateSystemData(er.SystemData, *sc.SystemData())
	if er.CreatedAPIVersion != "" {
		nr.CreatedAPIVersion = er.CreatedAPIVersion
	}
	nr.TenantID = sc.HomeTenantID
}