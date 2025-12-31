# 📝 Task Tracker API

A robust RESTful API built with **Golang**, **Gin Framework**, and **PostgreSQL**.
This project demonstrates a complete backend service with Authentication, CRUD operations, Pagination, and API Documentation.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin_Framework-v1.9-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-v1.25-red?style=for-the-badge)
![Swagger](https://img.shields.io/badge/Swagger-UI-85EA2D?style=for-the-badge&logo=swagger)

## 🚀 Features

- **Authentication**: Secure User Registration & Login using **JWT** (JSON Web Token) and **Bcrypt**.
- **Task Management**: Create, Read, Update, and Delete (CRUD) tasks.
- **Data Integrity**: Users can only manage *their own* tasks.
- **Pagination**: Efficient data loading for task lists.
- **API Documentation**: Integrated **Swagger UI** for testing endpoints directly.
- **Hot Reload**: configured with **Air** for smooth development.

## 🛠️ Tech Stack

- **Language**: Golang
- **Framework**: Gin Gonic
- **Database**: PostgreSQL
- **ORM**: GORM
- **Docs**: Swaggo (Gin-Swagger)

## 📂 Project Structure

```bash
task-tracker/
├── config/         # Database connection setup
├── controllers/    # Request handlers (logic)
├── docs/           # Swagger generated files
├── helpers/        # Utility functions (Response, Error, JWT)
├── inputs/         # Structs for request validation
├── middlewares/    # Auth middleware (JWT Guard)
├── models/         # Database structs
├── main.go         # Entry point & Routes
└── .env            # Environment variables