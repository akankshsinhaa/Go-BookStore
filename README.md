
# About

Robust online book store to handle user authentication, authorization, and access management.

## Project Overview:

The project is a web application written in Go (Golang) that serves as a platform for managing and reviewing books. It provides features for user registration, authentication, adding and viewing books, as well as adding and viewing book reviews. The application uses a PostgreSQL database for data storage. Below, I'll explain the key aspects of the project:

Planned Structure:<br>
![diagram](https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/f66a8f9c-2555-41f1-b808-bc894c266b09)



**Security Measures:**  

1.  **User Authentication:** The project uses session-based user authentication to securely manage user sessions. It leverages the "gorilla/sessions" library to store session data securely, ensuring that user data remains confidential and protected.

  

2.  **Password Hashing:** When users register, their passwords are securely hashed using the bcrypt hashing algorithm. This adds an additional layer of security by storing only the hash of the password in the database, making it challenging for attackers to obtain plaintext passwords.

  

3.  **SQL Injection Prevention:** To prevent SQL injection attacks, the project uses parameterized queries when interacting with the PostgreSQL database. This ensures that user input is properly sanitized and eliminates the risk of SQL injection vulnerabilities.

  

4.  **Access Control:** The application implements access control to restrict certain functionalities to authenticated users. Unauthorized access to sensitive areas of the application is prevented through proper authentication and authorization checks.

  

5.  **Error Handling:** Comprehensive error handling is employed throughout the application to handle unexpected situations gracefully, without revealing sensitive information to users. Error messages are appropriately logged for debugging purposes.

  
**Libraries and Implementations:**
 

1.  **gorilla/mux:** The "gorilla/mux" library is used for routing and URL handling. It allows for creating flexible and RESTful routes for different application functionalities.

  

2.  **gorilla/sessions:** The "gorilla/sessions" library is utilized for session management. It helps store and manage user session data securely, including user authentication information.

  

3.  **database/sql:** The standard "database/sql" package is used for interacting with the PostgreSQL database. It provides a safe and efficient way to perform database operations.

  

4.  **github.com/lib/pq:** The "github.com/lib/pq" package is the PostgreSQL driver for Go, enabling seamless interaction with the PostgreSQL database.

  

5.  **html/template:** The built-in "html/template" package is used for rendering HTML templates. It facilitates the separation of logic and presentation in web pages.

  

6.  **golang.org/x/crypto/bcrypt:** The "golang.org/x/crypto/bcrypt" package is used for secure password hashing. It ensures that user passwords are hashed using a strong and salted hashing algorithm.

  

7.  **Frontend (HTML/CSS):** The project uses HTML and CSS for the frontend to create user-friendly interfaces for registration, login, book management, and review submissions.


## Development To-Do

  

### User Interface

- [x] User Login Page

- [x] User Login Route

- [x] User Register Page

- [x] User Register Route

- [x] All Books Page

- [ ] Search Functionality

- [x] All Books Route

- [x] Admin Page

- [x] Admin Route

- [x] Add Review Page

- [ ] Add Review Route

- [ ] View Review Page

- [x] View Review Route

- [ ] Buy/Download Route

- [ ] Buy/Download Page

  

### Database To-Do

- [x] Setup Postgres User

- [x] Setup Postgres Books

- [x] Create Table User

- [x] Create Table Books

  

### Security To-Do

- [ ] Input Sanitization using go-sanitize

- [x] Logging

- [x] Hashing Password

- [x] SSL Mode

  

### Docker To-Do

- [ ] Postgres

- [x] Reverse Proxy (NGINX)

  

### Miscellaneous To-Do

- [ ] Unit-Testing

- [x] Middleware

- [x] Comments


  

# Setup

  

## Postgres

Details to serve as follows:-

  

- Hostname: localhost

- Username: app_user

- Password: app_password

- Databasename: app_database

- Portnumber: 5432

- SSL Mode: enable

  

### User Table

```sql

CREATE  TABLE  IF  NOT  EXISTS bcrypt (

UserID SERIAL  PRIMARY KEY,

Username VARCHAR(50) UNIQUE  NOT NULL,

Hash  TEXT  NOT NULL

);

```

  
  
  
  

### Books Table

```sql

CREATE  TABLE  books (

id SERIAL  PRIMARY KEY,

title TEXT,

author TEXT,

publisher TEXT,

reviews TEXT[]

);

```

  

### Add Dummy Data

```sql

INSERT INTO books (title, author, publisher, reviews) VALUES

('Book 1', 'Author 1', 'Publisher A', '{"Review 1 for Book 1", "Review 2 for Book 1"}'),

('Book 2', 'Author 2', 'Publisher B', '{"Review 1 for Book 2", "Review 2 for Book 2"}'),

('Book 3', 'Author 3', 'Publisher C', '{"Review 1 for Book 3", "Review 2 for Book 3"}');

```

## Go Setup

```go

go get <<package name>>(mentioned in mod file)

```

  

## Running The App

To run the application locally use the following command in terminal.

```go

go run main.go

```



<em>Web app should be available at localhost:8088/register</em>

### Docker Reverse Proxy Setup(NGINX)

```zsh
sudo docker run -d --name nginx-base -p 80:80 nginx:latest
```

Now, replace file at nginx-base:/etc/nginx/conf.d/default.conf with the one in the repo.

```zsh
sudo docker exec nginx-base nginx -t
sudo docker exec nginx-base nginx -s reload
sudo docker commit nginx-base nginx-proxy
```

Ensure that the Dockerfile resides within the same directory as the command terminal's current path.

```zsh
sudo docker build -t nginx-reverse-proxy .
```

<em>Web app should now be available at http://localhost</em>

# Screenshots

<img width="1280" alt="Screenshot 2023-09-02 at 7 46 52 PM" src="https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/09a7a05a-4737-4c5e-a091-5bee0208e6be">
<img width="1280" alt="Screenshot 2023-09-02 at 7 46 17 PM" src="https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/8d578d10-5261-4177-a21f-4aab2cd4b66f">
<img width="1280" alt="Screenshot 2023-09-02 at 7 42 41 PM" src="https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/62c44ad9-7c1c-452b-9d2d-5a0d0d57bcfb">
<img width="1273" alt="Screenshot 2023-09-02 at 7 41 04 PM" src="https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/227ffd9d-f712-439d-929f-e161ad6cc942">
<img width="1279" alt="Screenshot 2023-09-02 at 7 41 41 PM" src="https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-akankshsinhaa/assets/91271963/77f9cd18-d3fd-4ec5-aab6-e8919d46a89b">


 
