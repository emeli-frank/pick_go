package seeder

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/emeli-frank/internships/backend/pkg/database"
	"io/ioutil"
	"math/rand"
	"net/http"
	"syreclabs.com/go/faker"

	"github.com/emeli-frank/internships/backend/pkg/domain/organization"
	//"github.com/emeli-frank/internships/backend/pkg/domain/offer"
	"github.com/emeli-frank/internships/backend/pkg/domain/users"
	"github.com/emeli-frank/internships/backend/pkg/domain/users/admin"
	//orgUser "github.com/emeli-frank/internships/backend/pkg/domain/users/organization"
	"log"
)

type Seeder struct {
	nextUserId func() int
	DB *sql.DB
	errorLog *log.Logger
	request *request
}

func NewSeeder(
	db *sql.DB,
	errorLog *log.Logger,
	/*organizationService *orgPkg.Service,
	offerService *offer.Service,
	organizationUserService *organization.Service*/) *Seeder {

	return &Seeder{
		nextUserId: nextUserIdGen(),
		DB: db,
		errorLog: errorLog,
		/*organizationService: organizationService,
		offerService: offerService,
		organizationUserService: organizationUserService,*/
		request: NewRequest(),
	}
}

func (s *Seeder) Mock() error {
	var err error
	err = database.ExecScripts(
		s.DB,
		"./pkg/storage/mysql/.db_setup/teardown.sql",
		"./pkg/storage/mysql/.db_setup/tables.sql",
		"./pkg/storage/mysql/.db_setup/data.sql")
	if err != nil {
		return err
	}

	s.userService.SetUser(&users.User{})
	s.userService.SetUser(&users.User{ID: 1})

	// create admin
	err = createAdmin(s.request)
	if err != nil {
		return err
	}

	// create orgPkg
	organizationIds, err := createOrganization(1, s.request)
	if err != nil {
		return err
	}

	_ = organizationIds

	// create orgPkg staff admin
	//organizationUserIds := 3
	/*err = createOrganizationStaff(s.organizationUserService, organizationIds, s.nextUserId)
	if err != nil {
		return err
	}*/

	/*// create offer
	_, err = createOffers(s.offerService, s.userService, organizationIds, organizationUserId)
	if err != nil {
		return err
	}*/

	// create applicant

	// create application

	return nil
}

func createAdmin(requestService *request) error {
	//u := &users.User{ID: nextId()}
	u := &users.User{Email: "email@example.com"}
	code, header, body, err := requestService.Post("admins", u)
	if err != nil {
		return err
	}
	fmt.Println("code:", code)
	fmt.Println("header:", header)
	fmt.Println("body:", string(body))

	type response struct {
		ID int `json:"id"`
	}
	resp := new(response)
	if err := json.Unmarshal(body, resp); err != nil {
		return nil
	}

	return nil
}

func createOrganization(count int, requestService *request) ([]int, error) {
	var organizationIds []int
	for i := 0; i < count; i++ {
		o := organization.Organization{
			Name:        faker.Company().Name(),
			Description: faker.Lorem().Sentence(20),
			State:       organization.State{
				ID:   (rand.Int() % 36) + 1,
			},
			Address:     faker.Address().StreetAddress(),
			Email:       faker.Internet().Email(),
			Phone:       faker.PhoneNumber().PhoneNumber(),
		}

		code, header, body, err := requestService.Post("organizations", o)
		if err != nil {
			return nil, err
		}
		fmt.Println("code:", code)
		fmt.Println("header:", header)
		fmt.Println("body:", string(body))
		type response struct {
			ID int `json:"id"`
		}
		resp := new(response)
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, nil
		}

		organizationIds = append(organizationIds, resp.ID)
	}

	return organizationIds, nil
}

/*func createOrganizationStaff(organizationUserService organization.Service, orgIds[]int, nextId func() int) error {
	for _, orgId := range orgIds {
		orgUser := orgUser.User{User: users.User{ID: nextId()}}
		err := organizationUserService.CreateAdmin(&orgUser, orgId)
		if err != nil {
			return err
		}

		roleIds := []int{users.RoleOrganizationAdmin,
			users.RoleOfferManager,
			users.RoleScreeningManager,
			users.RoleScoringManager,
			users.RoleFinalizingOfficer}

		for _, r := range roleIds {
			orgUser := orgUser.User{User: users.User{ID: nextId()}}
			err := organizationUserService.CreateStaff(&orgUser, orgId, []int{r})
			if err != nil {
				return err
			}
		}
	}

	return nil
}*/

// TODO:: REMOVE organizationUserId AND GET ID FROM AUTH
/*func createOffers(offerService *offer.Service, userService *users.Service, organizationId int, organizationUserId int) (int, error){
	id, err := offerService.Create(
		userService, strconv.Itoa(organizationId), "0", "1",
		"This is an internship", "1", []string{"1", "2", "3", "4", "5"},
		"", "1580515200", "1583020800", organizationUserId,
	)
	if err != nil {
		if validationErr, ok := err.(validation.Error); ok {
			fmt.Println(validationErr.ValidationErrors())
		}
		return 0, err
	}

	return int(id), nil
}*/

func nextUserIdGen() func() int {
	id := 1
	return func() int {
		id = id + 1
		return id
	}
}

type request struct {
	BaseUrl string
	ContentType string
}

func NewRequest() *request {
	return &request{
		BaseUrl: "http://localhost:5000/api/",
		ContentType: "application/json",
	}
}

func (r request) Post(path string, data interface{}) (int, http.Header, []byte, error) {
	url := r.BaseUrl + path
	dataJson, err := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		return 0, nil, nil, err
	}
	req.Header.Set("Content-Type", r.ContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	header := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)

	return code, header, body, nil
}
