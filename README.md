# tryout-watch-completed-pod

- go build -o main *.go
- ./main --kubeconfig=$HOME/.kube/config

## A regular Pod

```yaml
# pod.yaml
kind: Pod
apiVersion: v1
metadata:
  name: pause
spec:
  containers:
  - name: pause
    image: k8s.gcr.io/pause:3.6
```

- k apply -f pod.yaml

```
I1007 11:19:08.655488   32451 main.go:61] podAdd(): name->pause nodeName-> phase->Pending
I1007 11:19:08.655694   32451 main.go:71] podUpdate()
I1007 11:19:08.655706   32451 main.go:72] old: name->pause nodeName-> phase->Pending
I1007 11:19:08.655714   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Pending
I1007 11:19:08.677944   32451 main.go:71] podUpdate()
I1007 11:19:08.677967   32451 main.go:72] old: name->pause nodeName->kind-control-plane phase->Pending
I1007 11:19:08.677979   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Pending
I1007 11:19:09.320245   32451 main.go:71] podUpdate()
I1007 11:19:09.320262   32451 main.go:72] old: name->pause nodeName->kind-control-plane phase->Pending
I1007 11:19:09.320269   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Running```
```

- k delete -f pod.yaml

```
I1007 11:19:20.071331   32451 main.go:71] podUpdate()
I1007 11:19:20.071349   32451 main.go:72] old: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:20.071358   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:20.359158   32451 main.go:71] podUpdate()
I1007 11:19:20.359181   32451 main.go:72] old: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:20.359196   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:26.507247   32451 main.go:71] podUpdate()
I1007 11:19:26.507285   32451 main.go:72] old: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:26.507299   32451 main.go:73] new: name->pause nodeName->kind-control-plane phase->Running
I1007 11:19:26.515936   32451 main.go:93] podDelete(): name->pause nodeName->kind-control-plane phase->Running
```

## A Pod spawned by Job

```yaml
# job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(100)"]
      restartPolicy: Never
  backoffLimit: 4
```

- k apply -f job.yaml

```
I1007 11:21:27.794954   32451 main.go:61] podAdd(): name->pi-bndc8 nodeName-> phase->Pending
I1007 11:21:27.803693   32451 main.go:71] podUpdate()
I1007 11:21:27.803711   32451 main.go:72] old: name->pi-bndc8 nodeName-> phase->Pending
I1007 11:21:27.803719   32451 main.go:73] new: name->pi-bndc8 nodeName->kind-control-plane phase->Pending
I1007 11:21:27.825824   32451 main.go:71] podUpdate()
I1007 11:21:27.825842   32451 main.go:72] old: name->pi-bndc8 nodeName->kind-control-plane phase->Pending
I1007 11:21:27.825851   32451 main.go:73] new: name->pi-bndc8 nodeName->kind-control-plane phase->Pending
I1007 11:22:10.944678   32451 main.go:93] podDelete(): name->pi-bndc8 nodeName->kind-control-plane phase->Pending
```

The last line shows that Pi pod enters Completed state. This proves the client would receive a `Delete` event.

- k delete -f job.yaml

```
# Nothing outputs
```