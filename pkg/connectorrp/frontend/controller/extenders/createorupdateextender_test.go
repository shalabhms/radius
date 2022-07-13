// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package extenders

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/project-radius/radius/pkg/armrpc/asyncoperation/statusmanager"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/project-radius/radius/pkg/connectorrp/api/v20220315privatepreview"
	"github.com/project-radius/radius/pkg/connectorrp/frontend/deployment"
	"github.com/project-radius/radius/pkg/connectorrp/renderers"
	radiustesting "github.com/project-radius/radius/pkg/corerp/testing"
	"github.com/project-radius/radius/pkg/radrp/outputresource"
	"github.com/project-radius/radius/pkg/rp"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/stretchr/testify/require"
)

func getDeploymentProcessorOutputs() (renderers.RendererOutput, deployment.DeploymentOutput) {
	rendererOutput := renderers.RendererOutput{
		SecretValues: map[string]rp.SecretValueReference{
			"secretname": {
				Value: "secretvalue",
			},
		},
		ComputedValues: map[string]renderers.ComputedValueReference{
			"foo": {
				Value: "bar",
			},
		},
	}

	deploymentOutput := deployment.DeploymentOutput{
		Resources: []outputresource.OutputResource{},
	}

	return rendererOutput, deploymentOutput
}

func TestCreateOrUpdateExtender_20220315PrivatePreview(t *testing.T) {
	setupTest := func(tb testing.TB) (func(tb testing.TB), *store.MockStorageClient, *statusmanager.MockStatusManager, *deployment.MockDeploymentProcessor, renderers.RendererOutput, deployment.DeploymentOutput) {
		mctrl := gomock.NewController(t)
		mDeploymentProcessor := deployment.NewMockDeploymentProcessor(mctrl)
		rendererOutput, deploymentOutput := getDeploymentProcessorOutputs()
		mds := store.NewMockStorageClient(mctrl)
		msm := statusmanager.NewMockStatusManager(mctrl)

		return func(tb testing.TB) {
			mctrl.Finish()
		}, mds, msm, mDeploymentProcessor, rendererOutput, deploymentOutput
	}
	createNewResourceTestCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
	}{
		{"create-new-resource-no-if-match", "If-Match", "", "", http.StatusOK, false},
		{"create-new-resource-*-if-match", "If-Match", "*", "", http.StatusPreconditionFailed, true},
		{"create-new-resource-etag-if-match", "If-Match", "random-etag", "", http.StatusPreconditionFailed, true},
		{"create-new-resource-*-if-none-match", "If-None-Match", "*", "", http.StatusOK, false},
	}

	for _, testcase := range createNewResourceTestCases {
		t.Run(testcase.desc, func(t *testing.T) {
			teardownTest, mds, msm, mDeploymentProcessor, rendererOutput, deploymentOutput := setupTest(t)
			defer teardownTest(t)
			input, dataModel, expectedOutput := getTestModels20220315privatepreview()
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(context.Background(), http.MethodGet, testHeaderfile, input)
			req.Header.Set(testcase.headerKey, testcase.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mds.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return nil, &store.ErrNotFound{}
				})
			mDeploymentProcessor.EXPECT().Render(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(rendererOutput, nil)
			mDeploymentProcessor.EXPECT().Deploy(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(deploymentOutput, nil)
			expectedOutput.SystemData.CreatedAt = expectedOutput.SystemData.LastModifiedAt
			expectedOutput.SystemData.CreatedBy = expectedOutput.SystemData.LastModifiedBy
			expectedOutput.SystemData.CreatedByType = expectedOutput.SystemData.LastModifiedByType

			if !testcase.shouldFail {
				mds.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) error {
						obj.ETag = "new-resource-etag"
						obj.Data = dataModel
						return nil
					})
			}

			opts := ctrl.Options{
				StorageClient:  mds,
				AsyncOperation: msm,
				GetDeploymentProcessor: func() deployment.DeploymentProcessor {
					return mDeploymentProcessor
				},
			}

			ctl, err := NewCreateOrUpdateExtender(opts)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, req)
			require.NoError(t, err)
			_ = resp.Apply(ctx, w, req)
			require.Equal(t, testcase.expectedStatusCode, w.Result().StatusCode)

			if !testcase.shouldFail {
				actualOutput := &v20220315privatepreview.ExtenderResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "new-resource-etag", w.Header().Get("ETag"))
			}
		})
	}

	updateExistingResourceTestCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
	}{
		{"update-resource-no-if-match", "If-Match", "", "resource-etag", http.StatusOK, false},
		{"update-resource-*-if-match", "If-Match", "*", "resource-etag", http.StatusOK, false},
		{"update-resource-matching-if-match", "If-Match", "matching-etag", "matching-etag", http.StatusOK, false},
		{"update-resource-not-matching-if-match", "If-Match", "not-matching-etag", "another-etag", http.StatusPreconditionFailed, true},
		{"update-resource-*-if-none-match", "If-None-Match", "*", "another-etag", http.StatusPreconditionFailed, true},
	}

	for _, testcase := range updateExistingResourceTestCases {
		t.Run(testcase.desc, func(t *testing.T) {
			teardownTest, mds, msm, mDeploymentProcessor, rendererOutput, deploymentOutput := setupTest(t)
			defer teardownTest(t)
			input, dataModel, expectedOutput := getTestModels20220315privatepreview()
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(context.Background(), http.MethodGet, testHeaderfile, input)
			req.Header.Set(testcase.headerKey, testcase.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mds.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return &store.Object{
						Metadata: store.Metadata{ID: id, ETag: testcase.resourceETag},
						Data:     dataModel,
					}, nil
				})

			mDeploymentProcessor.EXPECT().Render(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(rendererOutput, nil)
			mDeploymentProcessor.EXPECT().Deploy(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(deploymentOutput, nil)

			if !testcase.shouldFail {
				mds.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) error {
						obj.ETag = "updated-resource-etag"
						obj.Data = dataModel
						return nil
					})
			}

			opts := ctrl.Options{
				StorageClient:  mds,
				AsyncOperation: msm,
				GetDeploymentProcessor: func() deployment.DeploymentProcessor {
					return mDeploymentProcessor
				},
			}

			ctl, err := NewCreateOrUpdateExtender(opts)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, req)
			_ = resp.Apply(ctx, w, req)
			require.NoError(t, err)
			require.Equal(t, testcase.expectedStatusCode, w.Result().StatusCode)

			if !testcase.shouldFail {
				actualOutput := &v20220315privatepreview.ExtenderResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "updated-resource-etag", w.Header().Get("ETag"))
			}
		})
	}
}