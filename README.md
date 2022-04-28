# Anansi
## An API supporting diversity in children's literature

Named after a fairytale spider who came to distribute all the stories in the world. Read more [here](https://en.wikipedia.org/wiki/Anansi#Akan-Ashanti_Anansi_stories)

## Tech stack

- Golang
- Serverless Framework
- AWS Lambda
- AWS DynamoDB

## API Documentation

### Auth API

This service enables the registration and authentication of users.

#### Registration

POST - https://p8aeakuo0l.execute-api.eu-central-1.amazonaws.com/auth/register

```typescript
export interface RegisterRequest {
	emailAddress: string;
	password: string;
}

export interface RegisterResponse {
	message: string;
}
```

#### Login

POST - https://p8aeakuo0l.execute-api.eu-central-1.amazonaws.com/auth/login

```typescript
export interface LoginRequest {
	emailAddress: string;
	password: string;
}

export interface LoginResponse {
	token: string;
	message: string;
}
```

#### Logout

POST - https://p8aeakuo0l.execute-api.eu-central-1.amazonaws.com/auth/logout

```typescript
export interface LogoutRequest {
	emailAddress: string;
	key: string;
}
```

#### Refresh

POST - https://p8aeakuo0l.execute-api.eu-central-1.amazonaws.com/auth/refresh

```typescript
export interface RefreshRequest {
	emailAddress: string;
	key: string;
}
```

#### Authorizing Requests

Once the user has been logged in a temporary access token will be returned. This token needs to be included in all requests made to non Auth-API endpoints (i.e. to the Books, Translations APIs etc...)

Set the header `Authorization: Bearer <Token>` on all non-auth requests.

### Language API

This service enables the retrieval of currently supported languages.

#### List Languages

GET - https://aih9h92cb8.execute-api.eu-central-1.amazonaws.com/language

```typescript
export type Language = string;

export type ListLanguagesResponse = Language[] | null;
```

### Books API

This service enables the creation and administration of book records, including their constituent translations and pages.

#### Create Book

POST - https://bkzq9e40g6.execute-api.eu-central-1.amazonaws.com/book

```typescript
export interface CreateBookRequest_Translation_Page {
	image_url: string;
	page_number: number;
}

export interface CreateBookRequest_Translation {
	localised_title: string;
	lang: Language;
	pages: CreateBookRequest_Translation_Page[] | null;
}

export interface CreateBookRequest {
	internal_title: string;
	authors: string;
	translations: CreateBookRequest_Translation[] | null;
}

export interface GetBookResponse_Translation_Page {
	page_id: PageID;
	image_url: string;
	page_number: number;
}

export interface GetBookResponse_Translation {
	id: TranslationID;
	localised_title: string;
	lang: Language;
	pages: GetBookResponse_Translation_Page[] | null;
}

export interface CreateBookResponse {
	id: BookID;
	internal_title: string;
	authors: string;
	translations: GetBookResponse_Translation[] | null;
}
```

#### Delete Book

DELETE - https://bkzq9e40g6.execute-api.eu-central-1.amazonaws.com/book/{id}

#### Get Book

GET - https://bkzq9e40g6.execute-api.eu-central-1.amazonaws.com/book/{id}

```typescript
export interface GetBookResponse_Translation_Page {
	page_id: PageID;
	image_url: string;
	page_number: number;
}

export interface GetBookResponse_Translation {
	id: TranslationID;
	localised_title: string;
	lang: Language;
	pages: GetBookResponse_Translation_Page[] | null;
}

export interface GetBookResponse {
	id: BookID;
	internal_title: string;
	authors: string;
	translations: GetBookResponse_Translation[] | null;
}
```

#### List Books

GET - https://bkzq9e40g6.execute-api.eu-central-1.amazonaws.com/book

```typescript
export type ListBooksResponse = GetBookResponse[] | null;
```

#### Update Book

PATCH - https://bkzq9e40g6.execute-api.eu-central-1.amazonaws.com/book/{id}

```typescript
export interface GetBookResponse_Translation_Page {
	page_id: PageID;
	image_url: string;
	page_number: number;
}

export interface GetBookResponse_Translation {
	id: TranslationID;
	localised_title: string;
	lang: Language;
	pages: GetBookResponse_Translation_Page[] | null;
}

export interface UpdateBookRequest {
	id: BookID;
	internal_title: string;
	authors: string;
	translations: GetBookResponse_Translation[] | null;
}

export interface UpdateBookResponse {
	id: BookID;
	internal_title: string;
	authors: string;
	translations: GetBookResponse_Translation[] | null;
}
```