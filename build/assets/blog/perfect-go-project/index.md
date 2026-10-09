---
title: Go project structure IMO!!!
date: "2026-08-10"
summary: | This is a very opinionated golang(ish) project structure; that from  my experience works best on scaled team

---

# Go project structure IMO!!!

| This is a very opinionated golang(ish) project structure; that from  my experience works best on scaled team
| it's Mix of patterns and principles one learns over the year also convinced this can be similarly built in other languages.

## Baselines

### Folder structure should be as following:
```
root/
|_ pkg/ <- contains shared libs functionality that are designed to hide out external pkg calls behind a SOLID desinged package that can be reused internally without tiding internals to external packages/modules

|_ test/ <- common test helper functions, each function here must receive *testing.T and call t.Helper() as first action
    |_ fixtures/ <-- loading fixtures options
        |_ fixture.go <-- helper funcitons to load fixtures
        |_ db/ <- datastore fixtures to be loaded
        |_ json/ <-- json based fixtures that can be used to compare API response
    |_ grpc/ <-- grpc mock servers if we need to mock external dep servers calls
    |_ http/ <-- http mock servers if we need to mock external dep servers calls
|
|
|_  cmd/
    |-- root.go <- should be the starting entry point for all commands availble, through a helper function ExcuteWithConfig(cmd, config) this enables run any of the commands passing config
    |-- <Command>.go <-- each command should live in it's own file
    |-- <Command_test>.go <-- each command file should have an equal command test file that use ExecuteWithConfig through test/
|
|_ internal/
  |--> <app_name>/ <- apps should never leak into each others and only allowed to call each others through servers, pubsub messaging or gRPC.
        |--> app/ <-- pkg that implements pkg/app interfaces so it can be registerd as app directly. should never recieve be created from container, instead it uses container to glue everything together creating the http, grpc, publishers, subscribers, migrator and cronJobs from one place, at the end cmd should register this app, and able to run the app command/server directly
        |
        |--> config/ <- pkg that is responsible on only reading environment variables, settings files, or any configuration, but strictly without init any object or model, at the end it should deliver one config struct with all values either read or default or empty
        |
        |--> container/ <- responsible only on create or get from cached deps map any service, controller, adapter etc using the configs passed, but not responsible on starting routines, schedule events. if the service it starts needs shutdown, it should register it to shutdown lib
        if the package being created needs a client, tracer, logger, metrics etc it has to be provided through the container to unify the creation of deps
            |
            |--> externals.go <- external clients init
            |--> infrastructure.go <-- infra SRE related deps, like tracer, logger, metric
            |--> services.go <-- init of all services
            |--> adapters.go <-- init all adapters
            |--> controllers.go <-- init controllers
            |--> endpoints.go <-- init endpoints with decoders, encoders and controllers
            |--> db.go <-- init data-stores connections only (read, writes)
            |--> repositories.go <-- init repos only
        |
        |--> repositories/ <-- Responsible to translate from/to entity to/from DataStore model without knowing anything about the connection itself and do operation.
            |
            |--> http/ <-- http clients like other services
                    |
                    |--> DTOs/ <-- responsible to translate from entity to data model only if data model is needed
                        |-> exampleDTO/ <-- has all translations from example entity to any data type model, model structs are also saved here
                 |--> ExampleRepo <-- Name should never indicate how this repo call it's data store
            |
            |--> grpc/ <-- grpc clients
            |
            |--> sql/
                |
                |--> DTOs/ <-- responsible to translate from entity to data model only if data model is needed
                    |-> exampleDTO/ <-- has all translations from example entity to any data type model, model structs are also saved here
                |--> example.go <- exampleRepo Interface implementation with sqldb, SQL queries should be rendered/formed here only.
            |--> file/
                |--> example.go <- exampleRepo Interface implementation with file store system 
            |--> cache/
                |--> inmemory/
                    |--> example.go <- exampleRepo Interface implementation with local inmemory cache build, with lazyloading at the start calling sql.ExampleRepo to load whatever then cache it. to serve from cache later, writes through this layer should be in order sql write then cache write
                |--> redis/
                    |--> example.go <- exampleRepo Interface implementation with remote cache build, with lazyloading at the start calling sql.ExampleRepo to load whatever then cache it. to serve from cache later, writes through this layer should be in order sql write then cache write
        |
        |--> controllers/ <-- responsible to call usecase in order, from incoming request struct, return error or response struct, this request/response structs are entities and must not have any tags
        |
        |--> entities/ <-- each entity is a type business definition contract that can be shared across controllers->services->repository, entity might have the entities as type as well. and only embed self contained simple functions are allowed to belong to entity example Entities.find(id) will find the entity with id through the Entities list
        |
        |--> services/ <-- a small functions grouped based on business case for example ExampleService responsible to create tasks, workID related to task through right repo; Example2Service responsible to follow on session lifecycle and report back, if session got cancelled or ask questions Example2Service should take the correct corresponding actions
        |
        |--> usecases/ <-- it's thin line between service responsibiliy and usecase; usecases is a flow of complete business case.
        |
        |--> transport/
            | middlewares/
                |--> <middlewareExample>
                    |- middlewareExample.go <-- the main functionality that should happen receive the parameters it needs and verify example JWT key not the headers, 
                    |- http.go <-- the middleware functionality in http mux middleware calls the main funciton in middlewareExample.go
                    |- grpc.go <-- the middleware functionality in http mux middleware calls the main funciton in middlewareExample.go
            |
            |--> http/
            |--> views/ <-- for http handler that render templates, or serve HTML, CSS and JS files directly
                    |--> exampleView.go
                    |--> exampleView2.go
            |--> api/
                    |--> <endpoints>/ 
                        |--> endpoint.go <-- should be the endpoint handler to call encoder -> controller -> decoder
                        |--> encoder.go <-- should translate the incoming http request data type (ex json, form, params) to Controller Request Struct
                        |--> decoder.go <-- should translate the outgoing http request data type (ex json, form, params) from Controller Response Struct
                    |--> decoder/ <-- general decoder for example decode Error Response
            |
            |--> grpc/
                |--> <endpoints>/ 
                    |--> endpoint.go <-- should be the endpoint handler to call encoder -> controller -> decoder
                    |--> encoder.go <-- should translate the incoming grpc request proto type to Controller Request Struct
                    |--> decoder.go <-- should translate the outgoing grpc request proto type from Controller Response Struct
            |
            |--> socketIO/
                |--> <RPC>/ 
                    |--> endpoint.go <-- should be the endpoint handler to call encoder -> controller -> decoder
                    |--> encoder.go <-- should translate the incoming io request type to Controller Request Struct
                    |--> decoder.go <-- should translate the outgoing io request type from Controller Response Struct
            |
            |-> servers/
                |_ http.go <- create the http server router/mux, register endpoints(handlers), middlewares, 
                |_ grpc.go <- create the grpc server impl, registering right endpoints, middlewares etc
                |_ socket.go <- create socketIO server implemntaion, registering RPC endpoints, middlewares etc
                |_ pubsubs/ <- create publisher or subscriber for certain topic
                    <examplePublisher>/
                    <exampleSubscriber>/
            |-> schedulars/ <-- autostart cronjobs based on timetickers the assumption is caller is calling from goroutine (or not depend on need setup), must outcall functions that requires only ctx as input and are self sufficient from configurations read, controlling/managing the schedule, overlaps actions, error status etc.
                |
                |--> You might want to create schedular registery in pkg/schedular where any internal package can easily register a function to be cron scheduled or scheduled to run one time in the future.
                registery should have the ability to gracefuly shutdown all running jobs in the background, when the schedular main context is cancelled or receive shutdown call on schedularRegistery.Close() it should wait until all running jobs finish, but prevent any new one from starting

```

