package db

import (
	"context"
	"math"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Data struct {
	Time      string      `json:"time"`
	Container interface{} `json:"container"`
}

// DB stores information about the firebase database.
//
// Used in notification structure.
var DB Database

// Database structure stores information about firebase database.
//
// Functionality: Setup, Add, Get, Delete
type Database struct {
	Ctx    context.Context
	Client *firestore.Client
}

// Setup sets up database.
func (database *Database) Setup() error {
	//connect to firebase with key and branch if an error occurred
	database.Ctx = context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(database.Ctx, nil, sa)
	if err != nil {
		return err
	}
	//open firestore and branch if an error occurred
	database.Client, err = app.Firestore(database.Ctx)
	if err != nil {
		return err
	}
	// //update Notifications with data from database and branch if and error occurred
	// err = database.Get()
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Add adds a new webhook to database.
func (database *Database) Add(name string, id string, data Data) (string, string, error) {
	data.Time = time.Now().Format(time.RFC822)
	if id == "" {
		//add data to database and get a UUID from firebase and branch if an error occurred
		docRef, _, err := database.Client.Collection(name).Add(database.Ctx, data)
		if err != nil {
			return "", "", err
		}
		id = docRef.ID
		//update webhook ID in database with UUID and branch if an error occurred
		_, err = database.Client.Collection(name).Doc(id).Update(database.Ctx, []firestore.Update {{
			Path: "Container.ID",
			Value: id,
		}})
		if err != nil {
			return "", "", err
		}
	} else {
		//add data to database and get a UUID from firebase and branch if an error occurred
		_, err := database.Client.Collection(name).Doc(id).Set(database.Ctx, data)
		if err != nil {
			return "", "", err
		}
	}
	// //update Notifications with data from database and branch if and error occurred
	// err = database.Get()
	// if err != nil {
	// 	return err
	// }
	return data.Time, id, nil
}

func (database *Database) Get(name string, id string) (map[string]interface{}, bool) {
	iter, err := database.Client.Collection(name).Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, false
	}
	data := iter.Data()
	return data, true
}

func (database *Database) GetAll(name string) ([]map[string]interface{}, error) {
	var arrData []map[string]interface{}
	iter := database.Client.Collection(name).Documents(database.Ctx)
	for {
		//go to next element in array and break loop if there are no elements, branch if an error occurred
		elem, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		arrData = append(arrData, elem.Data())
	}
	return arrData, nil
}

// // Get gets all webhooks from database
// func (database *Database) Get() error {
// 	//clear current webhooks stored in Notifications
// 	Notifications = make(map[string]Notification)
// 	//iterate through database and get each webhook
// 	iter := database.Client.Collection("notification").Documents(database.Ctx)
// 	var notification Notification
// 	for {
// 		//go to next element in array and break loop if there are no elements, branch if an error occurred
// 		elem, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		} else if err != nil {
// 			return err
// 		}

// elem.Ref.ID for ID

// 		//convert data from interface and set in structure
// 		data := elem.Data()
// 		notification.ID = fmt.Sprintf("%v", data["ID"])
// 		notification.URL = fmt.Sprintf("%v", data["URL"])
// 		notification.Timeout = data["Timeout"].(int64)
// 		notification.Information = fmt.Sprintf("%v", data["Information"])
// 		notification.Country = fmt.Sprintf("%v", data["Country"])
// 		notification.Trigger = fmt.Sprintf("%v", data["Trigger"])
// 		//add structure to map
// 		Notifications[notification.ID] = notification
// 	}
// 	return nil
// }

// Delete deletes specific webhook.
func (database *Database) Delete(id string) error {
	//get only element that has the same ID as specified and branch if an error occurred
	iter := database.Client.Collection("notification").Where("ID", "==", id).Documents(database.Ctx)
	elem, err := iter.Next()
	if err != nil {
		return err
	}
	//delete webhook and branch if an error occurred
	_, err = elem.Ref.Delete(database.Ctx)
	if err != nil {
		return err
	}
	// //update Notifications with data from database and branch if and error occurred
	// err = database.Get()
	// if err != nil {
	// 	return err
	// }
	return nil
}

func CheckDate(dataTime string, expectedHours int) (bool, error) {
	then, err := time.Parse(time.RFC822, dataTime)
	if err != nil {
		return true, err
	}
	//get current time and subtract inputted date
	currentTime := time.Now()
	diffTime := currentTime.Sub(then)
	//convert the difference to integer of hours
	diffHours := int(math.Floor(diffTime.Hours()))
	return diffHours <= expectedHours, nil
}

//CheckIfDateOfEventPassed checks if the date of the event has passed, and if so returns true
func CheckIfDateOfEventPassed(dateOfEvent time.Time) bool {
	currentTime := time.Now()

	if currentTime.Before(dateOfEvent) {
		return false
	} else {
		return true
	}
}
