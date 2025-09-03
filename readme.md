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
   git clone <repo-url>
   cd <repo-name>

2. *Run the App **
    `docker compose up --build`
    