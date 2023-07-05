# Project Structure

[back to main](../../README.md)

## How the project structure is organized

The structure follows a Technical Design Grouping. Let's explore the directories and their responsibilities in more detail:

- `cmd`: This directory contains the application's entry point(s). The files in this directory are responsible for initializing and starting your application. or any cli commands. ( using [cobra](https://github.com/spf13/cobra) cli )

- `docs`: This directory contains documentation for the project. Any additional resources like images used in the documentation are located in the `img` subdirectory.

- `mocks`: Contains mock functions ( using [mockery](https://github.com/vektra/mockery) )

- `pkg`: This directory includes the main application code.

  - `entities`: This subdirectory contains your domain entities, which represent the Primary Data structures and related functions.
  
  - `http`: The code related to HTTP resides here. The subdirectories are:
  
    - `handlers`: Code that handles HTTP requests and responses.
    
    - `transport`: Code related HTTP payload transport, operations with request and response body.
    
    - `validator`: This directory contains code for validating HTTP request payloads.
    
  - `infra`: This subdirectory contains the infrastructure layer code, like configurations shared across the application.
  
    - `constants`: Contains constant values that are shared across your project.
  
  - `services`: This directory contains the service layer, which contains business logic.
  
    - `helpers`: Contains helper functions that can be utilized across different services.
  
  - `storage`: This subdirectory contains code related to storage, like database or file system interactions.
  
    - `sqlite`: Specific implementations for SQLite storage.
  
  - `svrerr`: This directory includes server error types definitions for whole application.

- `scripts`: This directory holds scripts for tasks like building, linting, testing, or deployment.

- `server`: The directory typically contains the main server code, server initialization, binding handlers, service, storage layers ( dependancy injection ) and adding routes and middleware.

