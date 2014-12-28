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

/* -------------------------- *
 *   For Dictionary project   *
 * -------------------------- */

type Dict struct {
	Id    int64
	Name  string     `sql:"size:255;unique;not null;"`
	Words []UserWord `gorm:"many2many:user_languages;`
}

type UserWord struct {
	Id     int64
	DictId int64
	WordId int64
	Word   string `sql:"size:255;not null"`
	Extra  string `sql:"size:255"`
}

type Word struct {
	Id    int64
	Word  string `sql:"size:255;unique;not null"`
	Trans []Trans
}

type Trans struct {
	Id     int64
	WordId int64
	Trans  string
}
