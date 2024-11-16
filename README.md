
# Retail Pulse Image Processing Service

This project implements a backend service that processes thousands of images collected from stores, fulfilling the requirements outlined in the assignment. The service is built using Go with Dockerized support for seamless deployment and scalability.

## Project Description

The service provides two primary APIs:
1. `/api/submit`: Accepts job submissions containing store IDs, image URLs, and visit times. It validates inputs, creates jobs, and processes images asynchronously.
2. `/api/status`: Retrieves the status of a job (ongoing, completed, or failed) and provides detailed error messages if any image processing fails.

The system calculates the perimeter of images (`2 * [Height + Width]`) and simulates GPU processing using random sleep times (0.1 to 0.4 seconds). Results are stored at the image level.

## Assumptions

1. The `stores.csv` file provided contains the store data (`store_id`, `store_name`, `area_code`), and we load this data into the database (`stores` table) for validation.
2. Jobs are associated with valid store IDs, enforced via foreign key constraints.
3. Image dimensions are fetched dynamically during processing.
4. Both Dockerized and local (non-Docker) setups are supported.

## Installation and Setup Instructions

### Prerequisites
- **Go** (version 1.20 or higher)
- **Docker** and **Docker Compose**
- **PostgreSQL** (if running without Docker)

### Steps to Run Locally
1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd retail_pulse_project
   ```

2. Set up environment variables:
   Create a `.env` file in the root directory with the following content:
   ```
   DATABASE_URL=postgresql://<username>:<password>@localhost:5432/retail_pulse
   ```

3. Load the `stores.csv` file:
   Ensure that the `stores.csv` file is present in the root directory. The service will automatically load the data into the `stores` table when starting.

4. Run the PostgreSQL database (if not using Docker):
   ```bash
   docker run --name retail_pulse_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=retail_pulse -p 5432:5432 -d postgres:14
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

### Steps to Run with Docker
1. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. The application will be available at `http://localhost:8080`.

3. The PostgreSQL database will be available at `localhost:5432` with the following credentials:
   ```
   User: postgres
   Password: password
   Database: retail_pulse
   ```

### Testing the APIs
Use tools like **Postman** or **curl** to test the APIs.

#### 1. Submit Job API
**Endpoint:** `POST /api/submit`  
**Request Payload:**
```json
{
  "count": 2,
  "visits": [
    {
      "store_id": "RP00001",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/2.jpg",
        "https://www.gstatic.com/webp/gallery/3.jpg"
      ],
      "visit_time": "2024-11-15T10:00:00Z"
    },
    {
      "store_id": "RP00002",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/4.jpg"
      ],
      "visit_time": "2024-11-15T11:00:00Z"
    }
  ]
}
```
**Response:**
```json
{
  "job_id": "1234567890123456789"
}
```

#### 2. Get Job Status API
**Endpoint:** `GET /api/status?jobid=<job_id>`  
**Response (Completed):**
```json
{
  "status": "completed",
  "job_id": "1234567890123456789"
}
```
**Response (Failed):**
```json
{
  "status": "failed",
  "job_id": "1234567890123456789",
  "error": [
    {
      "store_id": "RP00001",
      "error": "image: unknown format"
    }
  ]
}
```

## Work Environment

- **Operating System:** Windows 10
- **IDE:** Visual Studio Code
- **Programming Language:** Go (v1.20)
- **Database:** PostgreSQL (v14)
- **Containerization:** Docker (v24.0.5)
- **Other Tools:** Postman for API testing

## Future Improvements

If given more time, the following improvements can be made:
1. **Scalability:** Implement worker queues (e.g., RabbitMQ) for handling image processing jobs more efficiently.
2. **Monitoring:** Add logging and monitoring tools like Prometheus and Grafana.
3. **Authentication:** Secure the APIs with user authentication and authorization.
4. **Enhanced Validation:** Add more robust validation for input payloads.
5. **Error Handling:** Improve error reporting by storing detailed logs in the database or a logging service.
6. **Testing:** Add unit and integration tests for the service.

## Conclusion

This project implements all the requirements outlined in the assignment, providing a robust service for processing images collected from stores. Both Dockerized and local setups are supported, ensuring flexibility in deployment and testing.
