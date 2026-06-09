# Superfluous WriteHeader from gzipResponseWriter resolve

The SSE issue was caused by the `/my/logs` route being registered in a way that it still received the `GzipMiddleware`. In the previous attempt, the developer tried to fix this by creating a new route group with the same prefix (`/my`) as the existing group. However, in many router implementations, creating multiple groups with the same prefix can lead to them being merged or the middlewares being applied to all routes under that prefix.

I have resolved this by:
1. Creating a single base group for `/my` with the `auth` middleware.
2. Registering the SSE `/logs` route directly on this base group, so it does not receive the `GzipMiddleware`.
3. Creating a sub-group of the `/my` group specifically for the other routes, and applying the `GzipMiddleware` only to this sub-group.

This ensures that `/my/logs` is served without gzip, while all other `/my/*` routes continue to benefit from gzip compression.
