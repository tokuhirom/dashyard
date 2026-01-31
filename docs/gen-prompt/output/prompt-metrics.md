# Available Metrics

## go_gc

- `go_gc_duration_seconds` (summary) — A summary of the pause duration of garbage collection cycles.
  Labels: instance (1 values), job (1 values), quantile (5 values)
- `go_gc_duration_seconds_count`
  Labels: instance (1 values), job (1 values)
- `go_gc_duration_seconds_sum`
  Labels: instance (1 values), job (1 values)

## go_goroutines

- `go_goroutines` (gauge) — Number of goroutines that currently exist.
  Labels: instance (1 values), job (1 values)

## go_info

- `go_info` (gauge) — Information about the Go environment.
  Labels: instance (1 values), job (1 values), version (1 values)

## go_memstats

- `go_memstats_alloc_bytes` (gauge) — Number of bytes allocated and still in use.
  Labels: instance (1 values), job (1 values)
- `go_memstats_alloc_bytes_total` (counter) — Total number of bytes allocated, even if freed.
  Labels: instance (1 values), job (1 values)
- `go_memstats_buck_hash_sys_bytes` (gauge) — Number of bytes used by the profiling bucket hash table.
  Labels: instance (1 values), job (1 values)
- `go_memstats_frees_total` (counter) — Total number of frees.
  Labels: instance (1 values), job (1 values)
