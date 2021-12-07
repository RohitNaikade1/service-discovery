package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"service-discovery/database"
	"service-discovery/helpers"
	"service-discovery/models"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

func CronTask(credsid string) {

	var r models.Registration

	fmt.Println(credsid)

	col := database.RegistrationCollection()
	err := col.FindOne(context.TODO(), bson.M{"accounts.credsid": credsid}).Decode(&r)

	if err != nil {
		fmt.Println(err.Error())
	} else {

		for i := 0; i < len(r.Categories); i++ {

			for j := 0; j < len(r.Categories[i].Resource_info.Resources); j++ {

				x := r.Categories[i].Resource_info.Resources[j]
				fmt.Println(x)

				url := "http://localhost:8080/servicediscovery/cloudresources/azure/service/" + x + "?credsid=" + credsid
				method := "GET"

				client := &http.Client{}
				req, err := http.NewRequest(method, url, nil)
				if err != nil {
					fmt.Println(err)
				}

				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(body))

			}

		}
	}
}

type Set_DT struct {
	ti string
}
type Set_GT struct {
	year    int
	day     int
	hour    int
	min     int
	sec     int
	id      int
	month   time.Month
	credsid string
}

//Global{

var s *gocron.Scheduler
var cnt = 0
var s1 Set_DT
var s2 Set_GT

// }

func ExecuteCronJob() {
	fmt.Println("Inside execute cron job")
	s = gocron.NewScheduler(time.Now().Location())
	s.StartAsync()

}
func myTask() {

	fmt.Println("mytask")
	CronTask(s2.credsid)
	Collections(s2.credsid)
}

func IsCred(credsid string) (result bool) {
	var cred models.Credentials
	collection := database.CredentialCollection()
	err := collection.FindOne(context.Background(), bson.M{"credsid": credsid}).Decode(&cred)
	if err != nil {
		fmt.Println(err)
		result = false
	} else {
		fmt.Println("correct")
		result = true
	}
	return result
}

func SetJob(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
		s1 = Set_DT{}
		const DatelayOut = time.RFC3339

		var iscred bool
		cred := c.PostForm("credsid")
		iscred = IsCred(cred)
		fmt.Println(iscred)
		if iscred {
			s2.credsid = c.PostForm("credsid")
			if c.PostForm("Date_time") != "" {
				fmt.Println("date-time")
				//fmt.Println(c.PostForm("Date_time"))
				s1.ti = c.PostForm("Date_time")

				timeStampString := s1.ti

				timeStampDate, err := time.Parse(DatelayOut, timeStampString)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				cnt++

				tag := strconv.Itoa(cnt)

				hr, min, sec := timeStampDate.Clock()
				year, month, day := timeStampDate.Date()
				s2 := Set_GT{
					year:    year,
					month:   month,
					day:     day,
					hour:    hr,
					min:     min,
					sec:     sec,
					id:      cnt,
					credsid: cred,
				}

				t := time.Date(s2.year, s2.month, s2.day, s2.hour, s2.min, s2.sec, 0, time.Now().Location())

				j, _ := s.Every(1).Days().LimitRunsTo(1).At(t).Do(myTask)
				j.Tag(tag)

				fmt.Println("\nJob added, will run at:", t, "\nID:", tag)
				fmt.Println(s.Jobs())

			} else if c.PostForm("Periodic_hr") != "" {

				s1.ti = c.PostForm("Periodic_hr")
				timeStampString := s1.ti

				timeStampHour, err2 := time.Parse(time.Kitchen, timeStampString)
				if err2 != nil {
					fmt.Println(err2)
					os.Exit(1)
				}

				cnt++
				// var tag string
				tag := strconv.Itoa(cnt)
				// fmt.Println(time.Now())

				year, month, day := time.Now().Date()
				hr, min, sec := timeStampHour.Clock()
				s2 := Set_GT{
					year:  year,
					month: month,
					day:   day,
					hour:  hr,
					min:   min,
					sec:   sec,
					id:    cnt,
				}
				t := time.Date(s2.year, s2.month, s2.day, s2.hour, s2.min, s2.sec, 0, time.Now().Location())

				str := fmt.Sprintf("%d%s%d%s%d", s2.hour, ":", s2.min, ":", s2.sec)

				j, _ := s.Every(1).Day().At(t).Do(myTask)
				j.Tag(tag)

				fmt.Println("\nJob added it will run at :", str, "\nID:", tag)
				fmt.Println(s.Jobs())

			} else if c.PostForm("Periodic_min") != "" {

				s1.ti = c.PostForm("Periodic_min")
				timeStampString := s1.ti

				timeStampMin, err3 := strconv.Atoi(timeStampString)
				if err3 != nil {
					fmt.Println(err3)
					os.Exit(1)
				}

				cnt++

				tag := strconv.Itoa(cnt)

				year, month, day := time.Now().Date()
				hr, min, sec := time.Now().Clock()
				s2 := Set_GT{

					min: timeStampMin,
					id:  cnt,
				}

				t := time.Date(year, month, day, hr, min+s2.min, sec, 0, time.Now().Location())

				//j, _ := s.Every(str).Minutes().At(t).Do(myTask)
				//j.Tag(tag)

				// j1, _ := s.Every(str).Minutes().Do(myTask)
				// j1.Tag(tag)

				j, _ := s.Every(s2.min).Minutes().StartAt(t).Do(myTask)
				j.Tag(tag)

				fmt.Println("\nJob added it will run in ", s2.min, "minutes", "\nID:", tag)

				fmt.Println(s.Jobs())

			}
		} else {
			fmt.Println("Creds does not exist")
		}
	} else {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

}

