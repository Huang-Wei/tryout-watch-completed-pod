package main

import (
	"context"
	"flag"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
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

	ctx := context.Background()
	// Create a SharedInformerFactory and watch on Pod change.
	informerFactory := NewInformerFactory(cs, 0)
	informerFactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podAdded,
			UpdateFunc: podUpdated,
			DeleteFunc: podDeleted,
		},
	)

	// Start all informers.
	informerFactory.Start(ctx.Done())

	// Wait for all caches to sync before scheduling.
	informerFactory.WaitForCacheSync(ctx.Done())

	klog.Infof("Start")
	<-ctx.Done()
}

func podAdded(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("podAdd(): name->%v nodeName->%v phase->%v", pod.Name, pod.Spec.NodeName, pod.Status.Phase)
}

func podUpdated(o, n interface{}) {
	old, new := o.(*v1.Pod), n.(*v1.Pod)
	if old.ResourceVersion == new.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	klog.Infof("podUpdate()")
	klog.Infof("old: name->%v nodeName->%v phase->%v", old.Name, old.Spec.NodeName, old.Status.Phase)
	klog.Infof("new: name->%v nodeName->%v phase->%v", new.Name, new.Spec.NodeName, new.Status.Phase)
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
			klog.Errorf("cannot convert to *v1.Pod: %v", t.Obj)
			return
		}
	default:
		klog.Errorf("cannot convert to *v1.Pod: %v", t)
		return
	}

	klog.Infof("podDelete(): name->%v nodeName->%v phase->%v", pod.Name, pod.Spec.NodeName, pod.Status.Phase)
}
