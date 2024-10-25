# All-in-One Transport Application

## Overview

This is a backend system for a transport application focused on Kochi, Kerala. The application allows users to view routes, book transport services, and manage payments. Admins can manage routes and stops. The backend is built using **Go**, **GORM**, **Gin** and **PostgreSQL**, following **Clean Architecture** principles with JWT-based authentication.

## Features

- JWT authentication for users and admins
- Dynamic fare calculation based on distance or stops (Haversine Formula/OSM)
- Route management for admins (Create, Update Routes And Stops)
- User bookings, invoices, and payments handling
- Monthly and seasonal subscription plans
- Discount coupons for fare reductions
- Google Single Sign On
- OTP verification at the time of Sign up
- OSM (Open Street Map)
- Wallet (Refund & Payment)
- Razorpay Integration (Test Mode)
- Invoice generation (gofpdf )
- Admin dashboard

## Prerequisites

Before running the application, make sure you have the following installed:
- **Go**: [Installation Guide](https://golang.org/doc/install)
- **Gorm**: [Installation Guide](gorm.io/gorm)
- **PostgreSQL**: [Installation Guide](https://www.postgresql.org/download/)
- **Git**: [Installation Guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- **OSRM** **OSM**

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Prototype-1/xtrace.git
cd xtrace
