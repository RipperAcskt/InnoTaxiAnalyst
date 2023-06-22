# InnoTaxiDriver API

This is a microservice for analyse service.

## Installation

To install and run the service, follow these steps:

Install Go on your system if you haven't already done so.
Clone the repository to your local machine.
Run the following command in the root directory of the project:

    go run ./cmd/main.go

Also you can run project using docker-compose.

## Starting Dependency using Docker and Docker Compose

To start the dependency using Docker and Docker Compose, follow these steps:

- Make sure you have Docker and Docker Compose installed on your system.
- Open a terminal and navigate to the root directory of the project.
- In the terminal, run the following command to start the dependency:

        docker-compose up -d

    This command will start the dependency service in detached mode, allowing it to run in the background.
    You can now access the dependency from your application using the specified configuration (e.g., hostname, port) or environment variables.

    Note: Make sure your application is configured to connect to the dependency service using the appropriate hostname and port.

- To stop the dependency service, run the following command in the terminal:
  
      docker-compose down

    This command will stop and remove the running dependency container.

# Commit Rules

When making commits to the repository, please follow the commit rules outlined below

## Types
Use one of the following types to categorize your commit:

- API: Relevant changes to the API.
- feat: Commits that add a new feature.
- fix: Commits that fix a bug.
- refactor: Commits that rewrite/restructure your code without changing any behavior.
- perf: Commits that improve performance (special refactor commits).
- style: Commits that do not affect the meaning of the code (e.g., whitespace, formatting).
- test: Commits that add missing tests or correct existing tests.
- docs: Commits that affect documentation only.
- build: Commits that affect build components like build tools, CI pipeline, dependencies, project version, etc.
- ops: Commits that affect operational components like infrastructure, deployment, backup, recovery, etc.
- chore: Miscellaneous commits (e.g., modifying .gitignore).
## Scopes

Use an optional scope to provide additional contextual information. Allowed scopes depend on the specific project. Do not use issue identifiers as scopes.

## Subject
The subject should contain a succinct description of the change. Follow these guidelines:

The subject is a mandatory part of the commit format.
Use the imperative, present tense (e.g., "change" instead of "changed" or "changes").
Think of the commit as "This commit will <subject>".
Do not capitalize the first letter.
Do not include a dot (.) at the end.
## Body
The body is an optional part of the commit format. Include the motivation for the change and contrast it with the previous behavior. Follow these guidelines:

Use the imperative, present tense (e.g., "change" instead of "changed" or "changes").
This is the place to mention issue identifiers and their relations.
## Footer
The footer is an optional part of the commit format. Use it to contain any information about breaking changes and reference related issues. Follow these guidelines:

Optionally reference an issue by its ID.
Breaking changes should start with the phrase "BREAKING CHANGES:" followed by a space or two newlines. The rest of the commit message can be used to describe the breaking changes.
## Examples
Here are some examples of commit messages following the commit rules:

    feat(shopping cart): add the amazing button


    feat: remove ticket list endpoint

    Refers to JIRA-1337
    

    fix: add missing parameter to service call

    The error occurred because of <reasons>.


    build(release): bump version to 1.0.0


    build: update dependencies


    refactor: implement calculation method as recursion

    
    style: remove empty line

## Code Description

# Project structure

The code for the microservice is organized into several packages:

- cmd/main.go contains the main function for the service.
- internal/app/app.go contains functions which sets up the API routes and starts the server.
- internal/models/ contains the data models for the application. In this case, there is only one model - Driver.
- internal/repo/ contains the repository implementation for working with the databases.
- internal/services/ contains the business logic services for the application.
- internal/handlers/ contains the API request handlers for the application. Service provides handlers for registartion and auth driver, also handlers for working with driver's profile.

---

- Business logic services are located in the internal/services/ package. Service uses repository layer to get data.

### Conclusion

This microservice demonstrates a simple way to analyse and follow after service using Go. The code is organized into packages, making it easy to maintain and extend. The `AnalystService` and `AnalystHandler` objects provide the business logic and API endpoints, respectively.