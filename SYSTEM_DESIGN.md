# SYSTEM_DESIGN.md

## Idempotencia & Reprocesamiento
El ETL asegura idempotencia procesando datos por fecha y claves únicas (fecha, canal, campaña). Reprocesar un rango no genera duplicados en MongoDB.

## Particionamiento & Retención
Los datos se particionan por fecha y canal/campaña. La retención se gestiona a nivel de base de datos (MongoDB) con TTL o limpieza manual.

## Concurrencia & Throughput
Se usan goroutines y worker pools para paralelizar la extracción y carga de datos, maximizando throughput y aprovechando la concurrencia de Go.

## Calidad de datos (UTMs ausentes y fallbacks)
Si faltan UTMs, se aplican valores por defecto o se descartan registros según reglas de negocio. Se loguean los casos para análisis posterior.

## Observabilidad (logs y métricas útiles)
[TODO]El sistema registra logs estructurados (procesos, errores, métricas de ETL). Se pueden integrar métricas Prometheus y trazas para monitoreo.

## Evolución en el ecosistema Admira (data lake/ETL + contratos de API)
[TODO]GoETL puede integrarse con un data lake y evolucionar hacia pipelines más complejos. 
Los contratos de API y OpenAPI facilitan la interoperabilidad y evolución futura.
