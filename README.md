# Load Balancer em Golang

Este repositório contém uma implementação de um load balancer em Golang para fins de estudos.

## Funcionamento

O load balancer é responsável por receber as requisições HTTP dos clientes e distribuí-las entre diversos servidores. Para isso, ele mantém uma lista de URLs de servidores e, a cada requisição, escolhe uma delas utlizando neste caso até então a estratégia de **Round Robin** para enviar a requisição.

## Futuras Melhorias

* Implementar um algoritmo de balanceamento de carga mais sofisticado como Least Connections (baseado no número de conexões ativas em cada servidor).
* Implementar um HealthCheck para monitorar o status dos servidores.
* Implementar um mecanismo de fallback para lidar com falhas de servidor.