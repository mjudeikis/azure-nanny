package subscription

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// DefinitionsOperationMetadataClient is the subscription definitions client provides an interface to create, modify
// and retrieve azure subscriptions programmatically.
type DefinitionsOperationMetadataClient struct {
	BaseClient
}

// NewDefinitionsOperationMetadataClient creates an instance of the DefinitionsOperationMetadataClient client.
func NewDefinitionsOperationMetadataClient() DefinitionsOperationMetadataClient {
	return NewDefinitionsOperationMetadataClientWithBaseURI(DefaultBaseURI)
}

// NewDefinitionsOperationMetadataClientWithBaseURI creates an instance of the DefinitionsOperationMetadataClient
// client.
func NewDefinitionsOperationMetadataClientWithBaseURI(baseURI string) DefinitionsOperationMetadataClient {
	return DefinitionsOperationMetadataClient{NewWithBaseURI(baseURI)}
}

// List lists all of the available Microsoft.Subscription API operations.
func (client DefinitionsOperationMetadataClient) List(ctx context.Context) (result OperationListResultPage, err error) {
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.olr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "List", resp, "Failure sending request")
		return
	}

	result.olr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client DefinitionsOperationMetadataClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	const APIVersion = "2017-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.Subscription/operations"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DefinitionsOperationMetadataClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DefinitionsOperationMetadataClient) ListResponder(resp *http.Response) (result OperationListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client DefinitionsOperationMetadataClient) listNextResults(lastResults OperationListResult) (result OperationListResult, err error) {
	req, err := lastResults.operationListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscription.DefinitionsOperationMetadataClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client DefinitionsOperationMetadataClient) ListComplete(ctx context.Context) (result OperationListResultIterator, err error) {
	result.page, err = client.List(ctx)
	return
}
