# EmailEngine REST API
When the Engine instance is started with `e.StartHTTP(...)`, it exposes a REST API that allows external applications to enqueue outbound emails for delivery.

This document describes the available data models and endpoints provided by the API.

## Table of Contents
- [EmailEngine REST API](#emailengine-rest-api)
	- [Table of Contents](#table-of-contents)
- [Objects](#objects)
	- [Email](#email)
	- [Address](#address)
	- [Attachment](#attachment)
- [Endpoints](#endpoints)
	- [Queue Outbound Emails](#queue-outbound-emails)
		- [Request Body](#request-body)
		- [Responses](#responses)

# Objects

## Email
Represents a complete email message ready to be queued and sent.

| Field       | Type                        | Description                                                           |
| ----------- | --------------------------- | --------------------------------------------------------------------- |
| to          | [Address[]](#address)       | One or more recipients of the email. Must include at least one entry. |
| from        | [Address](#address)         | The sender’s name and email address.                                  |
| subject     | string                      | The subject line of the email (max 255 characters).                   |
| content     | string                      | The main message body. Can be plain text or HTML depending on `html`. |
| html        | boolean                     | Set to `true` if `content` is HTML; `false` for plain text.           |
| attachments | [Attachment[]](#attachment) | Optional. One or more file attachments or inline images.              |


## Address
Represents either a sender or recipient email address.

| Field   | Type   | Description                                                                                                   |
| ------- | ------ | ------------------------------------------------------------------------------------------------------------- |
| name    | string | Display name of the sender or recipient (1–128 characters)                                                    |
| address | string | The email address in [RFC 5322](https://datatracker.ietf.org/doc/html/rfc5322#section-3.4.1) format (max 128) |


## Attachment
Defines a file or inline resource attached to an email.

| Field        | Type    | Description                                                              |
| ------------ | ------- | ------------------------------------------------------------------------ |
| filename     | string  | The name of the attached file (max 255 characters)                       |
| content_type | string  | The MIME type of the attachment (e.g. `image/png`, `text/plain`)         |
| data         | string  | Binary content of the file encoded in Base64                             |
| inline       | boolean | Set to `true` to embed the attachment inline (e.g. for logos and images) |

> **TIP:** Inline attachments like images can be referenced in HTML emails using a Content-ID URL (e.g. `cid:logo.png`)


<br>


# Endpoints

All REST API requests must include:
- A JSON-encoded payload (`Content-Type: application/json`)
- A maximum payload size of **10 MB** (or a custom limit defined by `IncomingMaxBytes`)

## Queue Outbound Emails
`POST /queue`

Adds one or more emails to the outbound queue for later processing and delivery.

### Request Body
An array of [Email](#email) objects:
```json
[
	{
		"to": [
			{
				"name": "bakonpancakz",
				"address": "bakonpancakz@gmail.com"
			}
		],
		"from": {
			"name": "emailengine",
			"address": "emailengine@example.org"
		},
		"html": true,
		"subject": "Setup Complete!",
		"content": "<h1>If you're reading this, your email system is working correctly!</h1><p>Take this art as your reward!</p><img src='cid:teto.png' alt='Kasane Teto'/>",
		"attachments": [
			{
				"content_type": "image/png",
				"filename": "teto.png",
				"data": "iVBORw0KGgoAAAANSUhEUgAAAMgAAADgBAMAAAC0iTT2AAAAKlBMVEUAAAAAAADVACCIABX87dHDw8Pw8PBPT0/v26/LX19/f3/tHCT/rsmGKysrGzxDAAAAAXRSTlMAQObYZgAAAadJREFUeNrt2kFxwzAQheFSCAVTMAVTCAVTKIVSKIVQCIVSKJfuTvWmbxTZTa6r/102kqP9fNlx3cnbq7lEXr0GAjITkk3ukay+p7pHsoKAzIqMADVWBQGZFbl0WSK+FgACUh25DKL95YkIGvUBAamAHDXSQ+qzJfey9mt9PusDAlIRWSNqsjwRfU/nlD0CAlIFWQ9yj+jg8Ea8tuwWEJAqSGQI6aA37ZEjSBEAAlIF6Rv0wEckq6/7vf5GQECqIasyAL4jfcMtMtrfLSAglZAjSEOXDbKpGmbNtfZH5wSAgFRDfADvEQ1dft5aBPj+6I9BEJCqiA+S1l8t2yC65mf2FhCQioggHyQd3k6im/Eb1XkQkIqIvviwd5Kj8yAglZH/BlUvP2fNQEBmRLLhLXKNrJGsuc59EBCQX+Da8h7xmsnrICAzI3qp0SAK8IHUSw8IyMyIHlQCtPYKAjIrokH0Iby1LBEfUA0kCMiMyBbxh5TizfP6FgEBmRERlBHkEeD/hAYBmRFRBPXx5iAgIH8/BvA1CAjI4zBm9TUICIghFl+DgEyC/AAl2AjG7TSPSgAAAABJRU5ErkJggg==",
				"inline": true
			}
		]
	}
]
```

### Responses
| Code                               | Meaning                                                 |
| :--------------------------------- | :------------------------------------------------------ |
| **`201 Created`**                  | Emails were successfully queued.                        |
| **`400 Bad Request`**              | One or more emails failed validation and were rejected. |
| **`401 Unauthorized`**             | The `AuthHandler` rejected the request.                 |
| **`413 Request Entity Too Large`** | Payload exceeds the maximum allowed size.               |
| **`415 Unsupported Media Type`**   | The `Content-Type` header is not `application/json`.    |
| **`422 Unprocessable Entity`**     | The payload is invalid or malformed JSON.               |
| **`507 Insufficient Storage`**     | The queue is full and cannot accept additional emails.  |
