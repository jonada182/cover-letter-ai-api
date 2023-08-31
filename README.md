# CoverLetterAI - API

## What is CoverLetterAI

Introducing **CoverLetterAI**, your new best friend in the job application process. Utilizing the advanced capabilities of ChatGPT, our web app crafts compelling and personalized cover letters in minutes. Simply answer a few questions or give us some key points, and we'll turn them into a professional cover letter that stands out. Say goodbye to staring at a blank screen wondering how to start; let **CoverLetterAI** do the heavy lifting so you can focus on landing your dream job.

## Installation

1. Make sure you have Go installed. If not, you can download it from [here](https://golang.org/dl/).

2. Clone this repository

3. Create a `.env` file by running `cp .env.example .env` with your environment variables

4. Start the application services (`task run`)

5. The API routes will be available on `http://localhost:8080`

## Testing

**Note** To generate/update mocks, run `task mock`

Run `task test` to run all the existing tests

## Taskfile Commands

Here are the available commands from the [Taskfile](https://taskfile.dev/):

* `run`: Start application services
* `stop`: Stop application services
* `dev`: Run dev environment
* `test`: Run application tests
* `mock`: Generate mocks using gomock

## API Documentation

Check out the [API documentation](https://github.com/jonada182/cover-letter-ai-api/blob/main/docs/api.md#api-documentation)
