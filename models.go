package main

type User struct {
	Id       int64
	Email    string `sql:"size:255;not null;unique"`
	NickName string `sql:"size:255;not null;unique"`
	Active   bool
	Accounts []Account // One to many, A
}

type Account struct {
	Id      int64
	UserId  int64 // A
	GroupId int64 // B
	Credit  int64
	Records []Record // One to manay E
}

type Group struct {
	Id      int64
	Name    string    `sql:"size:255"`
	Members []Account // One to many, B
	Meals   []Meal    // One to many, C
}

type Meal struct {
	Id      int64
	GroupId int64    // C
	Records []Record // One to many, D
}

type Record struct {
	MealId    int64 // D
	AccountId int64 // E
	Pay       int64
	Owe       int64
	Credit    int64
}
