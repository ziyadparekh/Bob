package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"

	jenkins "github.com/ziyadparekh/golang-jenkins"
)

type User struct {
	Username string `json:username`
	ApiToken string `json:apitoken`
}

func NewBob() (*User, error) {
	// current system user (need to get home dir)
	usr, _ := user.Current()
	// Pointer to User struct
	var data *User
	// Check if jenkins config file exists
	if _, err := os.Stat(usr.HomeDir + "/.bob/config.json"); os.IsNotExist(err) {
		fmt.Println("Config file does not exist, lets create one!")
		if _, err := CreateConfigFile(); err != nil {
			return nil, err
		}
	}

	// Read file
	file, err := ioutil.ReadFile(usr.HomeDir + "/.bob/config.json")
	if err != nil {
		panic(err)
	}
	// Unmarshal json into User Struct
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	// Return new instance of User and nil for error
	return data, nil
}

func CreateConfigFile() (t bool, e error) {
	var (
		username string
		apitoken string
		FileData map[string]string
	)
	fmt.Println("Enter jenkins username: ")
	if _, err := fmt.Scanf("%s", &username); err != nil {
		return false, err
	}
	fmt.Println("Enter jenkins apitoken: ")
	if _, err := fmt.Scanf("%s", &apitoken); err != nil {
		return false, err
	}

	usr, _ := user.Current()
	fmt.Println("-----> Creating directory " + usr.HomeDir + "/.bob")
	if err := os.MkdirAll(usr.HomeDir+"/.bob", 0777); err != nil {
		return false, err
	}

	FileData["username"] = username
	FileData["apitoken"] = apitoken

	j, jerr := json.MarshalIndent(&FileData, "", "  ")
	if jerr != nil {
		fmt.Println("jerr:", jerr.Error())
	}
	fmt.Println("-----> Writing file config.json to directory")
	ferr := ioutil.WriteFile(usr.HomeDir+"/.bob/config.json", j, 0777)
	if ferr != nil {
		return false, ferr
	}
	fmt.Println("-----> Done!")
	return true, nil
}

func RunJob(d *jenkins.Auth, s, b, t, env string, o bool) {
	// Instantiate Jenkins
	j := jenkins.NewJenkins(d, DibsyJenkins)
	// Create job struct
	job := jenkins.Job{
		Name: Jobs[s],
	}

	fmt.Println("Preparing to build job")
	fmt.Println(job.Name)

	// Create params for the build
	params := url.Values{}
	// Jenkins host key is SERVER_NAME
	// Jenkins server name is deathstar.1stdibs.com
	// Jenking branch key is BRANCH_NAME
	// b is the branch arg
	switch t {
	case "services":
		params.Set(JenkinsBranchKey, b)
		params.Set(JenkinsHostKey, JenkinsServerName)
	case "client":
		params.Set(JenkinsClientHost, JenkinsClientServerName)
		params.Set(JenkinsBranchKey, b)
		params.Set(JenkinsClientEnv, env)
	}

	// Try get info about the next build
	var i jenkins.Job
	i, e := j.GetJob(job.Name)

	if e != nil {
		fmt.Println(e)
	}

	fmt.Println(fmt.Sprintf("-----> Next build number is %d", i.NextBuildNumber))
	// Link to the job page and not the build because otherwise will get 404 page
	// until jenkins registers the job
	url := fmt.Sprintf("%s", i.Url)
	fmt.Println(fmt.Sprintf("-----> Next build url is %s", url))

	// BUILD!
	err := j.Build(job, params)

	if err != nil {
		fmt.Println(err)
	}
	// If the open url flag is set open the url
	// WARNING:: if you build all services and set this to true
	// it will open a TON of pages
	if o == true {
		exec.Command("open", url).Start()
	}

}

/**
 * Builds a single service. Takes the service name, branch and
 * open flag (to open thebrowser window or not)
 * @param  {[type]} u User)         BuildService(s, b string, o bool [description]
 * @return {[type]}   [description]
 */
