# Worker Workflow

### `main.go`

### `func processJobs`

1. Receives message from worker queue with job data

2. Job execution entry is created in DB

3. Job code blocks are parsed

4. Sandbox reserved for user

5. Job blocks executed in Sandbox

6. Job results written to DB

7. Sandbox released back to pool
