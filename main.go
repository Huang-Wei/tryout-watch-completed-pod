package main

import (
	"context"
	"flag"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/Huang-Wei/tryout-watch-completed-pod/pkg"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s.", err.Error())
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s.", err.Error())
	}

	// Create a SharedInformerFactory and watch on Pod change.
	informerFactory := pkg.NewInformerFactory(cs, 0)
	informerFactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podAdded,
			UpdateFunc: podUpdated,
			DeleteFunc: podDeleted,
		},
	)

	klog.InfoS("Start")

	ctx := context.Background()
	// Start all informers.
	informerFactory.Start(ctx.Done())
	// Wait for all caches to sync before scheduling.
	informerFactory.WaitForCacheSync(ctx.Done())

	<-ctx.Done()
}

func podAdded(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.InfoS("Pod added", "pod", klog.KObj(pod), "assignedNodeName", pod.Spec.NodeName, "phase", pod.Status.Phase)
}

func podUpdated(o, n interface{}) {
	old, new := o.(*v1.Pod), n.(*v1.Pod)
	if old.ResourceVersion == new.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	klog.InfoS("Pod updated")
	klog.InfoS("Old pod", "pod", klog.KObj(old), "assignedNodeName", old.Spec.NodeName, "phase", old.Status.Phase)
	klog.InfoS("New pod", "pod", klog.KObj(new), "assignedNodeName", new.Spec.NodeName, "phase", new.Status.Phase)
}

func podDeleted(obj interface{}) {
	var pod *v1.Pod
	switch t := obj.(type) {
	case *v1.Pod:
		pod = t
	case cache.DeletedFinalStateUnknown:
		var ok bool
		pod, ok = t.Obj.(*v1.Pod)
		if !ok {
			klog.ErrorS(nil, "Cannot convert to *v1.Pod", "pod", t.Obj)
			return
		}
	default:
		klog.ErrorS(nil, "Cannot convert to *v1.Pod", "pod", t)
		return
	}

	klog.InfoS("Pod deleted", "pod", klog.KObj(pod), "assignedNodeName", pod.Spec.NodeName, "phase", pod.Status.Phase)
}
