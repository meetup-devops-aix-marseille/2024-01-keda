# 2024-01-keda

Keda is a Kubernetes-based Event Driven Autoscaler. With Keda, you can drive the scaling of any container in Kubernetes based on the number of events needing to be processed.

## Sources

- [KEDA](https://keda.sh/)
- [KEDA GitHub](https://github.com/kedacore/keda)
- [KEDA Helm Chart](https://artifacthub.io/packages/helm/kedacore/keda)

### Install KEDA

```bash
helm upgrade -i --create-namespace --namespace keda --repo https://kedacore.github.io/charts keda keda
helm upgrade -i --create-namespace --namespace keydb --repo https://enapter.github.io/charts keydb keydb --set nodes=1
```

## WERF BUILD

```bash
export WERF_BUILDAH_MODE=auto
werf build --repo  europe-west1-docker.pkg.dev/meetup-devops/images/meetup-keda --dev
werf helm dependency update .helm
werf converge --repo  europe-west1-docker.pkg.dev/meetup-devops/images/meetup-keda --dev
```

k port-forward -n meetup-keda svc/api 8080:8080 &
k port-forward -n meetup-keda svc/api 8081:8081 &

## Add a star

```bash
curl -X POST "localhost:8080/add-task"
```

## Add 100 stars

```bash
for i in {1..100}; do curl -X POST "localhost:8080/add-task"; done
```