func (u User) BuildService(s, b string, o bool) {
	service, err := FormatService(s, "services")
	if err != nil {
		log.Fatal(err)
		return
	}

	branch := FormatBranch(b)
	if _, err := EnsureBranchExists(service, branch); err != nil {
		log.Fatal("Requested branch does not exist on remote")
		return
	}

	info := fmt.Sprintf("Attempting to build %s service on branch %s", service, branch)
	fmt.Println(info)

	var data = jenkins.Auth{
		Username: u.Username,
		ApiToken: u.ApiToken,
	}

	RunJob(&data, service, branch, "services", "", o)
}

/**
 * Builds all services in separate go routines
 * @param  {[type]} u User)         BuildAllServices(b string, o bool [description]
 * @return {[type]}   [description]
 */
func (u User) BuildAllServices(b string, o bool) {
	fmt.Println("Attempting to build all Services")

	var data = jenkins.Auth{
		Username: u.Username,
		ApiToken: u.ApiToken,
	}

	branch := FormatBranch(b)

	for _, service := range Services {
		if _, err := EnsureBranchExists(service, branch); err != nil {
			log.Fatal("Requested branch does not exist on remote")
			return
		}
		info := fmt.Sprintf("Attempting to build %s service on branch %s", service, branch)
		fmt.Println(info)
		go RunJob(&data, service, branch, "services", "", o)
	}
}

/**
 * Client Version of BuildServices
 * @param  {[type]} u User)         BuildClient(s, b, env string, o bool [description]
 * @return {[type]}   [description]
 */
func (u User) BuildClient(s, b, env string, o bool) {
	service, err := FormatService(s, "client")
	if err != nil {
		log.Fatal(err)
		return
	}

	branch := FormatClientBranch(b)

	if _, err := EnsureBranchExists(service, branch); err != nil {
		log.Fatal("Requested branch does not exist on remote")
		return
	}

	var data = jenkins.Auth{
		Username: u.Username,
		ApiToken: u.ApiToken,
	}

	fmt.Printf("Attempting to build %s on branch %s", service, branch)
	RunJob(&data, service, branch, "client", env, o)
}

/**
 * Takes a string and if its empty defualts to master
 * @param {[type]} b string) (c string [description]
 */
func FormatBranch(b string) (c string) {
	if b == "" {
		b = "master"
	}
	return b
}

/**
 * If an empty branch is specified, then default it
 * to feature-deathstar-fully-operational, otherwise return
 * the specified branch
 * @param {[type]} b string) (c string [description]
 */
func FormatClientBranch(b string) (c string) {
	if b == "" {
		b = AdminV2DefaultBranch
	}
	return b
}

/**
 * Takes a string and if its empty returns an error. If its not
 * empty, it checks a map to make sure the service exists and returns
 * the formatted service job name if its there
 * @param {[type]} s string) (c string, e error [description]
 */
func FormatService(s string, t string) (c string, e error) {
	if s == "" {
		err := errors.New("Specify service to build")
		return "", err
	}
	switch t {
	case "services":
		s = JenkinsServices[s]
	case "client":
		s = JenkinsClients[s]
	}
	if s == "" {
		err := errors.New("Specified service does not exist")
		return "", err
	}

	return s, nil
}

/**
 * Makes a call to the remote repo to check if the branch exists
 * @param {[type]} s Service name
 * @param {[type]} b string)       (t bool, e error [description]
 */
func EnsureBranchExists(s, b string) (t bool, e error) {
	cmdArgs := formatCmdArgs(s, b)
	fmt.Println("-----> Running cmd to check if branch exists")
	fmt.Println(cmdArgs)

	_, err := exec.Command("/bin/sh", "-c", cmdArgs).Output()

	if err != nil {
		return false, err
	}

	return true, nil
}

func formatCmdArgs(s, b string) (c string) {
	return fmt.Sprintf("git ls-remote git@github.com:1stdibs/%s.git | grep -sw '%s'", s, b)
}
