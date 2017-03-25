### Reblog - Some simple blog API built in Go

Features 

- [x] User authentication via JWT
- [x] Multi tenant
  - [x] Admin can add new collaborators
  - [x] Collaborators can sign up after admin sends them a link to signup
  - [x] Posts can be created by the admin and collaborators
  - [x] Admin can delete collaborators
  - [x] Admin can delete posts
  - [x] Admin can mark a post as unpublished


> The admin user can be manually created by running an insert query into the `users` table with the type field set to 1.

  
