![Recording2024-12-27at02 19 14-ezgif com-optimize](https://github.com/user-attachments/assets/379acace-061a-420f-9c5a-c6121d17c860)

# NAITA Report Generator

A simple tool for generating weekly training reports for NAITA . This application allows trainees to input their daily progress and automatically generates a formatted report document.

## Features

- Simple, user-friendly web interface
- Daily progress tracking for 5 working days
- Key learnings summary
- Automatic report generation in DOCX format

## Prerequisites

- Go 1.21 or higher
- Google Cloud API key for Gemini

## Installation

1. Clone the repository:
```bash
git clone https://github.com/lakpahana/write-my-naita
cd write-my-naita
```

2. Set up environment variables:
```bash
cp .env.local .env
# API_KEY = "your_api_key"
```

3. Install Go dependencies:
```bash
go mod download
```

## Running the Application

1. Start the server:
```bash
cd cmd/api	
go run main.go
```
The server will start on http://localhost:8080

2. Open the frontend:

Then visit http://localhost:8080

If you love CLI like me, you can go to 
```bash
cd cmd/cli
go run main.go
```

## Usage

1. Fill in the project details (name and technologies used)
2. Enter your daily progress for each of the 5 working days
3. Add your key learnings for the week
4. Click "Generate Report"
5. Once generated, click "Download Report" to get your DOCX file
