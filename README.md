# Rest-Api-Task-List
A Basic Rest API written in GoLang with CRUD operations for List and its Tasks

This repository implements the below functionality:

A REST API for a TO-DO application in Golang, which will have
the following features:

- Create a task list
- Add task into a list
- Delete task from a list
- Update an existing task
- Delete a list

Added a MySQL persistence mechanism to save tasks and lists.

- Create a Task List
  - Endpoint - /list
  - Method - POST
  - Body - {"name": "List1"}

- Add Task into a List
  - Endpoint - /list/:list_id/task
  - Method - POST
  - Body - {"description": "Sample Task 1"}

- Delete Task from a List
  - Endpoint - /list/:list_id/task/:task_id
  - Method - DELETE

- Update an Existing Task
  - Endpoint - /task/:task_id
  - Method - PUT
  - Body - {"list_id": 1, "description": "Updated Sample Task 1"}

- Delete a List (MySQL ON DELETE CASCADE deletes all the tasks in the deleted list as well)
  - Endpoint - /list/:list_id
  - Method - DELETE

The Postman collection to test the APIs can be found here: [Postman Collection For Task List] (https://www.getpostman.com/collections/ad4437e76861ce09d219)
