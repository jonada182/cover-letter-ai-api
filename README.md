# Resu-mate API

This Golang project offers a set of API endpoints designed to streamline common job application tasks. For instance, it enables you to effortlessly create cover letters and perform other related actions.

## Installation

1. Make sure you have Go installed. If not, you can download it from [here](https://golang.org/dl/).

2. Clone this repository

3. Run the app `go run main.go`

## API Endpoints

`POST`: `/cover-letter`

Creates a cover letter.

#### Request:
```json
{
  "company_name": "Acme Inc.",
  "job_role": "Software Engineer",
  "job_details": "Developing cutting-edge software applications.",
  "skills": "Golang, Java, JavaScript"
}
```
#### Response:
```json
{
  "message": "Message from OpenAI"
}
```
