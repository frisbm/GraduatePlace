# Graduate Place
#### A place for graduate students (soon)

## Technology
* Backend
  * [go-chi](https://github.com/go-chi/chi) - routing
  * [sqlc](https://sqlc.dev/) - sql-go code generation
  * [asynq](https://github.com/hibiken/asynq) - async task queue
  * [postgresql](https://www.postgresql.org/) - database
  * [redis](https://redis.io/) - message broker
  * [air](https://github.com/cosmtrek/air) - hot reloading

## Running Locally

It is easiest to use the provided [docker-compose.yaml](docker-compose.yaml) file to run everything locally.
The servers run with air so any changes made will hot reload for a streamlined development experience.
If you like Makefiles like I do, feel free to use/peruse the make targets to assist.
