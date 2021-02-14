package main

import (
	"TestProject/Server/GolangServer/collection"
	"TestProject/Server/GolangServer/confighelper"
	"TestProject/Server/GolangServer/dbhelper"
	"TestProject/Server/GolangServer/model"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"regexp"

	"github.com/labstack/echo/v4"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	RegExp       = `s/<[a-zA-Z\/][^>]*>//g`
	layoutISO    = "2006-01-02T15:04"
	layoutUS     = "January 2, 2006"
	secondRegExp = `[^a-z ^A-Z]`
)

var wg sync.WaitGroup
var primeNumberList []int

func main() {

	confighelper.InitViper()
	e := echo.New()
	e.POST("/getColumns", getColumns)
	e.POST("/getWordCountService", GetWordCountService)
	e.POST("/getLastDayOfInputDate", GetLastDayOfInputDate)
	e.POST("/getPrimeNumberListService", GetPrimeNumberListService)
	e.POST("/updateCandidateRecordService", UpdateCandidateRecordService)
	e.POST("/deleteRecordService", DeleteRecordService)
	e.GET("/getAllRecordsService", GetAllRecordsService)
	e.Logger.Fatal(e.Start(":4000"))
}

//This service is return the colums and rows by processing input
func getColumns(c echo.Context) error {
	body, _ := GetRequestBodyJson(c)
	columnToStart := body.Get("columnToStart").String()
	numRows := body.Get("numRows").Int()
	numCols := body.Get("numCols").Int()
	res := GetColList(columnToStart, int(numRows), int(numCols))
	return c.JSON(http.StatusOK, res)
}

//Function will return the last day of input date
func GetLastDayOfInputDate(c echo.Context) error {
	body, _ := GetRequestBodyJson(c)

	date := body.Get("Date").String()
	t, _ := time.Parse(layoutISO, date)
	now := time.Now()
	currentYear, currentMonth := t.Year(), t.Month()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	dateObj, _ := sjson.Set("{}", "Date", date)
	dateObj, _ = sjson.Set(dateObj, "LastDayOfMonth", lastOfMonth.Day())
	ldData := gjson.Parse(dateObj)
	return c.JSON(http.StatusOK, ldData.Value())
}

//Read the URL and retun the json string
func GetRequestBodyJson(c echo.Context) (body gjson.Result, err error) {
	bb, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Error Reading Request Body")
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(bb), nil
}

func GetColList(columnToStart string, numRows, numCols int) []string {
	colList := []string{}
	colList = append(colList, columnToStart)

	// fixed := 0
	pivot := 0
	// if len(columnToStart) > 1 {
	// 	fixed = len(columnToStart) - 2
	// }
	//ABC
	if len(columnToStart) > 0 {
		pivot = len(columnToStart) - 1
	}

	for true {
		if len(colList) == numCols*numRows {
			break
		}
		currentStr := colList[len(colList)-1]
		if rune(currentStr[pivot]) == 'Z' {
			currentStr = handleEndCase(currentStr)
			colList = append(colList, currentStr)
			pivot = len(currentStr) - 1
			// fixed = len(currentStr) - 2
			continue
		}
		colList = append(colList, currentStr[:pivot]+string(rune(currentStr[pivot])+1))
		currentStr = currentStr[:pivot] + string(rune(currentStr[pivot])+1)
	}
	fmt.Println(colList)

	finalRes := []string{}
	for i := 0; i < len(colList); i++ {
		fmt.Println(i, i+numCols)
		finalRes = append(finalRes, strings.Join(colList[i:i+numCols], " "))
		i += numCols - 1
	}
	// fmt.Println(strings.Join(finalRes, ","))

	return finalRes
}

//Handle the End case
//AA ZZ
//AAAAA
func handleEndCase(cs string) string {
	pre := ""
	zCount := 0
	for i := len(cs) - 1; i >= 0; {
		if cs[i] == 'Z' {
			zCount++
			i--
			continue
		}
		pre = cs[:i] + string(rune(cs[i])+1)
		i--
		break
	}
	for i := 0; i < zCount; i++ {
		pre = pre + "A"
	}
	if zCount == len(cs) {
		pre = pre + "A"
	}
	return pre
}

//To get the word count
func GetWordCountService(c echo.Context) error {

	body, _ := GetRequestBodyJson(c)
	url := body.Get("url").String()
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f GetWordCountService", r)
		}
	}()
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error While Reading data from URL GetWordCountService", r)
		}
	}()

	r := regexp.MustCompile(RegExp)
	stripped := strip.StripTags(string(html))
	stripped = r.ReplaceAllString(stripped, "")
	ss := regexp.MustCompile(secondRegExp)
	stripped = ss.ReplaceAllString(stripped, " ")
	totalCount := 0
	list := []interface{}{}

	for index, element := range wordCount(stripped) {
		totalCount++
		ldDetails, _ := sjson.Set("{}", "word", index)
		ldDetails, _ = sjson.Set(ldDetails, "count", element)
		ldData := gjson.Parse(ldDetails)
		list = append(list, ldData.Value())

		// fmt.Println(index, "=>", element)
	}
	// fmt.Println(":::::::::::::::::::::::", list)
	return c.JSON(http.StatusOK, list)
}

