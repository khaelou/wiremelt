# Wiremelt

Extendable utility for parallel concurrent worker-pool operations at scale.

## Prerequisites
* Node.js
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

> **Client Configurator** instructs setup details for the client prior to session configuration; connects required APIs to client.

```
wiremelt client
```
> **Session Configurator** instructs setup details for each session prior to client runtime; sessions handle all workload executions.
```
wiremelt session
```
> **Macro** displays the Macro Library arranged in the session configuation.
```
wiremelt macro
```
> The *macro* flag can also be used to import custom macros (external JavaScript).
```
wiremelt macro SayHi https://example.khaelou.com/sayHi.js
```

> **Del** deletes the target macro from the client's Macro Library.
```
wiremelt del
```
> **Shell** enables UNIX access, a “*force operator embeds system comamnds for further usability.
```
wiremelt shell
```
> **Web** launches Web UI client in default web browser, powered by WebAssembly.
```
wiremelt web
```

> **DND** (Do Not Disturb) dismisses Neural Network executions for sessions configured with NeuralEnabled set to true.
```
wiremelt dnd
```

> **NNET** (Neural Network) activates Neural Network executions for sessions configured with NeuralEnabled set to false.
```
wiremelt nnet
```

## Documentation

This project is not yet documented.

## License

This project is maintained under the *GNU GENERAL PUBLIC [LICENSE](/LICENSE)*.
