package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

type Error struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type WallPost struct {
	ID          int          `json:"id"`
	OwnerID     int          `json:"owner_id"`
	Date        int          `json:"date"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type  string          `json:"type"`
	Photo PhotoAttachment `json:"photo"`
	Video VideoAttachment `json:"video"`
}

type PhotoAttachment struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	AccessKey string `json:"access_key"`
}

type VideoAttachment struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Date    int    `json:"date"`
	Title   string `json:"title"`
}

func GetWallPosts(values url.Values) ([]WallPost, error) {
	u := composeURL("wall.get")

	rawData, err := sendRequest(u, values)
	if err != nil {
		return []WallPost{}, err
	}

	wallPosts, err := parseWallPosts(rawData)
	if err != nil {
		return []WallPost{}, err
	}

	return wallPosts, nil
}

type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ScreenName string `json:"screen_name"`
}

func GetUserInfo(values url.Values) (User, error) {
	u := composeURL("users.get")

	rawData, err := sendRequest(u, values)
	if err != nil {
		return User{}, err
	}

	user, err := parseUserInfo(rawData)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type Community struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Description string `json:"description"`
}

func GetCommunityInfo(values url.Values) (Community, error) {
	u := composeURL("groups.getById")

	rawData, err := sendRequest(u, values)
	if err != nil {
		return Community{}, err
	}

	community, err := parseCommunityInfo(rawData)
	if err != nil {
		return Community{}, err
	}

	return community, nil
}

func composeURL(method string) string {
	u := fmt.Sprintf("https://api.vk.com/method/%s", method)

	return u
}

func parseWallPosts(rawData []byte) ([]WallPost, error) {
	var data struct {
		Response struct {
			Count int        `json:"count"`
			Items []WallPost `json:"items"`
			Error Error      `json:"error"`
		} `json:"response"`
	}

	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return []WallPost{}, err
	}

	if (data.Response.Error != Error{}) {
		return []WallPost{}, fmt.Errorf("error %d: %s",
			data.Response.Error.ErrorCode, data.Response.Error.ErrorMsg)
	}

	return data.Response.Items, nil
}

func parseUserInfo(rawData []byte) (User, error) {
	var data struct {
		Response []User `json:"response"`
		Error    Error  `json:"error"`
	}

	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return User{}, err
	}

	if (data.Error != Error{}) {
		return User{}, fmt.Errorf("error %d: %s", data.Error.ErrorCode, data.Error.ErrorMsg)
	}

	return data.Response[0], nil
}

func parseCommunityInfo(rawData []byte) (Community, error) {
	var data struct {
		Response []Community `json:"response"`
		Error    Error       `json:"error"`
	}

	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return Community{}, err
	}

	if (data.Error != Error{}) {
		return Community{}, fmt.Errorf("error %d: %s", data.Error.ErrorCode, data.Error.ErrorMsg)
	}

	return data.Response[0], nil
}

func sendRequest(u string, values url.Values) ([]byte, error) {
	response, err := http.PostForm(u, values)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("%s\n\n%s\n", err, debug.Stack())
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
