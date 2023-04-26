package v1

import "time"

// CreatePostRequest 指定了 `POST /v1/posts` 接口的请求参数.
type CreatePostRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CreatePostResponse 指定了 `POST /v1/posts` 接口的请求参数.
type CreatePostResponse struct {
	PostID string `json:"postID"`
}

// GetPostResponse 指定了 `GET /v1/posts/{postID}` 接口的返回参数.
type GetPostResponse PostInfo

// PostInfo 指定了用户的详细信息.
type PostInfo struct {
	Username  string    `json:"username,omitempty"`
	PostID    string    `json:"postID,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListPostRequest 指定了 `GET /v1/posts` 接口的请求参数.
type ListPostRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// ListPostResponse 指定了 `GET /v1/posts` 接口的返回参数.
type ListPostResponse struct {
	TotalCount int64       `json:"totalCount"`
	Posts      []*PostInfo `json:"posts"`
}

// UpdatePostRequest 指定了 `PUT /v1/posts` 接口的请求参数.
type UpdatePostRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}
