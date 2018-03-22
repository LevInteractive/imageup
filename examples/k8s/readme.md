# Imageup w/ Kubernetes

```
kubectl apply -f imageup-deployment.yml -f imageup-service.yml
```

From a pod running in the cluster, you may now connect to imageup via
`http://imageup-service:31111`. This example is assuming you're using Google
Cloud.
