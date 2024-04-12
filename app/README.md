This project is aimed to provide a simple and easy to use data storage API for tasks.
It is written in Go and uses PostgreSQL as the database.
The API is designed to be RESTful and ready to be consumed by any client.

The following endpoints are available:
- `GET /api/tasks?page=1&size=10`: Returns all tasks in the database.
- `GET /api/tasks/{id}`: Returns the task with the given ID.
- `POST /api/tasks`: Creates a new task.
- `PUT /api/tasks/{id}`: Updates the task with the given ID.
- `DELETE /api/tasks/{id}`: Deletes the task with the given ID.

Considering the host and port of the server is `localhost:8080`, an example request to create a new task would look like this:
```bash
curl -X GET http://localhost:8080/api/tasks
```
> Note: Some endpoints require CSRF header to be set. For example, to create a new task, you need to set the `X-CSRF-Token` header with a valid CSRF token. So you might need a more capable tool like Postman that supports cookie storage to test these endpoints.

See `src/api/routers/task_router.go` for more details.


For more detailed documentation, see the following URL after running `godoc -http=127.0.0.1:6060` command in this directory (`./app`)

http://localhost:6060/pkg/github.com/emso-c/konzek-go-assignment/