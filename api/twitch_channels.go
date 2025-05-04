package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Channel struct {
	ID                          string   `json:"broadcaster_id"`
	Login                       string   `json:"broadcaster_login"`
	DisplayName                 string   `json:"broadcaster_name"`
	GameID                      string   `json:"game_id"`
	GameName                    string   `json:"game_name"`
	Title                       string   `json:"title"`
	Delay                       int      `json:"delay"`
	Tags                        []string `json:"tags"`
	ContentClassificationLabels []string `json:"content_classification_labels"`
	IsBrandedContent            bool     `json:"is_branded_content"`
}

type Followed struct {
	BroadcasterID    string    `json:"broadcaster_id"`
	BroadcasterLogin string    `json:"broadcaster_login"`
	BroadcasterName  string    `json:"broadcaster_name"`
	FollowedAt       time.Time `json:"followed_at"`
}

type Follower struct {
	UserID     string    `json:"user_id"`
	UserLogin  string    `json:"user_login"`
	UserName   string    `json:"user_name"`
	FollowedAt time.Time `json:"followed_at"`
}

type ChannelsResource struct {
	client *Client

	Followed  *FollowedResource
	Followers *FollowersResource
}

func NewChannelsResource(client *Client) *ChannelsResource {
	r := &ChannelsResource{client: client}

	r.Followed = NewFollowedResource(client)
	r.Followers = NewFollowersResource(client)

	return r
}

type ChannelsListCall struct {
	resource *ChannelsResource
	opts     []RequestOption
}

type ChannelsListResponse struct {
	Header http.Header
	Data   []Channel
}

// List creates a request to list channels based on the specified criteria.
func (r *ChannelsResource) List() *ChannelsListCall {
	return &ChannelsListCall{resource: r}
}

// BroadcasterID filters the results to the specified broadcaster ID.
func (c *ChannelsListCall) BroadcasterID(ids []string) *ChannelsListCall {
	for _, id := range ids {
		c.opts = append(c.opts, AddQueryParameter("broadcaster_id", id))
	}
	return c
}

// Do executes the request.
func (c *ChannelsListCall) Do(ctx context.Context, opts ...RequestOption) (*ChannelsListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channels", nil, append(opts, c.opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Channel](res)
	if err != nil {
		return nil, err
	}

	return &ChannelsListResponse{
		Header: res.Header,
		Data:   data.Data,
	}, nil
}

type FollowedResource struct {
	client *Client
}

type FollowedListCall struct {
	resource *FollowedResource
	opts     []RequestOption
}

type FollowedListResponse struct {
	Header http.Header
	Data   []Followed
	Cursor string
}

func NewFollowedResource(client *Client) *FollowedResource {
	return &FollowedResource{client}
}

func (r *FollowedResource) List() *FollowedListCall {
	return &FollowedListCall{resource: r}
}

// Use this parameter to see whether the user follows this broadcaster.
// If specified, the response contains this broadcaster if the user follows them.
// If not specified, the response contains all broadcasters that the user follows.
func (c *FollowedListCall) BroadcasterID(id string) *FollowedListCall {
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", id))
	return c
}

// UserID returns the list of broadcasters that this user follows
func (c *FollowedListCall) UserID(id string) *FollowedListCall {
	c.opts = append(c.opts, SetQueryParameter("user_id", id))
	return c
}

// The minimum page size is 1 item per page and the maximum is 100. The default is 20.
func (c *FollowedListCall) First(n int) *FollowedListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// The cursor used to get the next page of results.
func (c *FollowedListCall) After(n int) *FollowedListCall {
	c.opts = append(c.opts, SetQueryParameter("after", fmt.Sprint(n)))
	return c
}

func (c *FollowedListCall) Do(ctx context.Context, opts ...RequestOption) (*FollowedListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channels/followed", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[Followed](res)
	if err != nil {
		return nil, err
	}

	return &FollowedListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}

type FollowersResource struct {
	client *Client
}

type FollowersListCall struct {
	resource *FollowersResource
	opts     []RequestOption
}

type FollowersListResponse struct {
	Header http.Header
	Data   []Follower
	Cursor string
}

func NewFollowersResource(client *Client) *FollowersResource {
	return &FollowersResource{client}
}

func (r *FollowersResource) List() *FollowersListCall {
	return &FollowersListCall{resource: r}
}

// BroadcasterID filters the results to the specified broadcaster ID.
func (c *FollowersListCall) BroadcasterID(id string) *FollowersListCall {
	c.opts = append(c.opts, SetQueryParameter("broadcaster_id", id))
	return c
}

// UserID filters the results to the specified user ID.
func (c *FollowersListCall) UserID(id string) *FollowersListCall {
	c.opts = append(c.opts, SetQueryParameter("user_id", id))
	return c
}

// The minimum page size is 1 item per page and the maximum is 100. The default is 20.
func (c *FollowersListCall) First(n int) *FollowersListCall {
	c.opts = append(c.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return c
}

// The cursor used to get the next page of results.
func (c *FollowersListCall) After(n int) *FollowersListCall {
	c.opts = append(c.opts, SetQueryParameter("after", fmt.Sprint(n)))
	return c
}
func (c *FollowersListCall) Do(ctx context.Context, opts ...RequestOption) (*FollowersListResponse, error) {
	res, err := c.resource.client.doRequest(ctx, http.MethodGet, "/channels/followers", nil, append(c.opts, opts...)...)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := decodeResponse[Follower](res)
	if err != nil {
		return nil, err
	}

	return &FollowersListResponse{
		Header: res.Header,
		Data:   data.Data,
		Cursor: data.Pagination.Cursor,
	}, nil
}
