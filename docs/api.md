# API Documentation

For more information visit the [documentation page](https://documenter.getpostman.com/view/4425953/2s9Y5bQgpV)

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