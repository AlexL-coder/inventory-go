# README

## Inventory Management API

This project is a Go-based web API that provides inventory management features, including user authentication via JWT tokens, database operations using SQLite, and PDF report generation.

---

### Features
- **User Authentication**: Secure login using JWT tokens.
- **Generate Reports**: Fetch inventory data and generate PDF reports.
- **Database Integration**: SQLite database for storing inventory items.
- **Concurrent Processing**: Efficient data handling using goroutines and channels.

---

### Prerequisites

Ensure you have the following installed:

- **Go**: Version 1.213 or higher
- **Docker**: (Optional) For containerized deployment
- **SQLite**: Pre-installed or accessible within the Docker container

---

### Installation

#### 1. Clone the Repository
```bash
git clone https://github.com/your-username/inventory-management-api.git
cd inventory-management-api
```

#### 2. Install Dependencies
```bash
go mod download
```

#### 3. Run the Application
```bash
go run main.go
```

The API will start on `http://localhost:8080`.

---

### API Endpoints

#### 1. **Login**
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "user",
    "password": "***"
  }
  ```
- **Response**:
  ```json
  {
    "token": "<JWT_TOKEN>"
  }
  ```

#### 2. **Generate Report**
- **URL**: `/generate-report`
- **Method**: `GET`
- **Headers**:
    - `Authorization: Bearer <JWT_TOKEN>`
- **Response**: PDF file containing the inventory report.

---

### Docker Deployment

#### 1. Build the Docker Image
```bash
docker build -t inventory-management-api .
```

#### 2. Run the Docker Container
```bash
docker run -p 8080:8080 --name inventory-management-api inventory-management-api
```

---

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

### Contributing

Feel free to fork this repository and submit pull requests. All contributions are welcome!

---

### Contact

For questions or support, contact [your-email@example.com](mailto:your-email@example.com).

---

# LICENSE

MIT License

Copyright (c) 2024 Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
