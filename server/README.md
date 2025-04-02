# Environment Variables

- The service expects a fully-formed `DATABASE_URL` variable exposed to the process, which is a Postgres DB connection string.

- `API_KEY` should be set to an agreed-upon value, to be sent in the header `x-api-key` with POST requests.

# TODOs
Presently API Key is a single-tenant solution, which will eventually be revised 
