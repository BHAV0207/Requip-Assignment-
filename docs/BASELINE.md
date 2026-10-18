# Baseline record

| | |
|---|---|
| Student | Bhavya Jain |
| Student ID | _fill in_ |
| Project | Requip User Management, Advanced DevOps final project |
| Repository | https://github.com/BHAV0207/Requip-Assignment- |
| Baseline commit | [`59da7edbdad65f642915d796a8d33fe806830318`](https://github.com/BHAV0207/Requip-Assignment-/commit/59da7edbdad65f642915d796a8d33fe806830318) (also tagged `baseline`) |

## Application

A user-management system. A React 19 + TypeScript single-page app calls a Go (Gin + GORM) REST API that stores users in MySQL 8.4. It supports create, list (paginated), read, update and soft-delete, with validation for Aadhaar, PAN and mobile numbers.

## What already existed at the baseline

| Area | State at baseline |
|---|---|
| Tests | Unit tests for the service and handler layers using mocks |
| Containers | `Backend/docker-compose.yml` ran **only MySQL**, with passwords hard-coded in the file. No Dockerfile for the API or the frontend |
| CI/CD | None |
| Kubernetes / Helm | None |
| Health checks | A single `/health` endpoint that always returned 200 and never checked the database |
| Observability | Gin's default text access log only. No metrics, no tracing |
| Configuration | `.env` file. The default DB port was wrong (5432, the Postgres port, for a MySQL app) |
| Security | Hard-coded DB credentials in compose, CORS allowing every origin, no dependency or image scanning |

Everything else in the pull request was added on top of this commit.
