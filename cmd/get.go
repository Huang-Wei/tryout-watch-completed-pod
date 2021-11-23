package main

import (
	"context"
	"flag"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
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

// main tests the pods obtained from pre-built informerFactory.
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
	ctx := context.Background()
	// Start all informers.
	informerFactory.Start(ctx.Done())
	// Wait for all caches to sync before scheduling.
	informerFactory.WaitForCacheSync(ctx.Done())

	pods, err := informerFactory.Core().V1().Pods().Lister().List(labels.Everything())
	if err != nil {
		klog.Fatalf("Error obtaining pod from shared informer: %s.", err.Error())
	}

	klog.InfoS("Length of obtained pods", "length", len(pods))

	for _, pod := range pods {
		klog.InfoS("Obtained pod from shared informer", "pod", klog.KObj(pod))
	}
}
