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

POST - https://81qrzgok36.execute-api.eu-central-1.amazonaws.com/auth/register

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

POST - https://81qrzgok36.execute-api.eu-central-1.amazonaws.com/auth/login

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

```typescript
export interface LogoutRequest {
	emailAddress: string;
	key: string;
}
```

#### Refresh

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

This service enables the retrieval of currently supported languages

#### List Languages

GET - https://aih9h92cb8.execute-api.eu-central-1.amazonaws.com/language

```typescript
export type Language = string;

export type ListLanguagesResponse = Language[] | null;
```

### Books API

This service enables the creation and administration of book records.

#### Create Book

POST - https://8oufr7mu4f.execute-api.eu-central-1.amazonaws.com/book

```typescript
export interface CreateBookRequest {
	internal_title: string;
	authors: string;
}

export interface CreateBookResponse {
	id: string;
	internal_title: string;
	authors: string;
}
```

#### Delete Book

DELETE - https://8oufr7mu4f.execute-api.eu-central-1.amazonaws.com/book/{id}

#### Get Book

GET - https://8oufr7mu4f.execute-api.eu-central-1.amazonaws.com/book/{id}

```typescript
export interface GetBookResponse {
	id: string;
	internal_title: string;
	authors: string;
}
```

#### List Books

GET - https://8oufr7mu4f.execute-api.eu-central-1.amazonaws.com/book

```typescript
export interface Book {
	id: string;
	internal_title: string;
	authors: string;
}

export type ListBooksResponse = Book[] | null;
```

#### Update Book

PATCH - https://8oufr7mu4f.execute-api.eu-central-1.amazonaws.com/book/{id}

```typescript
export interface UpdateBookRequest {
	internal_title: string;
	authors: string;
}

export interface UpdateBookResponse {
	id: string;
	internal_title: string;
	authors: string;
}
```

### Translations API

This service enables the creation and administration of book translations.

#### Create Translation

POST - https://fgxdnldga8.execute-api.eu-central-1.amazonaws.com/translation

```typescript
export interface CreateTranslationRequest {
	book_id: string;
	localised_title: string;
	language: string;
}

export interface CreateTranslationResponse {
	id: string;
	book_id: string;
	localised_title: string;
	language: Language;
}
```

#### Delete Translation

DELETE - https://fgxdnldga8.execute-api.eu-central-1.amazonaws.com/translation/{id}

#### Get Translation

GET - https://fgxdnldga8.execute-api.eu-central-1.amazonaws.com/translation/{id}

```typescript
export interface GetTranslationResponse {
	id: string;
	book_id: string;
	localised_title: string;
	language: Language;
}
```

#### List Translations

GET - https://fgxdnldga8.execute-api.eu-central-1.amazonaws.com/translation

```typescript
export type Language = string;
export interface Translation {
	id: string;
	book_id: string;
	localised_title: string;
	language: Language;
}

export type ListBooksResponse = Book[] | null;
```

#### Update Translation

PATCH - https://fgxdnldga8.execute-api.eu-central-1.amazonaws.com/translation/{id}

```typescript
export interface UpdateTranslationRequest {
	book_id: string;
	localised_title: string;
	language: string;
}

export interface UpdateTranslationResponse {
	id: string;
	book_id: string;
	localised_title: string;
	language: Language;
}
```