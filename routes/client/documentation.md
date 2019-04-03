## Backpulse API

This is the official Backpulse API.  
API Endpoint:
```
https://api.backpulse.io/:sitename
```
#### Successful request
Example of a successful request.
```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": {...}
}
```

#### Unsuccessful request
Exemple of 404 error.
```json
{
    "status": "error",
    "code": 404,
    "message": "not_found",
    "payload": null
}
```

### Models
Examples of the API models.

#### Gallery
```json
{
    "short_id": "ef855uQmR",
    "site_name": "mysite",
    "title": string, // English title
    "titles": []Translation,
    "descriptions": []Translation,
    "photos": []Photo,
    "created_at": Date,
    "updated_at": Date
}
```

#### Project
```json
{
    "short_id": "NzFVjUwig",
    "site_name": "aureleoules",
    "title": string, // English title
    "titles": []Translation,
    "descriptions": []Translation,
    "url": string,
    "created_at": Date,
    "updated_at": Date
}
```

#### Translation
```json
{
    "language_name": "English",
    "language_code": "en",
    "content": string
}
```
#### Photo
```json
{
    "url": "https://res.cloudinary.com/zygtpradi/image/authenticated/x--3BfHaB4O--/v1547280893/le6bd7z918vc8qas1s.jpg",
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "index": 0,
    "created_at": Date
}
```

####


## Routes
List of all the API routes.

### Galleries

Lists all galleries.

```endpoint
GET /galleries
```

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/galleries
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": [
        {
            "short_id": "eS9u5RQmR",
            "site_name": "mysite",
            "title": string,
            "titles": []Translation,
            "descriptions": []Translation,
            "photos": []Photo,
            "created_at": Date,
            "updated_at": Date
        },
        ...
    ]
}
```

### Gallery

Get a specific gallery.

```endpoint
GET /gallery/:short_id
```
Route parameter | Description
--- | ---
`short_id` | id of an existing gallery

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/gallery/NzFVjUwig
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": {
        "short_id": "eS9u5RQmR",
        "site_name": "mysite",
        "title": string,
        "titles": []Translation,
        "descriptions": []Translation,
        "photos": []Photo,
        "created_at": Date,
        "updated_at": Date
    }
}
```

### Default gallery

Get the default gallery.

```endpoint
GET /galleries/home
```

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/galleries/home
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": {
        "short_id": "eS9u5RQmR",
        "site_name": "mysite",
        "title": string,
        "titles": []Translation,
        "descriptions": []Translation,
        "photos": []Photo,
        "created_at": Date,
        "updated_at": Date
    }
}
```

### Projects

Lists all projects.

```endpoint
GET /projects
```

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/projects
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": [
        {
            "short_id": "eS9u5RQmR",
            "site_name": "mysite",
            "title": string,
            "titles": []Translation,
            "descriptions": []Translation,
            "url": string,
            "created_at": Date,
            "updated_at": Date
        },
        ...
    ]
}
```

### Project

Get a specific project.

```endpoint
GET /project/:short_id
```
Route parameter | Description
--- | ---
`short_id` | id of an existing project

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/project/NzFVjUwig
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": {
        "short_id": "eS9u5RQmR",
        "site_name": "mysite",
        "title": string,
        "titles": []Translation,
        "descriptions": []Translation,
        "url": string,
        "created_at": Date,
        "updated_at": Date
    }
}
```

### Contact

Get contact informations.

```endpoint
GET /contact
```

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/contact
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": [
        {
            "site_name": "aureleoules",
            "name": "John Doe",
            "phone": "202-555-0199",
            "email": "contact@backpulse.io",
            "address": "355 Yukon Lane\nOak Park, MI 48237",
            "facebook_url": "https://facebook.com/...",
            "instagram_url": "https://instagram.com/...",
            "twitter_url": "https://twitter.com/...",
            "custom_fields": []CustomField
        }
    ]
}
```

### About
Get about informations.

```endpoint
GET /about
```

#### Example request

```curl
$ curl https://api.backpulse.io/mysite/about
```

#### Example response

```json
{
    "status": "success",
    "code": 200,
    "message": "success",
    "payload": [
        {
            "site_name": "aureleoules",
            "name": "John Doe",
            "descriptions": []Translation,
            "titles": []Translation
        }
    ]
}
```