---

### Rules that are not allowed to break

#### Design

- Must follow strictly SOLID principles (Single Responsibility, Open–Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion)
- use but not limited to: Builder, Factory, Separation of concern, pubsub messaging (using internal go channels) patterns

#### Importing

- importing should happen inwards only and never the other way around, so for the following (x --can call--> y but y can not call x)
```
cmd --> internal.app -> transport.servers -> transport.endpoints -> controllers -> usecases --> services --> repositories --> DTOs

OR

cmd --> internal.app -> transport.schedular -> controllers -> usecases --> services --> repositories --> DTOs
```
- any internal package only can import from internal/entities
- only `internal/<app>/container` can import from pkg/
- only `internal/<app>/container` and `cmd/` allowed to import `internal/<app>/app/`
- all Struct must receive interfaces if data struct is not a concern but functions are.

#### Testing

- in tests always use mocks and make sure functions has been called with expected passed input, and returning expected output to cleanly test function
- only `internal/<app>/repositories` allowed to use `internal/<app>/repositories/<type>/<DTO>s`
- every go files that has functionality must be tested in corresponding _test.go file
- only tracer can be provided as noop in testing across all code.
- DB / broker: integrate, don't mock
    For code that exercises SQL or a message broker, prefer **real** dependencies in
    `integration`-tagged tests over `sqlmock`:

    - `test/dbtest` — open a real Postgres from `TEST_DATABASE_DSN`, run fixtures/DDL
    (`Exec`, a schema `const`), `Truncate`, and assert with `Count`.
    - `test/pubsubtest` — a `Collector` handler + `Publish` over `pkg/pubsub` to assert
    a consumer acks/redelivers N messages on any backend (memory/NATS/GCP).

