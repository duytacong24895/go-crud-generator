package pkg

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Initialize Resty client with common settings
func NewHTTPClient(baseURl string) *resty.Client {
	client := resty.New()

	// Set common settings
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)

	// Set common headers
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "MyApp/1.0",
	})

	// Optional: Set base URL
	client.SetBaseURL(baseURl)

	return client
}

// GET request function
func GetRequest(client *resty.Client, url string, result interface{}) (serverAlive bool, err error) {
	resp, err := client.R().
		SetResult(result). // Automatically unmarshal success response
		Get(url)

	if err != nil {
		return false, fmt.Errorf("GET request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("GET request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// GET with query parameters
func GetWithParams(client *resty.Client, url string, params map[string]string, result interface{}) (serverAlive bool, err error) {
	resp, err := client.R().
		SetResult(result).
		SetQueryParams(params). // Set query parameters
		Get(url)

	if err != nil {
		return false, fmt.Errorf("GET request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("GET request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// POST request function
func PostRequest(client *resty.Client, url string, body interface{}, result interface{}) (
	serverAlive bool, err error) {
	resp, err := client.R().
		SetBody(body).
		SetResult(result).
		Post(url)

	if err != nil {
		return false, fmt.Errorf("POST request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("POST request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// PUT request function (for full updates)
func PutRequest(client *resty.Client, url string, body interface{}, result interface{}) (serverAlive bool, err error) {
	resp, err := client.R().
		SetBody(body).
		SetResult(result).
		Put(url)

	if err != nil {
		return false, fmt.Errorf("PUT request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("PUT request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// PATCH request function (for partial updates)
func PatchRequest(client *resty.Client, url string, body interface{}, result interface{}) (serverAlive bool, err error) {
	resp, err := client.R().
		SetBody(body).
		SetResult(result).
		Patch(url)

	if err != nil {
		return false, fmt.Errorf("PATCH request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("PATCH request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// DELETE request function
func DeleteRequest(client *resty.Client, url string) (serverAlive bool, err error) {
	resp, err := client.R().Delete(url)

	if err != nil {
		return false, fmt.Errorf("DELETE request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("DELETE request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return true, nil
}

// Advanced function with headers and auth
func RequestWithAuth(client *resty.Client, method, url, token string, body interface{}, result interface{}) (serverAlive bool, err error) {
	req := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	var resp *resty.Response

	switch method {
	case "GET":
		resp, err = req.Get(url)
	case "POST":
		resp, err = req.Post(url)
	case "PUT":
		resp, err = req.Put(url)
	case "PATCH":
		resp, err = req.Patch(url)
	case "DELETE":
		resp, err = req.Delete(url)
	default:
		return false, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return false, fmt.Errorf("%s request failed: %w", method, err)
	}

	if resp.StatusCode() >= 400 {
		return true, fmt.Errorf("%s request failed with status %d: %s", method, resp.StatusCode(), resp.String())
	}

	return true, nil
}

// // Example usage
// func main() {
// 	client := NewHTTPClient()

// 	// GET example
// 	fmt.Println("=== GET Request ===")
// 	var user User
// 	err := GetRequest(client, "/users/1", &user)
// 	if err != nil {
// 		log.Printf("GET error: %v", err)
// 	} else {
// 		fmt.Printf("User: %+v\n", user)
// 	}

// 	// GET with parameters
// 	fmt.Println("\n=== GET with Parameters ===")
// 	var users []User
// 	params := map[string]string{
// 		"userId": "1",
// 	}
// 	err = GetWithParams(client, "/posts", params, &users)
// 	if err != nil {
// 		log.Printf("GET with params error: %v", err)
// 	}

// 	// POST example
// 	fmt.Println("\n=== POST Request ===")
// 	newUser := User{
// 		Name:  "John Doe",
// 		Email: "john@example.com",
// 		Age:   30,
// 	}
// 	var createdUser User
// 	err = PostRequest(client, "/users", newUser, &createdUser)
// 	if err != nil {
// 		log.Printf("POST error: %v", err)
// 	} else {
// 		fmt.Printf("Created user: %+v\n", createdUser)
// 	}

// 	// PUT example (full update)
// 	fmt.Println("\n=== PUT Request ===")
// 	updateUser := User{
// 		ID:    1,
// 		Name:  "Jane Doe",
// 		Email: "jane@example.com",
// 		Age:   25,
// 	}
// 	var updatedUser User
// 	err = PutRequest(client, "/users/1", updateUser, &updatedUser)
// 	if err != nil {
// 		log.Printf("PUT error: %v", err)
// 	} else {
// 		fmt.Printf("Updated user: %+v\n", updatedUser)
// 	}

// 	// PATCH example (partial update)
// 	fmt.Println("\n=== PATCH Request ===")
// 	partialUpdate := map[string]interface{}{
// 		"name": "Updated Name",
// 	}
// 	var patchedUser User
// 	err = PatchRequest(client, "/users/1", partialUpdate, &patchedUser)
// 	if err != nil {
// 		log.Printf("PATCH error: %v", err)
// 	} else {
// 		fmt.Printf("Patched user: %+v\n", patchedUser)
// 	}

// 	// DELETE example
// 	fmt.Println("\n=== DELETE Request ===")
// 	err = DeleteRequest(client, "/users/1")
// 	if err != nil {
// 		log.Printf("DELETE error: %v", err)
// 	} else {
// 		fmt.Println("User deleted successfully")
// 	}

// 	// Request with authentication
// 	fmt.Println("\n=== Authenticated Request ===")
// 	var authUser User
// 	err = RequestWithAuth(client, "GET", "/users/1", "your-jwt-token", nil, &authUser)
// 	if err != nil {
// 		log.Printf("Auth request error: %v", err)
// 	}
// }
