# Observability in Kubernetes

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

