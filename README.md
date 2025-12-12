# 🐾 VetClinic REST API

Backend service for managing veterinary clinic operations — including users, pets, appointments, medical records, and treatments. <br>
Built with **Go**, **Gin**, **PostgreSQL**, and **JWT Authentication** along with **Railway** Deployment.

---

## 🌐 Live API Base URL (Railway)

https://vetclinic-rest-api-production.up.railway.app

Open the link to view the full API documentation, including all endpoints and example responses.

---

## 🌐 API Documentation (Postman)

https://documenter.getpostman.com/view/15890346/2sB3dSQot9

Open the domain to check the api documentation of this project to check the path with example result.

## ✅ Overview

VetClinic API is designed as an **internal system** for veterinary clinics to manage daily operations efficiently.
This backend provides secure and structured endpoints for:

-   Managing clinic staff and doctors
-   Registering and tracking pets
-   Scheduling and updating appointments
-   Recording medical diagnoses
-   Logging treatments and calculating costs

The system uses **role-based access control (RBAC)** to ensure Admin, Staff, and Doctor roles only access the features intended for them.

---

## ✅ Purpose

This project was created to:

-   Learn and implement **REST API architecture** using Go + Gin
-   Build a **realistic clinic management backend**
-   Practice **JWT authentication**, **middleware**, and **RBAC**
-   Understand **relational database design** with PostgreSQL
-   Simulate real clinic workflows: appointments → medical records → treatments
-   Provide a clean, maintainable backend for potential future frontend/mobile apps

---

## ✅ Features

### 🔐 Authentication & Authorization

-   JWT-based login
-   Password hashing with bcrypt
-   Role-Based Access Control (Admin, Staff, Doctor)
-   Protected routes via middleware

### 👥 User Management

-   Register new users
-   Login
-   Fetch and update users detail (all roles)
-   Fetch users data by role (all roles)
-   Update users password (all roles)
-   Update role (Admin only)
-   Soft delete user (Admin only)

### 🐶 Pet Management (CRUD)

-   Create, read, update, soft delete pets (Staff/Admin)
-   Fetch pets data by owner name and phone (Staff/Admin)

### 📅 Appointment Management (CRUD)

-   Create appointments (Staff/Admin)
-   Update appointment details (Staff/Admin)
-   Update appointment status (pending/cancelled/completed) (all roles)
-   Soft delete appointments (Staff/Admin)
-   Filter by pet, doctor, or date (all roles)
-   Full appointment detail (appointment + pet + medical record + treatments) (all roles)

### 🩺 Medical Records

-   One medical record per appointment
-   Create, update, soft delete (Doctor/Admin)
-   Fetch medical records by appointment id (all roles)

### 💊 Treatments

-   Multiple treatments per medical record
-   Create, update, soft delete (Doctor/Admin)
-   Fetch treatments by medical records id (all roles)

### 🛡️ Middleware

-   JWT validation
-   Role-based route protection
-   Logging & error handling

---

## ✅ Tech Stack

-   **Go (Golang)**
-   **Gin Web Framework**
-   **PostgreSQL**
-   **JWT Authentication**
-   **bcrypt**
-   **UUID for entity IDs**
-   **Railway Deployment**

---

## ✅ Running Locally

### ▶️ Start the server

```bash
go run main.go

Local server runs on:
http://localhost:8000
```

## ✅ API Path List

Use Postman or other API tools to test the API.<br>
Below is the complete list of available routes grouped by feature.

🔐 AUTH & USERS API
Base: `/api/users`

-   POST `/api/users/register` — Register new user
-   POST `/api/users/login` — Login and receive JWT token
    User Profile & Management
-   GET `/api/users/:id/profile` — Get user profile (Staff, Doctor, Admin)
-   GET `/api/users/role/:role` — Get users by role (Staff, Doctor, Admin)
-   PUT `/api/users/:id/update` — Update user info (Staff, Doctor, Admin)
-   PUT `/api/users/:id/change-password` — Change password (Staff, Doctor, Admin)
-   PUT `/api/users/:id/role` — Update user role (Admin only)
-   PUT `/api/users/:id/active-status` — Soft delete user (Admin only)

🐾 PETS API
Base: `/api/pets`

-   GET `/api/pets/:id/profile` — Get pet profile (Staff, Doctor, Admin)
-   GET `/api/pets/by-owner/:owner_name/:owner_phone` — Get pets by owner name & phone (Staff, Doctor, Admin)
-   POST `/api/pets` — Create new pet (Staff, Admin)
-   PUT `/api/pets/:id` — Update pet (partial update supported) (Staff, Admin)
-   PUT `/api/pets/:id/active-status` — Soft delete pet (Staff, Admin)

📅 APPOINTMENTS API
Base: `/api/appointments`

-   POST `/api/appointments` — Create appointment (Staff, Admin)
-   GET `/api/appointments/:id` — Get appointment by ID (Staff, Doctor, Admin)
-   PUT `/api/appointments/:id` — Update appointment details (Staff, Admin)
-   PUT `/api/appointments/:id/status` — Update appointment status (Staff, Doctor, Admin)
-   PUT `/api/appointments/:id/active-status` — Soft delete appointment (Staff, Admin)
-   GET `/api/appointments/pet/:pet_id` — Get appointments by pet (Staff, Doctor, Admin)
-   GET `/api/appointments/doctor/:doctor_id` — Get appointments by doctor (Staff, Doctor, Admin)
-   GET `/api/appointments/date/:date` — Get appointments by date (Staff, Doctor, Admin)
-   GET `/api/appointments/:id/full` — Get full appointment detail (pet + medical record + treatments)

🩺 MEDICAL RECORDS API
Base: `/api/medical-records`

-   POST `/api/medical-records` — Create medical record (Doctor, Admin)
-   GET `/api/medical-records/appointment/:appointment_id` — Get medical record by appointment ID (Doctor, Staff, Admin)
-   PUT `/api/medical-records/:id` — Update medical record (partial update supported) (Doctor, Admin)
-   PUT `/api/medical-records/:id/active-status` — Soft delete medical record (Admin)

💊 TREATMENTS API
Base: `/api/treatments`

-   POST `/api/treatments` — Create treatment (Doctor, Admin)
-   GET `/api/treatments/medicalrecord/:medicalrecord_id` — Get treatments by medical record (Doctor, Staff, Admin)
-   PUT `/api/treatments/:id` — Update treatment (partial update supported) (Doctor, Admin)
-   PUT `/api/treatments/:id/active-status` — Soft delete treatment (Doctor, Admin)

## 🚀 Future Improvements

Planned enhancements for future versions:<br>

✅ 1. Owner Role (End-User Access)

-   Add Owner role in Users table
-   Allow pet owners to view their pets, appointments, and medical history
-   Allow owners to request appointments directly

✅ 2. Transaction & Payment System

-   Add Transactions table
-   Track total treatment cost per appointment
-   Integrate payment gateway
-   Generate invoices for owners

✅ 3. Treatment Cost Breakdown

-   Replace single cost field with:
-   Treatment item
-   Quantity
-   Unit price
-   Total price

✅ 4. Appointment Request System

-   Owners can request appointment slots
-   Staff/Admin approve or reject requests
