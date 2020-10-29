/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"net/http"

	"github.com/emicklei/go-restful"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	goclientrest "k8s.io/client-go/rest"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dukovv1alpha1 "github.com/dukov/backend-api/api/v1alpha1"
)

// ApplicationManager for controlling application api
type ApplicationManager struct {
	Client  client.Client
	Context context.Context
}

// NewApplicationManager creates Applicationmanager
func NewApplicationManager(ctx context.Context, cfg *goclientrest.Config) (*ApplicationManager, error) {
	client, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return &ApplicationManager{
		Client:  client,
		Context: ctx,
	}, nil

}

func (am *ApplicationManager) WebService() *restful.WebService {
	ws := &restful.WebService{}
	ws.Path("/apps").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Param(
			ws.PathParameter("app-name", "identifier of the user").
				DataType("string"),
		)
	for _, bldr := range am.RouteBuilders(ws) {
		ws.Route(bldr)
	}
	return ws
}

// RouteBuilders list of pathes to handle
func (am *ApplicationManager) RouteBuilders(ws *restful.WebService) []*restful.RouteBuilder {
	return []*restful.RouteBuilder{
		ws.GET("/").
			To(am.List).
			Returns(http.StatusOK, "OK", []dukovv1alpha1.Application{}),
		ws.GET("/{app-name}").
			To(am.Get).
			Returns(http.StatusOK, "OK", []*unstructured.Unstructured{}),
	}
}

// List resources
func (am *ApplicationManager) List(request *restful.Request, response *restful.Response) {
	appList := &dukovv1alpha1.ApplicationList{}
	if err := am.Client.List(am.Context, appList); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
	response.WriteAsJson(appList.Items)
}

// Get resource
func (am *ApplicationManager) Get(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("app-name")
	ns := request.QueryParameter("namespace")
	if ns == "" {
		ns = "default"
	}
	appObj := &dukovv1alpha1.Application{}
	objectKey := client.ObjectKey{Name: name, Namespace: ns}
	if err := am.Client.Get(am.Context, objectKey, appObj); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
	result := make([]*unstructured.Unstructured, len(appObj.Spec.Resources))
	for i, objRef := range appObj.Spec.Resources {
		resource, err := external.Get(context.TODO(), am.Client, objRef, objRef.Namespace)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		result[i] = resource
	}
	response.WriteAsJson(result)
}
