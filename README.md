# HTTP Fan-out

A very small and simple application to fan-out HTTP requests to multiple servers. Written to deliver webhooks in my smart home to both my "production" machine and my development machine, but implemented in a very generic way.

## Configuration

Configuration is done through environment variables. The following variables are used:

- `BIND`: The address to bind to. Defaults to an empty string, which means all interfaces.
- `PORT`: The port to listen on. Defaults to `80`.
- `TARGET_MAIN`: The main target to fan-out to. Optional. If provided, the response from this target will be returned to the client. If not provided, the response will always be an empty `200 OK`.
- `TARGET_{N}`: A target to fan-out to. The `{N}` is a number starting at either `0` or `1`. The value is a fully qualified URL. Targets have to be numbered sequentially.
