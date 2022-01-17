package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/helpers"
	"service-discovery/models"
	"strconv"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

func GetResourcesCount(Categories []models.Category_info) (count int) {
	count = 0
	for i := 0; i < len(Categories); i++ {

		for j := 0; j < len(Categories[i].Resource_info.Resources); j++ {
			count++
		}
	}
	return count
}

func CronTask(credsid string) {
	start := time.Now()
	var wg = sync.WaitGroup{}

	var r models.Registration

	//fmt.Println(credsid)

	collection := database.RegistrationCollection()
	err := collection.FindOne(context.TODO(), bson.M{"accounts.credsid": credsid}).Decode(&r)

	count := GetResourcesCount(r.Categories)

	//	fmt.Println("Count: ", count)

	wg.Add(count)
	if err != nil {
		Logger.Error(err.Error())
	} else {
		for i := 0; i < len(r.Categories); i++ {

			for j := 0; j < len(r.Categories[i].Resource_info.Resources); j++ {
				resource := r.Categories[i].Resource_info.Resources[j]
				go GetResourceData(resource, credsid, &wg)
			}
		}

	}

	wg.Wait()
	elapsed := time.Since(start)

	Logger.Info("Sync took " + elapsed.String())

}

func GetResourceData(resource string, credsid string, wg *sync.WaitGroup) {
	port := env.GetEnvironmentVariable("PORT")
	url := "http://localhost:" + port + "/servicediscovery/cloudresources/azure/service/" + resource + "?credsid=" + credsid
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		Logger.Error(err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger.Error(err.Error())
	}

	Logger.Info(string(body))
	time.Sleep(time.Duration(10 * time.Millisecond))
	wg.Done()
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
	Logger.Info("Inside execute cron job")
	s = gocron.NewScheduler(time.Now().Location())
	s.StartAsync()

}
func myTask() {
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

		//fmt.Println(iscred)

		if iscred {
			s2.credsid = c.PostForm("credsid")
			if c.PostForm("Date_time") != "" {
				//fmt.Println("date-time")
				s1.ti = c.PostForm("Date_time")

				timeStampString := s1.ti

				timeStampDate, err := time.Parse(DatelayOut, timeStampString)

				if err != nil {
					Logger.Error(err.Error())
					//os.Exit(1)
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

				Logger.Info("\nJob added, will run at:" + t.String())
				Logger.Info("ID:" + tag)
				c.JSON(http.StatusOK, bson.M{"Job added, it will run at": t, "ID": tag})
				//	fmt.Println(s.Jobs())

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

				Logger.Info("Job added it will run at : " + str)
				Logger.Info("\nID:" + tag)
				//fmt.Println(s.Jobs())

			} else if c.PostForm("Periodic_min") != "" {

				s1.ti = c.PostForm("Periodic_min")
				timeStampString := s1.ti

				timeStampMin, err := strconv.Atoi(timeStampString)
				if err != nil {
					Logger.Info(err.Error())
					//os.Exit(1)
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

				//fmt.Println(s.Jobs())

			}
		} else {
			Logger.Error("Creds does not exist")
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
			Logger.Info("There are no jobs set")
		}
		if len(j) != 0 {

			Logger.Info("List of all jobs")
			for i := 0; i < len(j); i++ {

				fmt.Println("ID: ", j[i].Tags())
				fmt.Println("Is scheduled to run at: ", j[i].NextRun())

			}

		}

		// fmt.Println(j)

	case "runall":
		Logger.Info("\nRunning all jobs,\n This might take some time")
		s.RunAllWithDelay(1)

	case "Delete":
		id := c.PostForm("ID")

		// var j []*gocron.Job
		// j = s.Jobs()
		// var c []string
		err := s.RemoveByTag(id)

		Logger.Info("Deleted")

		if err != nil {
			Logger.Error(err.Error())
		}

	case "run_id":
		id := c.PostForm("ID")
		if id != "" {
			Logger.Info("\nRunning ID:" + id)
			s.RunByTag(id)
			Logger.Info("Done")
		}

	}

}

//*****change active status *****//

func GetOSType(data interface{}) string {
	p, err := json.Marshal(data)
	if err != nil {
		Logger.Error(err.Error())
	}

	var d = models.Properties{}
	json.Unmarshal(p, &d)

	ostype := d.StorageProfile.OSDisk.OSType
	return ostype
}

func Collections(creds string) {

	start := time.Now()
	var wg = sync.WaitGroup{}

	arr := database.ListCollectionNames()
	//fmt.Println(len(arr))
	wg.Add(len(arr))
	for i := 0; i < len(arr); i++ {
		Logger.Info(arr[i])
		go SyncResources(arr[i], creds, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)
	Logger.Info("Collections took " + elapsed.String())
}

func SyncResources(collection string, creds string, wg *sync.WaitGroup) {
	var Resp armresources.ResourcesGetByIDResponse
	var Err error
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		Logger.Error("failed to obtain a credential: " + err.Error())
	}

	client := armresources.NewResourcesClient(helpers.SubscriptionID(creds), cred, nil)
	results := database.ReadAll(collection)
	Logger.Info(collection + " Query Reult:")
	//if collection == database.UserCollectionName() || collection == database.CredentialCollectionName() || collection == database.RegistrationCollectionName() {
	//	fmt.Println("not a resource")
	//} else {
	for _, doc := range results {
		//fmt.Println(doc)
		name := doc["name"]
		n := fmt.Sprint(name)
		id := doc["id"]
		status := doc["status"]
		ID := fmt.Sprint(id)
		//fmt.Println("ID: ", ID)

		if status == "active" {
			fmt.Println("Inside active")
			if collection == "databases" || collection == "servers" {
				resp, err := client.GetByID(context.Background(), ID, "2021-08-01-preview", nil)
				Resp = resp
				Err = err
			} else {
				resp, err := client.GetByID(context.Background(), ID, "2021-04-01", nil)
				Resp = resp
				Err = err
			}

			if Err != nil {
				Logger.Error("failed to obtain a response: " + err.Error())

				var res []byte
				var er error
				if collection == "virtualmachines" {
					properties := doc["properties"]
					ostype := GetOSType(properties)
					r := models.VirtualMachine{Name: n, Type: collection, OsType: ostype}
					res, er = json.Marshal(r)
					if er != nil {
						Logger.Error(er.Error())
					}
				} else {
					r := models.Resource{Name: n, Type: collection}
					res, er = json.Marshal(r)
					if er != nil {
						Logger.Error(er.Error())
					}
				}

				DeActive(string(res))
				filter := bson.M{"name": n}
				data := bson.M{"status": "inactive"}
				update := bson.M{"$set": data}
				database.Update(collection, filter, update)

			}

			out, err := json.Marshal(Resp)
			if err != nil {
				Logger.Error(err.Error())
			}

			var d map[string]interface{}
			json.Unmarshal(out, &d)
			filter := bson.M{"name": n}
			update := bson.M{"$set": d}
			database.Update(collection, filter, update)
		}

	}
	//}
	wg.Done()
}
