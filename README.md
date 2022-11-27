# Wiremelt

Extendible automation utility for parallel concurrent worker-pool operations at scale.

## Prerequisites
* NodeJS
* GCC

## Installation
```
go install github.com/khaelou/wiremelt
```

## Usage

```
wiremelt
```

## Flags

> **Client Configurator** requires setup details for the client prior to session initialization; connects required APIs to client.

```
wiremelt client
```
> **Session Configurator** requires setup details for each session prior to client initialization, enforces workload executions.
```
wiremelt session
```
> **Macro** shows Macro Library or used to import custom macros (external JavaScript).
```
wiremelt macro
```
> **Del** deletes the target macro from the client's Macro Library.
```
wiremelt del
```
> **Shell** enables UNIX access, using a “*force operator, embeds built-in comamnds for further usability.
```
wiremelt shell
```
> **Web** launches Web UI client in default web browser, powered by WebAssembly.
```
wiremelt web
```

## Documentation

This project is not yet documented.

## License

This project is maintained under the *GNU GENERAL PUBLIC [LICENSE](/LICENSE)*.
