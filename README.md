# CoverLetterAI - API

## What is CoverLetterAI

Introducing **CoverLetterAI**, your new best friend in the job application process. Utilizing the advanced capabilities of ChatGPT, our web app crafts compelling and personalized cover letters in minutes. Simply answer a few questions or give us some key points, and we'll turn them into a professional cover letter that stands out. Say goodbye to staring at a blank screen wondering how to start; let **CoverLetterAI** do the heavy lifting so you can focus on landing your dream job.

## Installation

1. Make sure you have Go installed. If not, you can download it from [here](https://golang.org/dl/).

2. Clone this repository

3. Create a `.env` file by running `cp .env.example .env` with your environment variables

4. Start the docker containers (`docker compose up -d`)

5. The API routes will be available on `http://localhost:8080`

## Development

1. When making code changes, you can use `go run ./cmd/api` without using the application docker container

2. Make sure you are running the `mongodb` container (`docker compose up -d mongodb`)

## Testing

**Note** To generate/update mocks, run `go generate ./...`

1. Make you are running the `mongodb-test` container (`docker compose up -d mongodb-test`)
2. Run `go test ./...` to run all the existing tests

## API Endpoints

`POST`: `/cover-letter`

Creates a cover letter.

#### Request:
```json
{
  "email": "your@email.com",
    "job_posting": {
    "company_name": "Acme Inc.",
    "job_role": "CEO",
    "job_details": "Looking for a very experienced CEO",
    "skills": "sales, accounting, management"
  },
}
```
#### Response:
```json
{
  "data": "This is a cover letter for a CEO position"
}
```


`POST`: `/career-profile`

Creates a career profile (Recommended for more accurate cover letters)

#### Request:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "headline": "CEO",
  "experience_years": 10,
  "skills": ["management", "accounting", "sales"],
  "summary": "I am the best CEO",
  "contact_info": {
      "email": "your@email.com",
      "address": "Toronto, ON",
      "phone": "+1 555-555-5555",
      "website": "mywebsite.com"
  }
}
```
#### Response:
```json
{
  "data": {
      "id": "00000000-0000-0000-0000-000000000000",
      "first_name": "John",
      "last_name": "Doe",
      "headline": "CEO",
      "experience_years": 10,
      "summary": "I am the best CEO",
      "skills": [
          "management",
          "accounting",
          "sales"
      ],
      "contact_info": {
          "email": "your@email.com",
          "address": "Toronto, ON",
          "phone": "+1 555-555-5555",
          "website": "mywebsite.com"
      }
  },
  "message": "career profile has been created"
}
```

`GET`: `/career-profile/:email`

Returns an existing career profile for a given email address

#### Request: `/career-profile/your@email.com`

#### Response:
```json
{
  "data": {
      "id": "00000000-0000-0000-0000-000000000000",
      "first_name": "John",
      "last_name": "Doe",
      "headline": "CEO",
      "experience_years": 10,
      "summary": "I am the best CEO",
      "skills": [
          "management",
          "accounting",
          "sales"
      ],
      "contact_info": {
          "email": "your@email.com",
          "address": "Toronto, ON",
          "phone": "+1 555-555-5555",
          "website": "mywebsite.com"
      }
  }
}
```
