# Airbnb API

A small yet production-minded backend that models core Airbnb flows: user & property-owner authentication, property management, and bookings. Includes a clear path to scale from a single instance to tens of millions of users.

## Overview

The **Airbnb API** is a mini backend application that enables users to book properties or workspaces across the platform.  

- **User Section**:  
  Users can create an account and perform booking operations.  

- **Property Owner Section**:  
  Property owners can create an account, list/manage their properties, and handle bookings related to their properties.  

---

## Tech Stack

- **Language/Framework**: Go (Gin), GORM  
- **Database**: PostgreSQL  
- **Authentication**: JWT (HS256)  
- **Documentation**: Swagger (swaggo)  
- **Containerization**: Docker & Docker Compose  

---

## Project Setup

Since this is a demo system (and to avoid cloud costs), you can run everything locally using **Docker Compose**.  

### Steps
1. **Clone the repository**
   ```bash
   git clone git@github.com:Philip-21/AirBnb-Test-Api.git
   cd AirBnb-Test-Api

2. **Run the App**
    ```
    docker compose up --build

3. **Interact with the App**
   ```
    http://localhost:8080/swagger/index.html


## Architecture

This project follows a **monolithic MVC architecture**:

- **Models**: Define the data schema and ORM mappings (Users, Property Owners, Properties, Bookings).  
- **Controllers/Handlers**: Contain the business logic for handling API requests.  
- **Repositories**: Abstract database access using GORM.  
- **Middleware**: Handle JWT authentication and authorization.  
- **Routes**: Define endpoints grouped by user, property, and booking contexts.  

The application is packaged as a **Dockerized monolith** that connects to a PostgreSQL database.

---

## Current Capacity

- Designed as a lightweight demo backend for **20–50 concurrent users**.  
- Runs as a **single instance** with PostgreSQL as the persistence layer.  
- All operations are synchronous; no queues or background workers are used.  

---

## Challenges, Issues & Insights (Scaling Beyond Thousands of Users)

### 1. Monolith Limitations
- A single codebase is simple for development, but as traffic grows, deploy cycles, bug isolation, and scaling specific modules (e.g., bookings vs. auth) become difficult.
- Horizontal scaling requires replicating the entire app, even if only one part is under heavy load.

### 2. Database Bottlenecks
- PostgreSQL on a single instance will eventually hit **I/O and query throughput limits**.  
- Heavy joins (e.g., user + property + booking lookups) increase latency as data grows.  
- Write-heavy operations like bookings can cause row locking and contention.

### 3. Authentication & Security
- JWT is stateless and fast, but token revocation is tricky without extra infrastructure (e.g., Redis for blacklisting).  
- At scale, stronger monitoring and rotation of secret keys would be required.

### 4. Scaling to Millions of Users
To handle **100k+ concurrent users** or **millions of accounts**:  
- **Database Scaling**: Introduce **read replicas**, sharding, and caching (e.g., Redis).  
- **Service Decomposition**: Split monolith into microservices — Auth, Property, Booking — each independently scalable.  
- **Load Balancing**: Use Nginx/HAProxy + horizontal scaling with Kubernetes or ECS.  
- **Message Queues**: For async workflows (e.g., booking confirmations, email notifications).  
- **CDN & Caching**: To reduce DB hits and serve static metadata faster.  

### 5. Performance Without Increasing Cost
- Optimize queries with proper **indexes**.  
- Use **connection pooling**.  
- Cache frequently accessed data (property lists, booking statuses).  
- Offload heavy analytics to a separate data pipeline.  

### 6. Measuring TPS (Transactions Per Second)
- Use tools like **k6**, **Locust**, or **JMeter** for load and stress testing.  
- Baseline TPS depends on infrastructure (single instance vs. multi-node).  
- With proper horizontal scaling + DB optimization, system could reach **thousands of TPS**.  

---

## Insights
- **Monoliths are great to start small**: faster to build and test.  
- **For production-scale Airbnb-like systems**: inevitable move toward microservices, distributed databases, and asynchronous systems.  
- **Early investment in clean boundaries (repos, handlers, middleware)** in this project makes future refactoring easier.  
