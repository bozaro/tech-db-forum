package main

import "github.com/bozaro/tech-db-forum/tests/models"

type UserByNickname []models.User

func (a UserByNickname) Len() int           { return len(a) }
func (a UserByNickname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserByNickname) Less(i, j int) bool { return a[i].Nickname < a[j].Nickname }
