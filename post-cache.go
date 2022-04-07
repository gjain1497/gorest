package main

type UserCache interface{
	Set(key string, value User)
	Get(key string) *User
}