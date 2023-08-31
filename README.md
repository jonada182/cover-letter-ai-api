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

## API Documentation

Check out the [API documentation](https://github.com/jonada182/cover-letter-ai-api/blob/main/docs/api.md#api-documentation)
