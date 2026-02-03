# Available Metrics

## go_gc

- `go_gc_duration_seconds` (summary) — A summary of the pause duration of garbage collection cycles.
  Labels: instance (2 values), job (2 values), quantile (5 values)
  Variable candidates: quantile
  Legend candidates: instance, job
- `go_gc_duration_seconds_count`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_gc_duration_seconds_sum`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_gc_gogc_percent` (gauge) — Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `go_gc_gomemlimit_bytes` (gauge) — Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"

## go_goroutines

- `go_goroutines` (gauge) — Number of goroutines that currently exist.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## go_info

- `go_info` (gauge) — Information about the Go environment.
  Labels: instance (2 values), job (2 values), version (2 values)
  Legend candidates: instance, job, version

## go_memstats

- `go_memstats_alloc_bytes` (gauge) — Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_alloc_bytes_total` (counter) — Total number of bytes allocated, even if freed.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_buck_hash_sys_bytes` (gauge) — Number of bytes used by the profiling bucket hash table.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_frees_total` (counter) — Total number of frees.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_gc_sys_bytes` (gauge) — Number of bytes used for garbage collection system metadata. Equals to /memory/classes/metadata/other:bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_alloc_bytes` (gauge) — Number of heap bytes allocated and still in use.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_idle_bytes` (gauge) — Number of heap bytes waiting to be used. Equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_inuse_bytes` (gauge) — Number of heap bytes that are in use. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_objects` (gauge) — Number of allocated objects.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_released_bytes` (gauge) — Number of heap bytes released to OS. Equals to /memory/classes/heap/released:bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_heap_sys_bytes` (gauge) — Number of heap bytes obtained from system. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_last_gc_time_seconds` (gauge) — Number of seconds since 1970 of last garbage collection.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_lookups_total` (counter) — Total number of pointer lookups.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="traefik:8080", job="traefik"
- `go_memstats_mallocs_total` (counter) — Total number of mallocs.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_mcache_inuse_bytes` (gauge) — Number of bytes in use by mcache structures.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_mcache_sys_bytes` (gauge) — Number of bytes used for mcache structures obtained from system.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_mspan_inuse_bytes` (gauge) — Number of bytes in use by mspan structures.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_mspan_sys_bytes` (gauge) — Number of bytes used for mspan structures obtained from system.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_next_gc_bytes` (gauge) — Number of heap bytes when next garbage collection will take place.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_other_sys_bytes` (gauge) — Number of bytes used for other system allocations.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_stack_inuse_bytes` (gauge) — Number of bytes in use by the stack allocator.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_stack_sys_bytes` (gauge) — Number of bytes obtained from system for stack allocator.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `go_memstats_sys_bytes` (gauge) — Number of bytes obtained from system.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## go_sched

- `go_sched_gomaxprocs_threads` (gauge) — The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. Sourced from /sched/gomaxprocs:threads.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"

## go_threads

- `go_threads` (gauge) — Number of OS threads created.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## myapp_cache

- `myapp_cache_hits_total` (counter) — Total cache hits.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `myapp_cache_misses_total` (counter) — Total cache misses.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"

## myapp_db

