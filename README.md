# Wiremelt

Extendable AI / ML worker-pool orchestrator.

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

> **Client Configurator** instructs setup details for the client prior to session initialization.

```
wiremelt client
```

> **Session Configurator** instructs setup details for each session prior to workload executions.
```
wiremelt session
```
Pass a value to the *session* flag to activate the target session if it exists.
```
wiremelt session Example
```

> **Macro** displays the Macro Library arranged in the active session configuation.
```
wiremelt macro
```
The *macro* flag can also be used to import custom macros (external JavaScript).
```
wiremelt macro SayHi https://example.khaelou.com/sayHi.js
```

> **Del** removes the target macro from the Macro Library of the active session.
```
wiremelt del
```
> **Shell** provides a built-in command-line interface with commands for additional extendability.
```
wiremelt shell
```
A **force* operator enables UNIX access, which embeds system commands for further usability.
```
>_ ] wiremelt@wm-iMac.local

ls -a *force 
.               .etc
..              .gitignore            LICENSE
```
> **Web** launches API & Web UI in default web browser, powered by WebAssembly.
```
wiremelt web
```

> **Pilot** initiates a navigator for browser automated executions.
```
wiremelt pilot
```

> **DND** (Do Not Disturb) dismisses Neural Network executions for sessions configured with NeuralEnabled set to true.
```
wiremelt dnd
```

> **NNET** (Neural Network) activates Neural Network executions for sessions configured with NeuralEnabled set to false.
```
wiremelt nnet
```

> **Flush** resets all configurations and neural network metadata from Wiremelt, use such flag with caution. Custom macro imports will remain installed.
```
wiremelt flush
```

## Documentation

This project is not yet documented.

## License

This project is maintained under the *GNU GENERAL PUBLIC [LICENSE](/LICENSE)*.