//To get the prime number list
func GetPrimeNumberListService(c echo.Context) error {
	primeNumberList = []int{}
	body, _ := GetRequestBodyJson(c)
	number := body.Get("number").Int()
	for i := 0; i <= int(number); i++ {
		wg.Add(1)
		go CheckPrime(i)
		wg.Wait()
	}
	return c.JSON(http.StatusOK, primeNumberList)
}

//to delete the input login Id details PERMENTLY
func DeleteRecordService(c echo.Context) error {
	body, _ := GetRequestBodyJson(c)
	loginId := body.Get("loginId").String()

	_, serviceCallError := DeleteRecordDAO(loginId)
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, true)
}

//Update the user details
func UpdateCandidateRecordService(c echo.Context) error {

	EmployeeProfile := model.EmployeeProfile{}
	bindError := c.Bind(&EmployeeProfile)
	if bindError != nil {
		fmt.Println("BIND ERROR")
		return c.JSON(http.StatusInternalServerError, bindError)
	}
	loginId := EmployeeProfile.LoginId

	_, serviceCallError := UpdateService(EmployeeProfile, loginId)
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, true)
}
func GetAllRecordsService(c echo.Context) error {

	EmployeeProfileList := gjson.Result{}

	EmployeeProfileList, serviceCallError := GetAllRecordListService()
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, EmployeeProfileList.Value())
}

// This method get all  record  mapping List by calling DAO method.
func GetAllRecordListService() (gjson.Result, error) {
	return GetAllRecordListDAO()
}

//This Method Get all  recordList, Retrives data from MongoDB.
func GetAllRecordListDAO() (gjson.Result, error) {

	// db, ctx, err := dbhelper.GetMongoDB("dbIP", "DB_PORT", "DBNAME", "USENAME", "PASS", true)
	fmt.Println("ssss", confighelper.GetConfig("DBNAME"))
	db, ctx, err := dbhelper.GetMongoDB(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), "", "", false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return gjson.Result{}, err
	}

	selector := bson.M{"isDeleted": false}

	cursor, err := db.Collection(collection.EMPLOYEE_PROFILE).Find(ctx, selector)
	if err != nil {
		log.Fatal(err)
	}
	var records []bson.M
	if err = cursor.All(ctx, &records); err != nil {
		log.Fatal(err)
	}

	err = cursor.Decode(&records)
	if err != nil {
		log.Print(err)

	}
	bs, err := json.Marshal(records)
	if err != nil {
		log.Print(err)

		// continue
	}
	return gjson.ParseBytes(bs), nil
}

// This method delete personal information by calling DAO method.
func DeleteRecordDAO(loginId string) (bool, error) {

	flag, updateServiceError := DeleteDAO(loginId)
	if updateServiceError != nil {
		fmt.Println(" UpdateService Error")
		return false, updateServiceError
	}
	return flag, nil
}

// This method update personal information by calling DAO method.
func UpdateService(employeeProfileObj model.EmployeeProfile, loginId string) (bool, error) {

	flag, updateServiceError := UpdateDAO(employeeProfileObj, loginId)
	if updateServiceError != nil {
		fmt.Println(" UpdateService Error")
		return false, updateServiceError
	}
	return flag, nil
}

// This Method Delete personal information, Update data in MongoDB.
func DeleteDAO(loginId string) (bool, error) {
	db, ctx, err := dbhelper.GetMongoDB(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), "sampleTestDatabase", "", "", false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return false, err
	}
	selector := bson.M{"loginId": loginId}
	_, err = db.Collection(collection.EMPLOYEE_PROFILE).DeleteOne(ctx, selector)
	if err != nil {
		log.Fatal("Error While Deleting Record", err)
	}
	return true, nil
}

// This Method update personal information, Update data in MongoDB.
func UpdateDAO(templateModelObj model.EmployeeProfile, loginId string) (bool, error) {
	// db, ctx, err := dbhelper.GetMongoDB("dbIP", "DB_PORT", "DBNAME", "USENAME", "PASS", true)
	db, ctx, err := dbhelper.GetMongoDB(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), "sampleTestDatabase", "", "", false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return false, err
	}

	opts := options.Update().SetUpsert(true)
	selector := bson.M{"loginId": loginId}
	updator := bson.M{"$set": bson.M{
		"fullName":  templateModelObj.FullName,
		"isEnabled": templateModelObj.IsEnabled,
		"isDeleted": templateModelObj.IsDeleted}}
	result, err := db.Collection(collection.EMPLOYEE_PROFILE).UpdateOne(ctx, selector, updator, opts)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return false, nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	return true, nil
}
func CheckPrime(number int) {
	isPrime := true
	if number == 0 || number == 1 {
		wg.Done()
		// fmt.Printf(" %d is not a  prime number\n", number)
	} else {
		for i := 2; i <= number/2; i++ {
			if number%i == 0 {
				wg.Done()
				// fmt.Printf(" %d is not a  prime number\n", number)
				isPrime = false
				break
			}
		}
		if isPrime == true {
			// fmt.Printf("%d\t", number)
			primeNumberList = append(primeNumberList, number)
			wg.Done()
		}
	}
}

func wordCount(str string) map[string]int {
	wordList := strings.Fields(str)
	counts := make(map[string]int)
	for _, word := range wordList {
		_, ok := counts[word]
		if ok {
			counts[word]++
		} else {
			counts[word] = 1
		}
	}
	return counts
}
