# singlemod

Run

go run main.go serve -p 8080


## Project Structure

---

### **1. .github/workflows**
This directory manages GitHub Actions, providing automated workflows for continuous integration (CI), continuous deployment (CD), and other GitHub event-triggered tasks. Developers can define various workflows to run tests, build binaries, deploy applications, and more.

---

### **2. .vscode**
Holds configuration files for the Visual Studio Code editor, ensuring a consistent development environment for all contributors. Developers may find settings and recommendations for extensions that are conducive to the project’s development.

---

### **3. cmd**
The entry point for the application or any related command-line interfaces (CLI). These scripts initialize and run the application, utilizing the Cobra CLI library. Developers should define CLI commands and flags in this directory.

---

### **4. internals**
Dedicated to housing the core application logic, organized into various segments:

- **core**
  - **entity**: Holds domain entities, which represent primary data structures and related functionalities.
  - **serr**: Contains definitions and potentially, handling logic for server-specific errors.
  
- **http**
  - **handler**: Responsible for handling HTTP requests and responses, essentially controlling the flow of HTTP traffic.
  - **helpers**: A collection of helper functions and utilities that assist with HTTP-related logic and functionality.
  - **transport**: Manages the transport layer of HTTP, handling the payload data transmission between client and server.
  
- **service**: Contains the service layer, encapsulating business logic and dictating how data is processed and handled within the application.
  
- **storage**: Manages the storage layer, which is responsible for data persistence and retrieval.

---

### **5. server**
The server directory encompasses various elements related to the server-side of the application:

- **infra**: Incorporates the infrastructure layer, housing configurations, constants, and shared logic utilized throughout the application.

- **middleware**: Contains middleware components that process HTTP requests and responses in between client interaction and reaching the application's handler or route.

- **routing**: Manages the routing of the server, defining paths, associating handlers, and ensuring that the HTTP request is adhered to the correct logic path.

--- 

For testing, [mockery](https://github.com/vektra/mockery) is reccomended.