package post

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/errno"
	"{[.RootPackage]}/internal/pkg/known"
	"{[.RootPackage]}/internal/pkg/log"
	v1 "{[.RootPackage]}/pkg/api/{[.AppName]}/v1"
)

// List 返回博客列表.
//
// @Summary    List posts
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      request	 query	    v1.ListPostRequest	 true  "Param"
// @Success	   200		{object}	v1.ListPostResponse
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts [GET].
func (ctrl *PostController) List(c *gin.Context) {
	log.C(c).Infow("List post function called")

	var r v1.ListPostRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	resp, err := ctrl.b.Posts().List(c, c.GetString(known.XUsernameKey), r.Offset, r.Limit)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Create 创建一条博客.
//
// @Summary    Create a post
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      request	 body	    v1.CreatePostRequest	 true  "Param"
// @Success	   200		{object}	v1.GetPostResponse
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts [POST].
func (ctrl *PostController) Create(c *gin.Context) {
	log.C(c).Infow("Create post function called")

	var r v1.CreatePostRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	resp, err := ctrl.b.Posts().Create(c, c.GetString(known.XUsernameKey), &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Get 获取指定的博客.
//
// @Summary    Get post info
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      postID	  path	    string          	 true  "postID"
// @Success	   200		{object}	v1.ListPostResponse
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts/{postID} [GET].
func (ctrl *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called")

	post, err := ctrl.b.Posts().Get(c, c.GetString(known.XUsernameKey), c.Param("postID"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, post)
}

// Update 更新博客.
//
// @Summary    Update post info
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      postID	  path	    string          	  true  "postID"
// @Param      request	 query	    v1.UpdatePostRequest  true  "Param"
// @Success	   200		{object}	nil
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts/{postID} [PUT].
func (ctrl *PostController) Update(c *gin.Context) {
	log.C(c).Infow("Update post function called")

	var r v1.UpdatePostRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	err := ctrl.b.Posts().Update(c, c.GetString(known.XUsernameKey), c.Param("postID"), &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}

// Delete 删除指定的博客.
//
// @Summary    Delete a post
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      postID	  path	    string          	  true  "postID"
// @Success	   200		{object}	nil
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts/{postID} [DELETE].
func (ctrl *PostController) Delete(c *gin.Context) {
	log.C(c).Infow("Delete post function called")

	if err := ctrl.b.Posts().Delete(c, c.GetString(known.XUsernameKey), c.Param("postID")); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}

// DeleteCollection 批量删除博客.
//
// @Summary    Batch delete posts
// @Security   Bearer
// @Tags       Post
// @Accept     application/json
// @Produce    json
// @Param      postID	  path	    array          	  true  "postID"
// @Success	   200		{object}	nil
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/posts/{postID} [DELETE].
func (ctrl *PostController) DeleteCollection(c *gin.Context) {
	log.C(c).Infow("Batch delete post function called")

	postIDs := c.QueryArray("postID")
	if err := ctrl.b.Posts().DeleteCollection(c, c.GetString(known.XUsernameKey), postIDs); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
