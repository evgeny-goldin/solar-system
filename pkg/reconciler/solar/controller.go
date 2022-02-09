package solar

import (
	"context"
	"k8s.io/client-go/tools/cache"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"solar-system/pkg/apis/solar/v1alpha1"
	starclient "solar-system/pkg/client/injection/client"
	starinformer "solar-system/pkg/client/injection/informers/solar/v1alpha1/star"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	starInformer := starinformer.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)
	r := &Reconciler{
		deploymentLister: deploymentInformer.Lister(),
		starLister:       starInformer.Lister(),
		starClient:       starclient.Get(ctx),
		KubeClientSet:    kubeclient.Get(ctx),
	}
	impl := controller.NewImpl(r, logger, "Star")
	logger.Info("Setting up event handlers.")
	starInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(v1alpha1.SchemeGroupVersion.WithKind("Star")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})
	return impl
}
