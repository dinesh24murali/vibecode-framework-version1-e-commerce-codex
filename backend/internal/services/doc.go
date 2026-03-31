// Package services contains business logic for the e-commerce API.
// Services accept input from handlers, enforce domain rules, and interact
// with repositories for persistence. Services must not import the handlers
// package to keep the dependency direction one-way: handler → service → repository.
package services
