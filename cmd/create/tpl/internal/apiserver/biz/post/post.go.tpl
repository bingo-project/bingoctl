package post

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"{[.RootPackage]}/internal/pkg/log"

	"github.com/jinzhu/copier"

	"{[.RootPackage]}/internal/apiserver/store"
	"{[.RootPackage]}/internal/pkg/errno"
	"{[.RootPackage]}/internal/pkg/model"
	v1 "{[.RootPackage]}/pkg/api/{[.AppName]}/v1"
)

// PostBiz 定义了 user 模块在 biz 层所实现的方法.
type PostBiz interface {
	Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error)
	Get(ctx context.Context, username, postID string) (*v1.GetPostResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListPostResponse, error)
	Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error
	Delete(ctx context.Context, username, postID string) error
	DeleteCollection(ctx context.Context, username string, postIDs []string) error
}

// The implementation of PostBiz interface.
type postBiz struct {
	ds store.IStore
}

// Make sure that postBiz implements the PostBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ PostBiz = (*postBiz)(nil)

func New(ds store.IStore) *postBiz {
	return &postBiz{ds: ds}
}

// Create is the implementation of the `Create` method in PostBiz interface.
func (b *postBiz) Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, r)
	postM.Username = username

	if err := b.ds.Posts().Create(ctx, &postM); err != nil {
		return nil, err
	}

	return &v1.CreatePostResponse{PostID: postM.PostID}, nil
}

// Get is the implementation of the `Get` method in PostBiz interface.
func (b *postBiz) Get(ctx context.Context, username, postID string) (*v1.GetPostResponse, error) {
	user, err := b.ds.Posts().Get(ctx, username, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.GetPostResponse
	_ = copier.Copy(&resp, user)

	return &resp, nil
}

// List is the implementation of the `List` method in PostBiz interface.
func (b *postBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListPostResponse, error) {
	count, list, err := b.ds.Posts().List(ctx, username, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list posts from storage", "err", err)

		return nil, err
	}

	posts := make([]*v1.PostInfo, 0, len(list))
	for _, item := range list {
		var post v1.PostInfo
		_ = copier.Copy(&post, item)

		posts = append(posts, &post)
	}

	return &v1.ListPostResponse{TotalCount: count, Data: posts}, nil
}

// Update is the implementation of the `Update` method in PostBiz interface.
func (b *postBiz) Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error {
	post, err := b.ds.Posts().Get(ctx, username, postID)
	if err != nil {
		return err
	}

	if r.Title != nil {
		post.Title = *r.Title
	}

	if r.Content != nil {
		post.Content = *r.Content
	}

	if err := b.ds.Posts().Update(ctx, post); err != nil {
		return err
	}

	return nil
}

// Delete is the implementation of the `Delete` method in PostBiz interface.
func (b *postBiz) Delete(ctx context.Context, username, postID string) error {
	return b.ds.Posts().Delete(ctx, username, []string{postID})
}

// DeleteCollection is the implementation of the `DeleteCollection` method in PostBiz interface.
func (b *postBiz) DeleteCollection(ctx context.Context, username string, postIDs []string) error {
	if err := b.ds.Posts().Delete(ctx, username, postIDs); err != nil {
		return err
	}

	return nil
}
