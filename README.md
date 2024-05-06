# Observability in Kubernetes

Diagrama de arquitetura:
![1-diagram.png](docs/images/1-diagram.png?raw=true "Architectural diagram")

## Prerrequisitos

- Docker
- kubectl
- k3d

## Configurar cluster

```
mkdir -p /tmp/k3dvol
k3d cluster create k3d-cluster --volume /tmp/k3dvol:/tmp/k3dvol --servers 1
```

## Containerização da aplicação Go

Build da imagem:

```
docker build -t go_web_app .
```

Importar imagem do Docker para cluster k3d:

```
k3d image import go_web_app -c k3d-cluster
```

Deploy da aplicação:
```
kubectl apply -f deploy/go-webapp-deployment.yaml
```

## Métricas com Prometheus e Grafana

Criação de ConfigMap para configurações do Prometheus:

```
kubectl create configmap prometheus-config \
    --from-file=prometheus=./config/prometheus.yml \
    --from-file=prometheus-alerts=./config/alerts.yml
```

Deploy do Prometheus:
```
kubectl apply -f deploy/prometheus-deployment.yaml
```

Deploy do Grafana:
```
kubectl apply -f deploy/grafana-deployment.yaml
```

Saúde dos serviços:
![2-prometheus-ui-healthy.png](docs/images/2-prometheus-ui-healthy.png?raw=true "Prometheus UI Healthy")

Utilizei o Apache Bench para realizar uma simulação de carga com 300 requisições para o serviço de Load Balancer Nginx, que conseguiu distribuir as requisições igualmente utilizando um algoritmo Round Robin dentre os 3 serviços com a aplicação Go.

Métricas no Prometheus:
![3-prometheus-ui-metrics.png](docs/images/3-prometheus-ui-metrics.png?raw=true "Prometheus UI Metrics")

Dashboard Grafana:
![4-grafana-ui.png](docs/images/4-grafana-ui.png?raw=true "Grafana UI")

## Logs com Graylog

Deploy das dependências (MongoDB e OpenSearch):
```
kubectl apply -f deploy/graylog-dependencies-deployment.yaml
```

Deploy do Graylog:
```
kubectl apply -f deploy/graylog-deployment.yaml
```

Configurar GELF via UDP:

- Usuário e senha padrão: _admin_
- **Graylog UI:** (System -> Input -> Select GELF UDP)

Logs:
![5-graylog-ui.png](docs/images/5-graylog-ui.png?raw=true "Graylog UI")

## Traces com Jaeger

Deploy do OpenTelemetry Collector (recebe, processa e envia os traces para o Jaeger):
```
kubectl apply -f deploy/opentelemetry-deployment.yaml
```

Deploy do Jaeger:
```
kubectl apply -f deploy/jaeger-deployment.yaml
```

**Jaeger UI:** http://localhost:16686

Traces no Console do OpenTelemetry Collector:
![6-opentelemetry-collector.png](docs/images/6-opentelemetry-collector.png?raw=true "OpenTelemetry Collector.png")

Jaeger UI:
![7-jaeger-ui.png](docs/images/7-jaeger-ui.png?raw=true "Jaeger UI")
