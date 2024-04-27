# BCC Nata de Coco Back-End

## API Specification

> [!NOTE]
> BaseURL is https://hackathon.miruza.my.id

### User API

```json
### User Sign-up (this will send email verif, please check your email)
POST {{baseURL}}/api/users

{
    "email": "{{email}}",
    "password": "{{password}}",
    "passwordConfirmation": "{{password}}",
    "name": "exquisite"
}


### User Log-in
POST {{baseURL}}/api/users/_login

{
    "email": "{{email}}",
    "password": "{{password}}",
    "rememberMe": true
}

### User Profile
GET {{baseURL}}/api/users/_login
Authorization: Bearer {{token}}

{
    "email": "mirza@gmail.com",
    "profilePicture": "{{url}}",
    "name": "Name"
}

```

### Animal API

```json
### What is this animal?
POST {{baseURL}}/api/users/_login
Content-Type: multipart/form-data

picture
lat
long

# Response
{
    "name": "string",
    "latin": "string",
    "countryOfOrigin": "string",
    "characteristics": ["string", "string", "string"],
    "category": "string",
    "lifespan": "string",
    "funfact": "string",
    "gotBonus": false
}
```