package models

import "errors"

var (
	// Errors
	ErrorGetValueFromDatabase = errors.New("Error get value from database")
	ErrorZeroValue            = errors.New("Getting zero value")
	ErrorUnmarshallJSON       = errors.New("Promblem with decoding from JSON")
	ErrorIncorrectData        = errors.New("Incorrect data")
	ErrorDatabase             = errors.New("Error with database")
	LongMessage               = errors.New("Long message. Length bigger 2300 symbols")
	ErrorMatchingSite         = errors.New("Error with matching choose site. Please check correct name")
	ErrorGoogleSheet          = errors.New("Error with getting value from google sheet")
	ErrorDeleteTask           = errors.New("Error deleting task")
	// Other
	GenderMale   = "мужской"
	GenderFemale = "женский"
)

type RowObject struct {
	UserId int `json:"userId"`
	Object struct {
		Project           string `json:"project"`
		Link              string `json:"link"`
		Gender            int    `json:"gender"` // 1 - женский, 2 - мужской
		TextDescription   string `json:"text_description"`
		DateOfPublication string `json:"date_of_publication"`
	} `json:"object"`
}

func NewRowObject(userId int, project, link string, gender int, textDescription, dateOfPublication string) *RowObject {
	return &RowObject{
		UserId: userId,
		Object: struct {
			Project           string `json:"project"`
			Link              string `json:"link"`
			Gender            int    `json:"gender"`
			TextDescription   string `json:"text_description"`
			DateOfPublication string `json:"date_of_publication"`
		}{
			Project:           project,
			Link:              link,
			Gender:            gender,
			TextDescription:   textDescription,
			DateOfPublication: dateOfPublication,
		},
	}
}

type Person struct {
	Name    string
	Details struct { // This is an inline anonymous struct
		Age  int
		City string
	}
}

// NewPerson is a factory function acting as a constructor
func NewPerson(name string, age int, city string) *Person {
	return &Person{
		Name: name,
		Details: struct { // Initialize the inline anonymous struct
			Age  int
			City string
		}{
			Age:  age,
			City: city,
		},
	}
}

type AddedTask struct {
	Name   string
	Price  float64
	TaskId int64
}

// type TaskObject struct {
// 	Task_ID                  int
// 	Name                     string
// 	Descr                    string
// 	Link                     string
// 	Need_for_report          string
// 	Price                    float64
// 	Tarif_id                 int
// 	Folder_id                int
// 	Need_screen              bool
// 	Time_for_work            int
// 	Time_for_check           int
// 	Targeting_gender         int
// 	Targeting_geo_country_id int
// }

// func NewTaskObj(task_id int,
// 	name string,
// 	descr string,
// 	link string,
// 	need_for_report string,
// 	price float64,
// 	tarif_id int,
// 	folder_id int,
// 	need_screen bool,
// 	time_for_work int,
// 	time_for_check int,
// 	targeting_gender int,
// 	targeting_geo_country_id int) TaskObject {
// 	return TaskObject{
// 		Task_ID:                  task_id,
// 		Name:                     name,
// 		Descr:                    descr,
// 		Link:                     link,
// 		Need_for_report:          need_for_report,
// 		Price:                    price,
// 		Tarif_id:                 tarif_id,
// 		Folder_id:                folder_id,
// 		Need_screen:              need_screen,
// 		Time_for_work:            time_for_work,
// 		Time_for_check:           time_for_check,
// 		Targeting_gender:         targeting_gender,
// 		Targeting_geo_country_id: targeting_geo_country_id,
// 	}
// }