- `myapp_db_connections_active` (gauge) — Number of active DB connections.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `myapp_db_connections_idle` (gauge) — Number of idle DB connections.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `myapp_db_query_duration_seconds_bucket`
  Labels: instance (1 values), job (1 values), le (12 values), operation (2 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Variable candidates: le
  Legend candidates: operation
- `myapp_db_query_duration_seconds_count`
  Labels: instance (1 values), job (1 values), operation (2 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: operation
- `myapp_db_query_duration_seconds_sum`
  Labels: instance (1 values), job (1 values), operation (2 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: operation

## myapp_errors

- `myapp_errors_total` (counter) — Total errors by type.
  Labels: instance (1 values), job (1 values), type (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: type

## myapp_http

- `myapp_http_request_duration_seconds_bucket`
  Labels: instance (1 values), job (1 values), le (12 values), method (2 values), path (4 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Variable candidates: le
  Legend candidates: method, path
- `myapp_http_request_duration_seconds_count`
  Labels: instance (1 values), job (1 values), method (2 values), path (4 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: method, path
- `myapp_http_request_duration_seconds_sum`
  Labels: instance (1 values), job (1 values), method (2 values), path (4 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: method, path
- `myapp_http_requests_in_flight` (gauge) — Number of HTTP requests currently in flight.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `myapp_http_requests_total` (counter) — Total HTTP requests.
  Labels: instance (1 values), job (1 values), method (2 values), path (4 values), status (2 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: method, status, path

## myapp_jobs

- `myapp_jobs_duration_seconds_bucket`
  Labels: instance (1 values), job (1 values), le (12 values), queue (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Variable candidates: le
  Legend candidates: queue
- `myapp_jobs_duration_seconds_count`
  Labels: instance (1 values), job (1 values), queue (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: queue
- `myapp_jobs_duration_seconds_sum`
  Labels: instance (1 values), job (1 values), queue (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: queue
- `myapp_jobs_processed_total` (counter) — Total background jobs processed.
  Labels: instance (1 values), job (1 values), queue (3 values), status (2 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: status, queue
- `myapp_jobs_queue_depth` (gauge) — Number of pending jobs in queue.
  Labels: instance (1 values), job (1 values), queue (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: queue

## myapp_orders

- `myapp_orders_created_total` (counter) — Total orders created.
  Labels: instance (1 values), job (1 values), status (3 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: status

## myapp_revenue

- `myapp_revenue_total` (counter) — Total revenue.
  Labels: currency (1 values), instance (1 values), job (1 values)
  Fixed: currency="USD", instance="dummyapp:3000", job="dummyapp"

## myapp_users

- `myapp_users_active` (gauge) — Number of currently active users.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `myapp_users_registered_total` (counter) — Total user registrations.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"

## process_cpu

- `process_cpu_seconds_total` (counter) — Total user and system CPU time spent in seconds.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## process_max

- `process_max_fds` (gauge) — Maximum number of open file descriptors.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## process_network

- `process_network_receive_bytes_total` (counter) — Number of bytes received by the process over the network.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `process_network_transmit_bytes_total` (counter) — Number of bytes sent by the process over the network.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"

## process_open

- `process_open_fds` (gauge) — Number of open file descriptors.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## process_resident

- `process_resident_memory_bytes` (gauge) — Resident memory size in bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## process_start

- `process_start_time_seconds` (gauge) — Start time of the process since unix epoch in seconds.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## process_virtual

- `process_virtual_memory_bytes` (gauge) — Virtual memory size in bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `process_virtual_memory_max_bytes` (gauge) — Maximum amount of virtual memory available in bytes.
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## promhttp_metric

- `promhttp_metric_handler_requests_in_flight` (gauge) — Current number of scrapes being served.
  Labels: instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
- `promhttp_metric_handler_requests_total` (counter) — Total number of scrapes by HTTP status code.
  Labels: code (3 values), instance (1 values), job (1 values)
  Fixed: instance="dummyapp:3000", job="dummyapp"
  Legend candidates: code

## redis_clients

- `redis_clients_blocked`
- `redis_clients_connected`
- `redis_clients_max_input_buffer_bytes`
- `redis_clients_max_output_buffer_bytes`

## redis_commands

- `redis_commands_per_second`
- `redis_commands_processed_total`

## redis_connections

- `redis_connections_received_total`
- `redis_connections_rejected_total`

## redis_cpu

- `redis_cpu_time_seconds_total`
  Labels: state (6 values)
  Variable candidates: state

## redis_keys

- `redis_keys_evicted_total`
- `redis_keys_expired_total`

## redis_keyspace

- `redis_keyspace_hits_total`
- `redis_keyspace_misses_total`

## redis_latest

- `redis_latest_fork_microseconds`

## redis_memory

- `redis_memory_fragmentation_ratio`
- `redis_memory_lua_bytes`
- `redis_memory_peak_bytes`
- `redis_memory_rss_bytes`
- `redis_memory_used_bytes`

## redis_net

- `redis_net_input_bytes_total`
- `redis_net_output_bytes_total`

## redis_rdb

- `redis_rdb_changes_since_last_save`

## redis_replication

- `redis_replication_backlog_first_byte_offset_bytes`
- `redis_replication_offset_bytes`

## redis_slaves

- `redis_slaves_connected`

## redis_uptime

- `redis_uptime_seconds_total`

## scrape_duration

- `scrape_duration_seconds`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## scrape_samples

- `scrape_samples_post_metric_relabeling`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job
- `scrape_samples_scraped`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## scrape_series

- `scrape_series_added`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

## system_cpu

- `system_cpu_time_seconds_total`
  Labels: cpu (2 values), state (8 values)
  Variable candidates: state
  Legend candidates: cpu

## system_disk

- `system_disk_io_bytes_total`
  Labels: device (7 values), direction (2 values)
  Variable candidates: device
  Legend candidates: direction
- `system_disk_io_time_seconds_total`
  Labels: device (7 values)
  Variable candidates: device
- `system_disk_merged_total`
  Labels: device (7 values), direction (2 values)
  Variable candidates: device
  Legend candidates: direction
- `system_disk_operation_time_seconds_total`
  Labels: device (7 values), direction (2 values)
  Variable candidates: device
  Legend candidates: direction
- `system_disk_operations_total`
  Labels: device (7 values), direction (2 values)
  Variable candidates: device
  Legend candidates: direction
- `system_disk_pending_operations`
  Labels: device (7 values)
  Variable candidates: device
- `system_disk_weighted_io_time_seconds_total`
  Labels: device (7 values)
  Variable candidates: device

## system_memory

- `system_memory_usage_bytes`
  Labels: state (6 values)
  Variable candidates: state

## system_network

- `system_network_connections`
  Labels: protocol (1 values), state (12 values)
  Fixed: protocol="tcp"
  Variable candidates: state
- `system_network_dropped_total`
  Labels: device (2 values), direction (2 values)
  Legend candidates: device, direction
- `system_network_errors_total`
  Labels: device (2 values), direction (2 values)
  Legend candidates: device, direction
- `system_network_io_bytes_total`
  Labels: device (2 values), direction (2 values)
  Legend candidates: device, direction
- `system_network_packets_total`
  Labels: device (2 values), direction (2 values)
  Legend candidates: device, direction

## traefik_config

- `traefik_config_last_reload_success` (gauge) — Last config reload success
  Labels: instance (1 values), job (1 values)
  Fixed: instance="traefik:8080", job="traefik"
- `traefik_config_reloads_total` (counter) — Config reloads
  Labels: instance (1 values), job (1 values)
  Fixed: instance="traefik:8080", job="traefik"

## traefik_entrypoint

- `traefik_entrypoint_request_duration_seconds_bucket`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"
  Variable candidates: le
- `traefik_entrypoint_request_duration_seconds_count`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"
- `traefik_entrypoint_request_duration_seconds_sum`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"
- `traefik_entrypoint_requests_bytes_total` (counter) — The total size of requests in bytes handled by an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"
- `traefik_entrypoint_requests_total` (counter) — How many HTTP requests processed on an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"
- `traefik_entrypoint_responses_bytes_total` (counter) — The total size of responses in bytes handled by an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
  Fixed: code="200", entrypoint="web", instance="traefik:8080", job="traefik", method="GET", protocol="http"

## traefik_open

- `traefik_open_connections` (gauge) — How many open connections exist, by entryPoint and protocol
  Labels: entrypoint (2 values), instance (1 values), job (1 values), protocol (1 values)
  Fixed: instance="traefik:8080", job="traefik", protocol="TCP"
  Legend candidates: entrypoint

## traefik_router

- `traefik_router_request_duration_seconds_bucket`
  Labels: code (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"
  Variable candidates: le
- `traefik_router_request_duration_seconds_count`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"
- `traefik_router_request_duration_seconds_sum`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"
- `traefik_router_requests_bytes_total` (counter) — The total size of requests in bytes handled by a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"
- `traefik_router_requests_total` (counter) — How many HTTP requests are processed on a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"
- `traefik_router_responses_bytes_total` (counter) — The total size of responses in bytes handled by a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", router="whoami@file", service="whoami@file"

## traefik_service

- `traefik_service_request_duration_seconds_bucket`
  Labels: code (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"
  Variable candidates: le
- `traefik_service_request_duration_seconds_count`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"
- `traefik_service_request_duration_seconds_sum`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"
- `traefik_service_requests_bytes_total` (counter) — The total size of requests in bytes received by a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"
- `traefik_service_requests_total` (counter) — How many HTTP requests processed on a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"
- `traefik_service_responses_bytes_total` (counter) — The total size of responses in bytes returned by a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
  Fixed: code="200", instance="traefik:8080", job="traefik", method="GET", protocol="http", service="whoami@file"

## up

- `up`
  Labels: instance (2 values), job (2 values)
  Legend candidates: instance, job

# Label Values

Full label value listings for each metric.

## go_gc_duration_seconds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik
- **quantile**: 0.0, 0.25, 0.5, 0.75, 1.0

## go_gc_duration_seconds_count

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_gc_duration_seconds_sum

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_gc_gogc_percent

- **instance**: dummyapp:3000
- **job**: dummyapp

## go_gc_gomemlimit_bytes

- **instance**: dummyapp:3000
- **job**: dummyapp

## go_goroutines

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_info

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik
- **version**: go1.24.5, go1.25.6

## go_memstats_alloc_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_alloc_bytes_total

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_buck_hash_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_frees_total

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_gc_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_alloc_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_idle_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_inuse_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_objects

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_released_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_heap_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_last_gc_time_seconds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_lookups_total

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mallocs_total

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_mcache_inuse_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_mcache_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_mspan_inuse_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_mspan_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_next_gc_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_other_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_stack_inuse_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_stack_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_memstats_sys_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## go_sched_gomaxprocs_threads

- **instance**: dummyapp:3000
- **job**: dummyapp

## go_threads

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## myapp_cache_hits_total

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_cache_misses_total

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_db_connections_active

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_db_connections_idle

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_db_query_duration_seconds_bucket

- **instance**: dummyapp:3000
- **job**: dummyapp
- **le**: +Inf, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 10.0, 2.5, 5.0
- **operation**: insert, select

## myapp_db_query_duration_seconds_count

- **instance**: dummyapp:3000
- **job**: dummyapp
- **operation**: insert, select

## myapp_db_query_duration_seconds_sum

- **instance**: dummyapp:3000
- **job**: dummyapp
- **operation**: insert, select

## myapp_errors_total

- **instance**: dummyapp:3000
- **job**: dummyapp
- **type**: internal, timeout, validation

## myapp_http_request_duration_seconds_bucket

- **instance**: dummyapp:3000
- **job**: dummyapp
- **le**: +Inf, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 10.0, 2.5, 5.0
- **method**: GET, POST
- **path**: /, /api/orders, /api/search, /api/users

## myapp_http_request_duration_seconds_count

- **instance**: dummyapp:3000
- **job**: dummyapp
- **method**: GET, POST
- **path**: /, /api/orders, /api/search, /api/users

## myapp_http_request_duration_seconds_sum

- **instance**: dummyapp:3000
- **job**: dummyapp
- **method**: GET, POST
- **path**: /, /api/orders, /api/search, /api/users

## myapp_http_requests_in_flight

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_http_requests_total

- **instance**: dummyapp:3000
- **job**: dummyapp
- **method**: GET, POST
- **path**: /, /api/orders, /api/search, /api/users
- **status**: 200, 500

## myapp_jobs_duration_seconds_bucket

- **instance**: dummyapp:3000
- **job**: dummyapp
- **le**: +Inf, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 10.0, 2.5, 5.0
- **queue**: cleanup, email, export

## myapp_jobs_duration_seconds_count

- **instance**: dummyapp:3000
- **job**: dummyapp
- **queue**: cleanup, email, export

## myapp_jobs_duration_seconds_sum

- **instance**: dummyapp:3000
- **job**: dummyapp
- **queue**: cleanup, email, export

## myapp_jobs_processed_total

- **instance**: dummyapp:3000
- **job**: dummyapp
- **queue**: cleanup, email, export
- **status**: failure, success

## myapp_jobs_queue_depth

- **instance**: dummyapp:3000
- **job**: dummyapp
- **queue**: cleanup, email, export

## myapp_orders_created_total

- **instance**: dummyapp:3000
- **job**: dummyapp
- **status**: cancelled, completed, failed

## myapp_revenue_total

- **currency**: USD
- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_users_active

- **instance**: dummyapp:3000
- **job**: dummyapp

## myapp_users_registered_total

- **instance**: dummyapp:3000
- **job**: dummyapp

## process_cpu_seconds_total

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_max_fds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_network_receive_bytes_total

- **instance**: dummyapp:3000
- **job**: dummyapp

## process_network_transmit_bytes_total

- **instance**: dummyapp:3000
- **job**: dummyapp

## process_open_fds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_resident_memory_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_start_time_seconds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_virtual_memory_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## process_virtual_memory_max_bytes

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## promhttp_metric_handler_requests_in_flight

- **instance**: dummyapp:3000
- **job**: dummyapp

## promhttp_metric_handler_requests_total

- **code**: 200, 500, 503
- **instance**: dummyapp:3000
- **job**: dummyapp

## redis_cpu_time_seconds_total

- **state**: sys, sys_children, sys_main_thread, user, user_children, user_main_thread

## scrape_duration_seconds

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## scrape_samples_post_metric_relabeling

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## scrape_samples_scraped

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## scrape_series_added

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

## system_cpu_time_seconds_total

- **cpu**: cpu0, cpu1
- **state**: idle, interrupt, nice, softirq, steal, system, user, wait

## system_disk_io_bytes_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc
- **direction**: read, write

## system_disk_io_time_seconds_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc

## system_disk_merged_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc
- **direction**: read, write

## system_disk_operation_time_seconds_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc
- **direction**: read, write

## system_disk_operations_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc
- **direction**: read, write

## system_disk_pending_operations

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc

## system_disk_weighted_io_time_seconds_total

- **device**: vda, vda1, vda15, vda16, vdb, vdb1, vdc

## system_memory_usage_bytes

- **state**: buffered, cached, free, slab_reclaimable, slab_unreclaimable, used

## system_network_connections

- **protocol**: tcp
- **state**: CLOSE, CLOSE_WAIT, CLOSING, DELETE, ESTABLISHED, FIN_WAIT_1, FIN_WAIT_2, LAST_ACK, LISTEN, SYN_RECV, SYN_SENT, TIME_WAIT

## system_network_dropped_total

- **device**: eth0, lo
- **direction**: receive, transmit

## system_network_errors_total

- **device**: eth0, lo
- **direction**: receive, transmit

## system_network_io_bytes_total

- **device**: eth0, lo
- **direction**: receive, transmit

## system_network_packets_total

- **device**: eth0, lo
- **direction**: receive, transmit

## traefik_config_last_reload_success

- **instance**: traefik:8080
- **job**: traefik

## traefik_config_reloads_total

- **instance**: traefik:8080
- **job**: traefik

## traefik_entrypoint_request_duration_seconds_bucket

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **le**: +Inf, 0.1, 0.3, 1.2, 5.0
- **method**: GET
- **protocol**: http

## traefik_entrypoint_request_duration_seconds_count

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http

## traefik_entrypoint_request_duration_seconds_sum

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http

## traefik_entrypoint_requests_bytes_total

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http

## traefik_entrypoint_requests_total

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http

## traefik_entrypoint_responses_bytes_total

- **code**: 200
- **entrypoint**: web
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http

## traefik_open_connections

- **entrypoint**: metrics, web
- **instance**: traefik:8080
- **job**: traefik
- **protocol**: TCP

## traefik_router_request_duration_seconds_bucket

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **le**: +Inf, 0.1, 0.3, 1.2, 5.0
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_router_request_duration_seconds_count

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_router_request_duration_seconds_sum

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_router_requests_bytes_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_router_requests_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_router_responses_bytes_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **router**: whoami@file
- **service**: whoami@file

## traefik_service_request_duration_seconds_bucket

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **le**: +Inf, 0.1, 0.3, 1.2, 5.0
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## traefik_service_request_duration_seconds_count

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## traefik_service_request_duration_seconds_sum

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## traefik_service_requests_bytes_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## traefik_service_requests_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## traefik_service_responses_bytes_total

- **code**: 200
- **instance**: traefik:8080
- **job**: traefik
- **method**: GET
- **protocol**: http
- **service**: whoami@file

## up

- **instance**: dummyapp:3000, traefik:8080
- **job**: dummyapp, traefik

