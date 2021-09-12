# Home24 Assignment

### Task description
Create a web application which takes a website URL as an input and provides general information
about the contents of the page:
- HTML Version
- Page Title
- Headings count by level
- Amount of internal and external links
- Amount of inaccessible links
- If a page contains a login form

## Instructions to RUN

- Install docker and start docker daemon
- Run the following command to start the server

``` shell 
make docker-run 
```

- Wait for the server to get started
- Run the following curl for making a request to the server

``` shell
curl --location --request POST 'localhost:8000/analyse' \
--header 'Content-Type: application/json' \
--data-raw '{"page_url": "https://www.facebook.com"}' 
```

- Replace the "page_url" value for the testing.

### Sample Request

```json
{
  "page_url": "https://www.facebook.com"
}
```

### Sample Response

```json
{
  "has_login_form": true,
  "header_count": {
    "h1": 1,
    "h2": 2,
    "h3": 3,
    "h4": 4,
    "h5": 5,
    "h6": 6
  },
  "html_version": "HTML 3",
  "links": {
    "external": {
      "count": 4,
      "urls": [
        "https://www.facebook.com",
        "https://www.google.com",
        "https://www.amazon.com",
        "https://www.example.com/inaccessible"
      ]
    },
    "internal": {
      "count": 2,
      "urls": [
        "/internal1",
        "/internal2"
      ]
    },
    "inaccessible": {
      "count": 1,
      "urls": [
        "https://www.example.com/inaccessible"
      ]
    }
  },
  "title": "Parvez Hassan Test Page"
}
```

### Notes for reviewer

- Have written integration and UT for each and every file but couldn't complete all test cases in UT for parser file due
  to time crunch.
- Currently, I have made the following assumption in my implementation as those requirements were not clear in the
  problem statement
  - Since the login form implementation could be different for different pages, currently I'm checking it based on the
    existence of password type field and submit type field.
  - Considering links as accessible only if they are returning 200 as response, we can also consider 3xx series as a
    valid response if req.
