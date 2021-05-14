package storage

import (
	"context"
	"fmt"
	"main/dict"
	"math"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Data is the root structure of everything stored in the database.
type Data struct {
	Time      string      `json:"time"`
	Container interface{} `json:"container"`
}

// Firebase is used as the global firebase connection.
var Firebase Database

// Database structure stores information about firebase database.
//
// Functionality: Setup, clean, Add, Get, GetAll, Delete
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
	//start database cleaner
	go database.clean()
	return nil
}

// clean database for some collections.
func (database *Database) clean() {
	var itemsCleaned int
	dict.MutexState.Lock()
	//clean ticketmaster events
	iter := database.Client.Collection(dict.EVENT_COLLECTION).Documents(database.Ctx)
	for {
		//go to next element in array and break loop if there are no elements, branch if an error occurred
		elem, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			fmt.Printf(
				"%v {\n\tError while cleaning database.\n\tRaw error: %v\n}\n",
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		} else {
			database.Delete(dict.EVENT_COLLECTION, elem.Ref.ID)
			itemsCleaned++
		}
	}
	//clean coordinates
	iter = database.Client.Collection(dict.LOCATION_COLLECTION).Documents(database.Ctx)
	for {
		//go to next element in array and break loop if there are no elements, branch if an error occurred
		elem, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			fmt.Printf(
				"%v {\n\tError while cleaning database.\n\tRaw error: %v\n}\n",
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		} else {
			database.Delete(dict.LOCATION_COLLECTION, elem.Ref.ID)
			itemsCleaned++
		}
	}
	fmt.Printf(
		"%v {\n\tSuccsesfully cleaned database.\n\tRemoved items: %v\n}\n",
		time.Now().Format("2006-01-02 15:04:05"), strconv.Itoa(itemsCleaned),
	)
	dict.MutexState.Unlock()
	//put program to sleep for 12 hours before cleaing again
	nextTime := time.Now().Truncate(time.Hour)
	nextTime = nextTime.Add(time.Duration(12) * time.Hour)
	time.Sleep(time.Until(nextTime))
	database.clean()
}

// Add a new webhook to database.
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
		_, err = database.Client.Collection(name).Doc(id).Update(database.Ctx, []firestore.Update{{
			Path:  "Container.ID",
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
	return data.Time, id, nil
}

// Get a webhook from the database.
func (database *Database) Get(name string, id string) (map[string]interface{}, bool) {
	iter, err := database.Client.Collection(name).Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, false
	}
	data := iter.Data()
	return data, true
}

// Get all webhooks from the database.
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

// Delete specific webhook.
func (database *Database) Delete(webhookdb string, id string) error {
	_, err := database.Client.Collection(webhookdb).Doc(id).Delete(database.Ctx)

	return err
}

// Check date of data in database.
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