func Task(c *gin.Context) {

	ch := c.PostForm("tasks")
	switch ch {

	case "list":
		var j []*gocron.Job
		j = s.Jobs()

		if len(j) == 0 {
			fmt.Println("\nThere are no jobs set")
		}
		if len(j) != 0 {

			fmt.Println("List of all jobs")
			for i := 0; i < len(j); i++ {

				fmt.Println("ID:", j[i].Tags())
				fmt.Println("Is scheduled to run at: ", j[i].NextRun())

			}

		}

		// fmt.Println(j)

	case "runall":
		fmt.Println("\nRunning all jobs,\n This might take some time")
		s.RunAllWithDelay(1)

	case "Delete":
		id := c.PostForm("ID")

		// var j []*gocron.Job
		// j = s.Jobs()
		// var c []string
		err := s.RemoveByTag(id)

		fmt.Println("Deleted")

		if err != nil {
			fmt.Println(err)
		}

	case "run_id":
		id := c.PostForm("ID")
		if id != "" {
			fmt.Println("\nRunning ID:", id)
			s.RunByTag(id)
			fmt.Println("Done")
		}

	}

}

//*****change active status *****//
var db = database.Database()

func GetOSType(data interface{}) string {
	p, e := json.Marshal(data)
	if e != nil {
		fmt.Println(e)
	}

	var d = models.Properties{}
	json.Unmarshal(p, &d)

	ostype := d.StorageProfile.OSDisk.OSType
	return ostype
}

func Collections(creds string) {
	fmt.Println("****************************************")
	arr := database.ListCollectionNames(db)
	//fmt.Println(len(arr))
	for i := 0; i < len(arr); i++ {
		fmt.Println("*******")
		fmt.Println(arr[i])
		SyncResources(arr[i], creds)
	}
}

func SyncResources(collection string, creds string) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourcesClient(helpers.SubscriptionID(creds), cred, nil)
	results := database.GetAllDocuments(db, collection)
	fmt.Println(collection + " Query Reult:")
	if collection == database.UserCollectionName() || collection == database.CredentialCollectionName() || collection == database.RegistrationCollectionName() {
		fmt.Println("not a resource")
	} else {
		for _, doc := range results {
			//fmt.Println(doc)
			name := doc["name"]
			n := fmt.Sprint(name)
			id := doc["id"]
			status := doc["status"]
			ID := fmt.Sprint(id)
			//fmt.Println("ID: ", ID)

			if status == "active" {
				resp, err := client.GetByID(context.Background(), ID, "2021-04-01", nil)
				if err != nil {
					fmt.Println("failed to obtain a response: ", err)

					var res []byte
					var er error
					if collection == "virtualmachines" {
						properties := doc["properties"]
						ostype := GetOSType(properties)
						r := models.VirtualMachine{Name: n, Type: collection, OsType: ostype}
						res, er = json.Marshal(r)
						if er != nil {
							fmt.Println(er)
						}
					} else {
						r := models.Resource{Name: n, Type: collection}
						res, er = json.Marshal(r)
						if er != nil {
							fmt.Println(er)
						}
					}

					DeActive(string(res))
					filter := bson.M{"name": n}
					data := bson.M{"status": "inactive"}
					update := bson.M{"$set": data}
					database.UpdateOne(db, collection, filter, update)

				}

				out, e := json.Marshal(resp)
				if e != nil {
					fmt.Println(err)
				}

				var d map[string]interface{}
				json.Unmarshal(out, &d)
				filter := bson.M{"name": n}
				update := bson.M{"$set": d}
				database.UpdateOne(db, collection, filter, update)
			}

		}
	}

}