- `go_memstats_gc_sys_bytes` (gauge) — Number of bytes used for garbage collection system metadata.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_alloc_bytes` (gauge) — Number of heap bytes allocated and still in use.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_idle_bytes` (gauge) — Number of heap bytes waiting to be used.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_inuse_bytes` (gauge) — Number of heap bytes that are in use.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_objects` (gauge) — Number of allocated objects.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_released_bytes` (gauge) — Number of heap bytes released to OS.
  Labels: instance (1 values), job (1 values)
- `go_memstats_heap_sys_bytes` (gauge) — Number of heap bytes obtained from system.
  Labels: instance (1 values), job (1 values)
- `go_memstats_last_gc_time_seconds` (gauge) — Number of seconds since 1970 of last garbage collection.
  Labels: instance (1 values), job (1 values)
- `go_memstats_lookups_total` (counter) — Total number of pointer lookups.
  Labels: instance (1 values), job (1 values)
- `go_memstats_mallocs_total` (counter) — Total number of mallocs.
  Labels: instance (1 values), job (1 values)
- `go_memstats_mcache_inuse_bytes` (gauge) — Number of bytes in use by mcache structures.
  Labels: instance (1 values), job (1 values)
- `go_memstats_mcache_sys_bytes` (gauge) — Number of bytes used for mcache structures obtained from system.
  Labels: instance (1 values), job (1 values)
- `go_memstats_mspan_inuse_bytes` (gauge) — Number of bytes in use by mspan structures.
  Labels: instance (1 values), job (1 values)
- `go_memstats_mspan_sys_bytes` (gauge) — Number of bytes used for mspan structures obtained from system.
  Labels: instance (1 values), job (1 values)
- `go_memstats_next_gc_bytes` (gauge) — Number of heap bytes when next garbage collection will take place.
  Labels: instance (1 values), job (1 values)
- `go_memstats_other_sys_bytes` (gauge) — Number of bytes used for other system allocations.
  Labels: instance (1 values), job (1 values)
- `go_memstats_stack_inuse_bytes` (gauge) — Number of bytes in use by the stack allocator.
  Labels: instance (1 values), job (1 values)
- `go_memstats_stack_sys_bytes` (gauge) — Number of bytes obtained from system for stack allocator.
  Labels: instance (1 values), job (1 values)
- `go_memstats_sys_bytes` (gauge) — Number of bytes obtained from system.
  Labels: instance (1 values), job (1 values)

## go_threads

- `go_threads` (gauge) — Number of OS threads created.
  Labels: instance (1 values), job (1 values)

## process_cpu

- `process_cpu_seconds_total` (counter) — Total user and system CPU time spent in seconds.
  Labels: instance (1 values), job (1 values)

## process_max

- `process_max_fds` (gauge) — Maximum number of open file descriptors.
  Labels: instance (1 values), job (1 values)

## process_open

- `process_open_fds` (gauge) — Number of open file descriptors.
  Labels: instance (1 values), job (1 values)

## process_resident

- `process_resident_memory_bytes` (gauge) — Resident memory size in bytes.
  Labels: instance (1 values), job (1 values)

## process_start

- `process_start_time_seconds` (gauge) — Start time of the process since unix epoch in seconds.
  Labels: instance (1 values), job (1 values)

## process_virtual

- `process_virtual_memory_bytes` (gauge) — Virtual memory size in bytes.
  Labels: instance (1 values), job (1 values)
- `process_virtual_memory_max_bytes` (gauge) — Maximum amount of virtual memory available in bytes.
  Labels: instance (1 values), job (1 values)

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
  Labels: instance (1 values), job (1 values)

## scrape_samples

- `scrape_samples_post_metric_relabeling`
  Labels: instance (1 values), job (1 values)
- `scrape_samples_scraped`
  Labels: instance (1 values), job (1 values)

## scrape_series

- `scrape_series_added`
  Labels: instance (1 values), job (1 values)

## system_cpu

- `system_cpu_time_seconds_total`
  Labels: cpu (12 values), state (8 values)

## system_disk

- `system_disk_io_bytes_total`
  Labels: device (10 values), direction (2 values)
- `system_disk_io_time_seconds_total`
  Labels: device (10 values)
- `system_disk_merged_total`
  Labels: device (10 values), direction (2 values)
- `system_disk_operation_time_seconds_total`
  Labels: device (10 values), direction (2 values)
- `system_disk_operations_total`
  Labels: device (10 values), direction (2 values)
- `system_disk_pending_operations`
  Labels: device (10 values)
- `system_disk_weighted_io_time_seconds_total`
  Labels: device (10 values)

## system_memory

- `system_memory_usage_bytes`
  Labels: state (6 values)

## system_network

- `system_network_connections`
  Labels: protocol (1 values), state (12 values)
- `system_network_dropped_total`
  Labels: device (2 values), direction (2 values)
- `system_network_errors_total`
  Labels: device (2 values), direction (2 values)
- `system_network_io_bytes_total`
  Labels: device (2 values), direction (2 values)
- `system_network_packets_total`
  Labels: device (2 values), direction (2 values)

## traefik_config

- `traefik_config_last_reload_success` (gauge) — Last config reload success
  Labels: instance (1 values), job (1 values)
- `traefik_config_reloads_total` (counter) — Config reloads
  Labels: instance (1 values), job (1 values)

## traefik_entrypoint

- `traefik_entrypoint_request_duration_seconds_bucket`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values)
- `traefik_entrypoint_request_duration_seconds_count`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
- `traefik_entrypoint_request_duration_seconds_sum`
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
- `traefik_entrypoint_requests_bytes_total` (counter) — The total size of requests in bytes handled by an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
- `traefik_entrypoint_requests_total` (counter) — How many HTTP requests processed on an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)
- `traefik_entrypoint_responses_bytes_total` (counter) — The total size of responses in bytes handled by an entrypoint, partitioned by status code, protocol, and method.
  Labels: code (1 values), entrypoint (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values)

## traefik_open

- `traefik_open_connections` (gauge) — How many open connections exist, by entryPoint and protocol
  Labels: entrypoint (2 values), instance (1 values), job (1 values), protocol (1 values)

## traefik_router

- `traefik_router_request_duration_seconds_bucket`
  Labels: code (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
- `traefik_router_request_duration_seconds_count`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
- `traefik_router_request_duration_seconds_sum`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
- `traefik_router_requests_bytes_total` (counter) — The total size of requests in bytes handled by a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
- `traefik_router_requests_total` (counter) — How many HTTP requests are processed on a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)
- `traefik_router_responses_bytes_total` (counter) — The total size of responses in bytes handled by a router, partitioned by service, status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), router (1 values), service (1 values)

## traefik_service

- `traefik_service_request_duration_seconds_bucket`
  Labels: code (1 values), instance (1 values), job (1 values), le (5 values), method (1 values), protocol (1 values), service (1 values)
- `traefik_service_request_duration_seconds_count`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
- `traefik_service_request_duration_seconds_sum`
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
- `traefik_service_requests_bytes_total` (counter) — The total size of requests in bytes received by a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
- `traefik_service_requests_total` (counter) — How many HTTP requests processed on a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)
- `traefik_service_responses_bytes_total` (counter) — The total size of responses in bytes returned by a service, partitioned by status code, protocol, and method.
  Labels: code (1 values), instance (1 values), job (1 values), method (1 values), protocol (1 values), service (1 values)

## up

- `up`
  Labels: instance (1 values), job (1 values)

# Label Values

Full label value listings for each metric.

## go_gc_duration_seconds

- **instance**: traefik:8080
- **job**: traefik
- **quantile**: 0.0, 0.25, 0.5, 0.75, 1.0

## go_gc_duration_seconds_count

- **instance**: traefik:8080
- **job**: traefik

## go_gc_duration_seconds_sum

- **instance**: traefik:8080
- **job**: traefik

## go_goroutines

- **instance**: traefik:8080
- **job**: traefik

## go_info

- **instance**: traefik:8080
- **job**: traefik
- **version**: go1.24.5

## go_memstats_alloc_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_alloc_bytes_total

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_buck_hash_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_frees_total

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_gc_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_alloc_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_idle_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_inuse_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_objects

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_released_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_heap_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_last_gc_time_seconds

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_lookups_total

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mallocs_total

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mcache_inuse_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mcache_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mspan_inuse_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_mspan_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_next_gc_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_other_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_stack_inuse_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_stack_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_memstats_sys_bytes

- **instance**: traefik:8080
- **job**: traefik

## go_threads

- **instance**: traefik:8080
- **job**: traefik

## process_cpu_seconds_total

- **instance**: traefik:8080
- **job**: traefik

## process_max_fds

- **instance**: traefik:8080
- **job**: traefik

## process_open_fds

- **instance**: traefik:8080
- **job**: traefik

## process_resident_memory_bytes

- **instance**: traefik:8080
- **job**: traefik

## process_start_time_seconds

- **instance**: traefik:8080
- **job**: traefik

## process_virtual_memory_bytes

- **instance**: traefik:8080
- **job**: traefik

## process_virtual_memory_max_bytes

- **instance**: traefik:8080
- **job**: traefik

## redis_cpu_time_seconds_total

- **state**: sys, sys_children, sys_main_thread, user, user_children, user_main_thread

## scrape_duration_seconds

- **instance**: traefik:8080
- **job**: traefik

## scrape_samples_post_metric_relabeling

- **instance**: traefik:8080
- **job**: traefik

## scrape_samples_scraped

- **instance**: traefik:8080
- **job**: traefik

## scrape_series_added

- **instance**: traefik:8080
- **job**: traefik

## system_cpu_time_seconds_total

- **cpu**: cpu0, cpu1, cpu10, cpu11, cpu2, cpu3, cpu4, cpu5, cpu6, cpu7, cpu8, cpu9
- **state**: idle, interrupt, nice, softirq, steal, system, user, wait

## system_disk_io_bytes_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0
- **direction**: read, write

## system_disk_io_time_seconds_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0

## system_disk_merged_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0
- **direction**: read, write

## system_disk_operation_time_seconds_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0
- **direction**: read, write

## system_disk_operations_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0
- **direction**: read, write

## system_disk_pending_operations

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0

## system_disk_weighted_io_time_seconds_total

- **device**: dm-0, dm-1, dm-2, loop0, nvme0n1, nvme0n1p1, nvme0n1p2, nvme0n1p3, nvme0n1p4, zram0

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

- **instance**: traefik:8080
- **job**: traefik

