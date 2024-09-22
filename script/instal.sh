cd api && docker build -t thaleseduardo/teste-go . && \
 docker push thaleseduardo/teste-go

# para testar
docker run --rm -p 3000:3000 thaleseduardo/teste-go

# depois de instalar o k3d no seu sistema operacional
# criar cluster e abrir a porta 8000 para acesso externo
k3d cluster create -p "8000:30000@loadbalancer" --agents 2 && kubectl get nodes

# depois de instalar o istio e istioctl no seu sistema operacional
# Instalando istio no cluster
istioctl install
kubectl get namespace
kubectl get services -n istio-system 
kubectl get pods -n istio-system 

# criar namespace
kubectl create namespace api-test-namespace-dev
kubectl get namespace

# aplicar a label do istio no namespaces
kubectl label namespace api-test-namespace-dev istio-injection=enabled

cd kubernetes && kubectl apply -f deployments.yaml && \
 kubectl apply -f server-service.yaml && kubectl apply -f virtual-service.yaml && \ 
 kubectl apply -f destination-rule.yaml && kubectl apply -f gateway.yaml && \
 kubectl apply -f statefulset.yaml && kubectl apply -f postgresql-service-h.yaml

# ou

kubectl apply -f .

# configurar gateway em localhost
# coloque a porta para 30000 no nodePort http2
kubectl edit service istio-ingressgateway -n istio-system
# deve estar assim 80:30000
kubectl get service -n istio-system
kubectl apply -f gateway.yaml
# localhost:8000

# metrics-server
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# no arquivo components.yaml pesquise por Deployment
# não tem tls instalado, então mude esse valor no arquivo components.yaml
# em -args:
# coloque esse valor - --kubelet-insecure-tls
kubectl apply -f metrics-server.yaml && kubectl apply -f hpa.yaml
kubectl get apiservices # v1beta1.metrics.k8s.io, kube-system/metrics-server, True

# Fortio
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.23/samples/httpbin/sample-client/fortio-deploy.yaml
export FORTIO_POD=$(kubectl get pods -l app=fortio -o 'jsonpath={.items[0].metadata.name}') && \
 echo $FORTIO_POD
kubectl exec "$FORTIO_POD" -c fortio -- fortio load -c 2 -qps 0 -t 200s -loglevel Warning http://api-test-service:8000
kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://api-test-service:8000"
watch -n1 kubectl get hpa

# com o istioctl instalado no seu sistema operacional
# baixar kiali, grafana, jaeger e prometheus
kubectl apply -f https://raw.githubusercontent.com/istio/istio/refs/heads/release-1.23/samples/addons/kiali.yaml && \
 kubectl apply -f https://raw.githubusercontent.com/istio/istio/refs/heads/release-1.23/samples/addons/kiali.yaml && \
 kubectl apply -f https://raw.githubusercontent.com/istio/istio/refs/heads/release-1.23/samples/addons/jaeger.yaml && \
 kubectl apply -f https://raw.githubusercontent.com/istio/istio/refs/heads/release-1.23/samples/addons/prometheus.yaml

# verificar os pods de kiali, grafana, jaeger e prometheus
kubectl get pods -n istio-system

# iniciar kiali
istioctl dashboard kiali 

# verificações
kubectl get pvc -n api-test-namespace-dev

curl http://api-test-service:8000
curl http://localhost:8000
kubectl logs postgresql-0 -n api-test-namespace-dev
kubectl get endpoints api-test-service
kubectl exec -it api-test-d8c4b859b-564dg -n api-test-namespace-dev -- bash
kubectl exec -it postgresql-0 -n api-test-namespace-dev -- psql -U postgres
kubectl edit deployment api-test -n api-test-namespace-dev
kubectl logs api-test-d8c4b859b-5vnmt  -c api-test -n api-test-namespace-dev

kubectl get pods -n api-test-namespace-dev
kubectl get services -n api-test-namespace-dev 

kubectl exec -it api-test-5b8fb8c49f-fbgxh -n api-test-namespace-dev -- bash
kubectl get endpoints postgresql-h -n api-test-namespace-dev   