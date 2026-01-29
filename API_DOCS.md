# POS Backend API Documentation

## Overview
This document describes all RESTful endpoints for the POS backend, including request/response formats, authentication, and business logic. Use this as a reference for Postman or any API client.

---

## Authentication
- All endpoints (except `/business/register` and `/business/login`) require a valid JWT access token in the `Authorization: Bearer <token>` header.

---

## Endpoints

### 1. Register Business
- **POST** `/api/business/register`
- **Body:**
```json
{
  "business_name": "string",
  "owner_full_name": "string",
  "email": "string",
  "phone_number": "string",
  "password": "string",
  "store_address": "string",
  "business_category": "string",
  "currency": "string",
  "store_icon": "string (optional)"
}
```
- **Response:**
```json
{
  "business_id": "string"
}
```
- **Logic:** Registers a new business, hashes password, creates main branch.

---

### 2. Login
- **POST** `/api/business/login`
- **Body:**
```json
{
  "email": "string",
  "password": "string"
}
```
- **Response:**
```json
{
  "access_token": "string",
  "refresh_token": "string",
  "business_id": "string",
  "role": "owner"
}
```
- **Logic:** Authenticates business owner, returns JWT tokens.

---

### 3. Create Branch
- **POST** `/api/branch/create?business_id=...`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
```json
{
  "branch_name": "string",
  "branch_address": "string",
  "is_main_branch": false
}
```
- **Response:**
```json
{
  "branch_id": "string",
  "business_id": "string",
  "branch_name": "string",
  "branch_address": "string",
  "is_main_branch": false,
  "created_at": 0,
  "updated_at": 0
}
```
- **Logic:** Creates a new branch for the business.

---

### 4. Delete Branch
- **DELETE** `/api/branch/delete?branch_id=...&business_id=...`
- **Headers:** `Authorization: Bearer <token>`
- **Response:** `204 No Content`
- **Logic:** Soft deletes a branch (main branch cannot be deleted).

---

### 5. Create Staff
- **POST** `/api/staff/create`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
```json
{
  "full_name": "string",
  "email": "string",
  "phone_number": "string",
  "role": "owner|manager|cashier|inventory_staff",
  "branch_id": "string",
  "status": "active|inactive",
  "photo_url": "string (optional)"
}
```
- **Response:**
```json
{
  "id": "string",
  "full_name": "string",
  "email": "string",
  "phone_number": "string",
  "role": "string",
  "branch_id": "string",
  "status": "string",
  "photo_url": "string (optional)",
  "created_at": 0,
  "updated_at": 0
}
```
- **Logic:** Adds a staff member to a branch.

---

### 6. Add Product
- **POST** `/api/product/add`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
```json
{
  "product_name": "string",
  "product_category": "string",
  "selling_price": 0,
  "cost_price": 0,
  "quantity_in_stock": 0,
  "low_stock_threshold": 0,
  "barcode_value": "string (optional)",
  "nafdac_reg_number": "string (optional)",
  "expiry_date": 0,
  "product_image_url": "string (optional)",
  "branch_id": "string"
}
```
- **Response:**
```json
{
  "id": "string",
  "product_name": "string",
  "nafdac_reg_number": "string (optional)",
  "selling_price": 0,
  "quantity_left": 0,
  "product_image_url": "string (optional)"
}
```
- **Logic:** Adds a product to a branch's inventory.

---

### 7. Sync Data (Offline-First)
- **POST** `/api/sync`
- **Headers:** `Authorization: Bearer <token>`
- **Body:**
```json
{
  "device_id": "string",
  "last_sync_timestamp": 0,
  "request_signature": "string",
  "data": [ ... ]
}
```
- **Response:**
```json
{
  "status": "success",
  "conflicts": [ ... ]
}
```
- **Logic:** Syncs offline data, supports batch uploads, idempotency, and conflict resolution.

---

### 8. Refresh Token
- **POST** `/api/auth/refresh`
- **Body:**
```json
{
  "refresh_token": "string"
}
```
- **Response:**
```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```
- **Logic:** Issues a new access token using a valid refresh token.

---

## Error Handling
- All errors are returned as JSON with an appropriate HTTP status code and message.

---

## Security
- All input is validated and sanitized to prevent XSS and injection attacks.
- JWT is required for all protected endpoints.
- Main branch cannot be deleted.
- All write endpoints are idempotent.

---

## Notes
- Replace `localhost:8080` with your server address in Postman.
- For full business logic, see the codebase.
