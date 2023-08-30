# CoverLetterAI - API

## What is CoverLetterAI

Introducing **CoverLetterAI**, your new best friend in the job application process. Utilizing the advanced capabilities of ChatGPT, our web app crafts compelling and personalized cover letters in minutes. Simply answer a few questions or give us some key points, and we'll turn them into a professional cover letter that stands out. Say goodbye to staring at a blank screen wondering how to start; let **CoverLetterAI** do the heavy lifting so you can focus on landing your dream job.

## Installation

1. Make sure you have Go installed. If not, you can download it from [here](https://golang.org/dl/).

2. Clone this repository

3. Create a `.env` file by running `cp .env.example .env` with your environment variables

4. Run the app `go run main.go`

## Testing

1. To generate/update mocks, run `go generate ./...`

2. To run tests: `go test`

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
