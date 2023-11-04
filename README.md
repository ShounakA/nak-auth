# nak-auth

Welcome to nak-auth. It's a personal project where I am trying to implement OAuth 2.0, and build myself an Authentication and Authorization Server. 

## Currently it supports:
 - session-based authentication
   - to view and manage your clients
 - client crendentials grant
 - authorization code grant
 - refresh token grant

## Future Goals
 - implement Open ID 
 - be able to create an account
 - be able to add my own client
    - manage my clients
    - revoke a secret
    - get token sign key
- be add permissions via claims
  - manage permissions via claims
  
## How to start the project

### Prerequiste 
 - Install Air

## Runing Dev Mode

```bash
$ make run
```
This will start up the dev server. 
It also watches for file changes, but it wont live-reload your webpage. This mean you need to reload the page on any template changes.

## See the makefile for more tooling info