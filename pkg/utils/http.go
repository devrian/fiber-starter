package utils

import (
	"bytes"
	"encoding/json"
	"fiber-starter/pkg/common"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Meta   Meta        `json:"meta"`
	Result interface{} `json:"result"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination Paginate    `json:"pagination"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Paginate struct {
	Total  int64               `json:"total"`
	Paging common.PaginateLink `json:"paging"`
}

const (
	ChannelApp = "app"
	ChannelWeb = "web"
)

func WithPagination(data interface{}, total int64, paging common.PaginateLink) interface{} {
	return PaginationResponse{
		Data: data,
		Pagination: Paginate{
			Total:  total,
			Paging: paging,
		},
	}
}

func APIResponse(c *fiber.Ctx, message string, code int, status string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  status,
		},
		Result: data,
	})
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func APIResponseErrorByValidationError(c *fiber.Ctx, err error) error {
	errors := FormatValidationError(err)
	errorMessage := fiber.Map{"errors": errors}
	response := APIResponse(c, fmt.Sprintf(`%s failed.`, c.Path()), fiber.StatusUnprocessableEntity, fiber.ErrUnprocessableEntity.Message, errorMessage)

	return response
}

func EnsureDir(dirName string, mode os.FileMode) error {
	err := os.Mkdir(dirName, mode)
	if err == nil || os.IsExist(err) {
		return nil
	}

	return err
}

func RequestHandler(bodyRequest map[string]interface{}, url, method string) (*http.Request, error) {
	dataValues, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	reqBody := []byte(string(dataValues))
	request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return request, err
	}

	return request, nil
}

func ResponseHandler(request *http.Request) (map[string]interface{}, error) {
	var result map[string]interface{}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func DoAsyncRequest(request *http.Request, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("err handle it")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("err handle it")
	}

	ch <- string(body)
}

func ResponseAsyncHandler(request *http.Request) ([]string, error) {
	// make a channel
	ch := make(chan string)
	var wg sync.WaitGroup

	// do async request with channel
	wg.Add(1)
	go DoAsyncRequest(request, ch, &wg)

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	// read from channel as they come in until its closed
	var responses []string
	for res := range ch {
		responses = append(responses, res)
	}

	return responses, nil
}

func RequestHandlerEntity(entity interface{}, url, method string) (*http.Request, error) {
	dataValues, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	reqBody := []byte(string(dataValues))
	request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return request, err
	}

	return request, nil
}
