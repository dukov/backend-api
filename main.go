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

package main

import (
	"context"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/dukov/backend-api/pkg/service"
	// +kubebuilder:scaffold:imports
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	mgr, err := service.NewApplicationManager(context.Background(), ctrl.GetConfigOrDie())
	if err != nil {
		setupLog.Error(err, "Failed to create Application manager")
		os.Exit(1)
	}

	restful.DefaultContainer.Add(mgr.WebService())
	setupLog.Info("Starting Backend API")
	http.ListenAndServe(":8080", nil)
}
