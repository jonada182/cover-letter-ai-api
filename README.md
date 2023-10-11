# CoverLetterAI - API

## What is CoverLetterAI

Introducing **CoverLetterAI**, your comprehensive companion in the job application journey. Not only does our platform harness the advanced capabilities of ChatGPT to craft compelling and personalized cover letters in minutes, but it also offers a robust tracking system for all your job applications. Monitor your progress with detailed event logs, from submission to interview, and even job offers. No more juggling multiple tools or staring at a blank screen wondering how to begin. With **CoverLetterAI**, you get a holistic solution that lets you focus on what matters most: landing your dream job.

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
