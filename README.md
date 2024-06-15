# Snippetbox

## About

**Snippetbox** is a project developed along with the book "Let's Go" by Alex Edwards. It is a web application built using the Go programming language and follows best practices in web development. The project covers various aspects of building a web application, such as routing, middleware, database integration, and authentication.

## Table of Contents

- [Demo](#demo)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Routes](#routes)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

## Demo
![](https://github.com/RafaelClaumann/snippetbox/blob/main/demo_snippetbox.gif)

## Features

- User authentication and authorization
- Create, view, and manage code snippets
- Secure password handling
- Middleware implementation for enhanced functionality
- Database integration with MySQL
- Graceful error handling and logging

## Installation

To set up Snippetbox locally, follow these steps:

1. **Clone the repository:**
    ```bash
    git clone https://github.com/RafaelClaumann/snippetbox.git
    cd snippetbox
    ```

2. **Run the MySQL container using Docker Compose:**

   Ensure you have Docker and Docker Compose installed on your machine.

   ```bash
   docker-compose up -d
   ```

   This command will start the MySQL database in detached mode (`-d`).

3. **Build and run the application:**
    ```bash
    go build ./cmd/web
    ./web

    # or

    go run ./cmd/web
    ```

4. **To run integrated tests you need to create 'test_snippetbox' database:**
    ```bash
    CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    ```

5. **Running application tests:**
    ```bash
    go test -v ./...
    ```    

## Usage

Once the application is running, you can access it via `https://localhost:4000`. The following routes are available:

## Routes

| Method | URL                        | Function                | Description                                                           | Access   |
|--------|----------------------------|-------------------------|-----------------------------------------------------------------------|----------|
| `GET`  | `/`                        | `app.home`              | Displays the home page of the application.                            | Public   |
| `GET`  | `/about`                   | `app.about`             | Displays the "About" page of the application.                         | Public   |
| `GET`  | `/snippet/view/:id`        | `app.snippetView`       | Displays a specific snippet based on the provided ID in the URL.      | Public   |
| `GET`  | `/user/signup`             | `app.userSignup`        | Displays the user signup form.                                        | Public   |
| `POST` | `/user/signup`             | `app.userSignupPost`    | Processes the user signup form data.                                  | Public   |
| `GET`  | `/user/login`              | `app.userLogin`         | Displays the user login form.                                         | Public   |
| `POST` | `/user/login`              | `app.userLoginPost`     | Processes the user login form data.                                   | Public   |
| `GET`  | `/account/password/update` | `app.updatePassword`    | Displays the form to update the user's password.                      | Protected|
| `POST` | `/account/password/update` | `app.updatePasswordPost`| Processes the password update form data.                              | Protected|
| `GET`  | `/account/view`            | `app.accountView`       | Displays the user's account information.                              | Protected|
| `GET`  | `/snippet/create`          | `app.snippetCreate`     | Displays the form to create a new snippet.                            | Protected|
| `POST` | `/snippet/create`          | `app.snippetCreatePost` | Processes the new snippet creation form data.                         | Protected|
| `POST` | `/user/logout`             | `app.userLogoutPost`    | Processes the user logout, ending the session.                        | Protected|

## Technologies Used

- **Go**: The main programming language used for the backend.
- **MySQL**: Used for storing user and snippet data.
- **HTML/CSS/JavaScript**: For the front-end interface.
- **Docker**: To facilitate containerization and deployment.
- **GitHub Actions**: To create CI(Continuous Integration).

## Contributing

Contributions are welcome! If you have any improvements or new features you'd like to add, please fork the repository, create a new branch, and submit a pull request. Ensure your code follows the existing style and includes tests where appropriate.

---

For more information, visit the [repository](https://github.com/RafaelClaumann/snippetbox) on GitHub.