#### pkgs

Libraries in `pkg/` MUST NOT leak their third-party SDK to callers. The app
depends only on interfaces you own, so the backend can be swapped (e.g. OTEL →
Datadog) by writing a new implementation, with zero caller changes.


1. **Own the interface.** Expose your own `Provider` / `Meter` / `Client` /
   instrument interfaces. Method signatures use only stdlib and your own types
   (`context.Context`, `map[string]any`, `string`, `time.Duration`). Never a
   vendor type in an exported signature or struct field that callers touch.
2. **Hide the implementation.** Keep the concrete vendor-backed structs and the
   vendor imports in unexported files (`provider.go`, `meter.go`, `utils.go`).
   Only this package imports the SDK.
3. **Options don't leak.** Functional options (`WithX`) accept plain values and
   convert internally (e.g. `WithAttributes(map[string]any)` builds vendor
   attributes inside the package). Accept loggers via an unexported interface so
   the app's logger satisfies it without this package importing it.
4. **Provide a no-op.** Ship a `NewNoopXxx()` and return it when the feature is
   disabled (`Config.Enabled == false`) or on non-fatal build issues, so
   dependent code always has a non-nil, safe implementation. Mirror
   `noopTracerProvider`.
5. **Config + constructor.** Add a `Config` with `mapstructure` tags and
   `NewConfig()` defaults. The `NewXxxProvider(ctx, cfg, ...opts)` returns your
   interface and a wrapped error; disabled config returns the no-op.
6. **Package naming.** If the package name would clash with the SDK import (otel
   `metric`/`trace`), name the package to avoid it (dir `pkg/metric` → package
   `metrics`, like dir `pkg/tracer` → package `tracing`).
7. **doc.go** explains the vendor-neutral intent and names the swap target.

#### Instrumentation

- logging must be used when appropriate with right amount of log level.
- every package (controllers, services, usecases, repositories) should have internal `tracing/` package that wraps all functions with tracing span using the provided tracer, adding good enough data to the tracer for debugging
for example:
```
usecases/
    <example>/
        tracing/<example>.go // implements same example interface, calls example.function within it's own tracingExample.function after wrapping it with span, and receive the response as well
        <example>.go // should never know anything about tracer at any given point
```

- OTLP-gRPC exporter: `NewMeterProvider(Enabled=true, ...)` succeeds with no
  collector; call `Meter().Counter(...).Add(...)` and `Stop()` (use a short
  `resourceTimeout` so `Stop` returns fast).
- `sql.OpenDB` / pgx: a syntactically valid DSN
  (`postgres://user:pass@localhost:5432/db?sslmode=disable`) opens without
  dialing; exercise the getter and caching.
- Invalid otel instrument name (empty string) → constructor returns an error →
  exercises the log-and-fallback-to-noop branch.
- everything must be derived or built from one original command context

#### Commands

- there should be always a command that can run all commands together directly (cronjobs, publisher and subscriber using goroutines and channels internal comms only, grpc, http) as a one self contained application binary

#### Dev environment

- the app must be runnable without the need to install any other tool locally than docker, use docker compose (or similar podman, minikube) to gain a fully working local dev environment without requiring of setup
- makefiles rules:
    - if setup is really needed (limited to copy files from somewhere to elsewhere before starting app ex: make copy of .env.sample file to .env) we must have make command that brings everything up
    - there must be a make command to clean the project from auto-generated files that is not needed (to keep the local env clean and tidy)
    - running all types of test should have each a command in make file

## Anti-patterns

### pkg/
- Returning `*sdk.Thing` "just for now" — it becomes load-bearing and blocks the swap.
- Re-exporting SDK option types.
- Putting SDK imports in the same file as the public interface.

