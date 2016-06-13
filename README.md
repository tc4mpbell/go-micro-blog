# go-micro-blog
Micro-blogging in Go.

## Usage

Run with two arguments (user and password) to create an 'account'.

    ./go-micro-blog admin somePassword

Run with no arguments to start the server.

Right now it tries to authenticate as 'taylor' (hard-coded), so you'd have to create a taylor account.

## Login

POST to localhost:8080/login with two arguments: username and password

If they match an account, it'll log you in and make a session file for you. This won't go away until you logout.

## Logout

POST to localhost:8080/login with one argument: username

It'll clear your session.

# Todo!

What, not production-ready? ;)

* Session handling, instead of nonsense.
