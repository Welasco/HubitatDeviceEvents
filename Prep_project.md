# SolidJS
Start project
Ref: https://www.solidjs.com/guides/getting-started
npx degit solidjs/templates/ts
npm install

# tailwind
Ref: https://tailwindcss.com/docs/guides/solidjs
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# Runtime logging and shutdown

Logging goes to stdout by default, including when `logPath` is set. Leave
`logToFile=false` (or unset) in Kubernetes. To also write a local log file,
explicitly set `logToFile=true`; its name remains `<logPath><HOSTNAME>.log`.
Optional file logging is synchronous and can block on an unhealthy filesystem.
An invalid `logToFile` value or an inaccessible enabled log file fails startup
instead of silently falling back.

The container startup script uses `exec` so the Go application receives SIGTERM
directly. SIGINT/SIGTERM stop new connections and allow up to 10 seconds for
in-flight requests to drain before canceling remaining dependency operations and
closing MySQL. The whole shutdown, including logging and database cleanup, has a
15-second application exit budget; timeout or cleanup failure exits nonzero.
Cancellation is best effort: a timed-out database write can have an uncertain
commit outcome, and forced exit cannot guarantee every received event is stored.

Each request has a shared 5-second dependency deadline, including the complete
unknown-device lookup, Hubitat fetch, device insert, and event retry. Database
operations and outbound HTTP calls also enforce a maximum of 5 seconds. MySQL
connection/read/write timeouts are capped at 5 seconds, preserving shorter DSN
values. HTTP read/write/idle timeouts are 5/10/30 seconds. Query iteration failures
return errors rather than partial success; canceled and expired operations return
503 and 504 respectively. Database initialization errors stop startup.

Keep the pod's `terminationGracePeriodSeconds` above 15 seconds plus any
`preStop` hook time and an operational margin (30 seconds is the Kubernetes
default). Stop this workload before shutting down storage it depends on.
Rebuild and deploy the image to apply these changes. Verify SIGTERM while idle,
during event ingestion, and with MySQL/Hubitat unavailable before rebooting a node.

These bounds cannot release Linux threads stuck in uninterruptible kernel I/O.
If node shutdown still reports hung tasks or filesystem-sync timeouts, capture
the complete kernel call traces and investigate the node's filesystem, container
storage, mounted volumes, and storage shutdown ordering.

Run backend regression tests with `go test ./...`.
