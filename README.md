# distributed-job-manager

## Setup

Pre-requisites:

- docker
- Golang >=1.24

Steps:

1. `git clone https://github.com/julianstephens/distributed-job-manager.git`
2. `docker compose up -d`

## TODO

### backend

- [ ] Update DB table(s) with worker results
- [ ] Assign worker threads by user id
- [ ] Create manager service
  - [ ] Add /register endpoint to initialize worker
  - [ ] Add heartbeat monitoring
  - [ ] Add worker cleanup
- [ ] Create coordinator service
- [ ] Add job cancellation endpoint

### frontend

- [x] Add timestamps to tables
