package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	o "reserve_restaurant/order"
	r "reserve_restaurant/restaurant"
	u "reserve_restaurant/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// NewUser godoc
// @Summary new user register
// @Description new user register
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param name body string true "name"
// @Param password body string true "password"
// @Success 200
// @Router /register [post]
func (handler *Handler) NewUser(c *gin.Context) {
	ctx := context.Background()

	user := &u.User{}
	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := handler.service.NewUser(ctx, user.Name, user.Password)
	if err != nil {
		if errors.Is(err, u.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": id,
	})
}

// @Summary user login
// @Description user login and return token
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param name body string true "name"
// @Param password body string true "password"
// @Success 200 string string "token and refresh token"
// @Router /login [POST]
func (handler *Handler) UserLogin(c *gin.Context) {
	ctx := context.Background()

	user := &u.User{}
	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, refreshToken, err := handler.service.UserLogin(ctx, user.Name, user.Password)
	if err != nil {
		if errors.Is(err, u.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	m := map[string]string{
		"token_type":    "bearer",
		"access_token":  token,
		"refresh_token": refreshToken,
	}
	b, _ := json.MarshalIndent(m, "", " ")
	_, err = c.Writer.Write(b)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary get restaurant'e menu
// @Description provide restaurant id to get it's menu
// @Tags restaurant
// @Accept application/json
// @Produce application/json
// @Param rid path int true "restaurant's id"
// @Success 200 string restaurant's menu
// @Router /restaurant/{rid} [GET]
func (handler *Handler) GetRestaurantMenu(c *gin.Context) {
	ctx := context.Background()

	rid, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	menu, err := handler.service.GetRestaurantMenu(ctx, rid)
	if err != nil {
		if errors.Is(err, r.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	menuJson, err := json.MarshalIndent(menu, "", "\t")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"menu": string(menuJson),
	})
}

// @Summary search restaurant
// @Description provide address to find nearby restaurant
// @Tags restaurant
// @Accept application/json
// @Produce applicatoin/json
// @Param location path string true "address"
// @Success 200
// @Router /restaurant/{location} [GET]
func (handler *Handler) GetNearbyRestaurant(c *gin.Context) {
	ctx := context.Background()

	location := c.Param("location")
	restaurants, err := handler.service.SearchRestaurant(ctx, location)
	if err != nil {
		if errors.Is(err, r.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	restaurantsJson, err := json.MarshalIndent(restaurants, "", " ")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"location":    location,
		"restaurants": string(restaurantsJson),
	})
}

// @Summary create food
// @Description restaurant insert new food
// @Tags restaurant
// @Accept application/json
// @Produce application/json
// @Param rid path int true "restaurant id"
// @Param name body string true "food's name"
// @Param description body string true "food's description"
// @Param price body float32 true "food's price"
// @Success 200
// @Router /restaurant/{rid}/food [POST]
func (handler *Handler) CreateFood(c *gin.Context) {
	ctx := context.Background()

	food := &r.Food{}
	if err := c.BindJSON(food); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := handler.service.CreateFood(ctx, food); err != nil {
		if errors.Is(err, r.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary create restaurant
// @Description register for restaurant
// @Tags restaurant
// @Accept application/json
// @Produce application/json
// @Param name body string true "restaurant name"
// @Param description body string true "restaurant description"
// @Param Location body string true "restaurant location"
// @Param Latitude body float32 true "restaurant latitude"
// @Param Longtitude body float32 true "restaurant longtitude"
// @Success 200
// @Router /restaurant [POST]
func (handler *Handler) CreateRestaurant(c *gin.Context) {
	ctx := context.Background()

	restaurant := &r.Restaurant{}
	if err := c.BindJSON(restaurant); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if err := handler.service.CreateRestaurant(ctx, restaurant); err != nil {
		if errors.Is(err, r.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary create order
// @Description user create order
// @Tags order
// @Accept application/json
// @Produce application/json
// @Param order body Order true "user's order"
// @Success 200
// @Router /order [POST]
func (handler *Handler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()
	userIds := ctx.Value(contextType("userId")).(string)
	uid, _ := strconv.Atoi(userIds)

	order := &o.Order{
		Uid:       uid,
		CreatedAt: time.Now(),
	}
	if err := c.BindJSON(order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	oid, err := handler.service.CreateOrder(ctx, order)
	if err != nil {
		if errors.Is(err, o.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"order_id": oid,
	})
}

// @Summary get order
// @Description fetch user's order
// @Tags order
// @Accept application/json
// @Produce application/json
// @Param oid path int true "order id"
// @Success 200
// @Router /order/{oid} [GET]
func (handler *Handler) GetOrder(c *gin.Context) {
	ctx := context.Background()
	userIds := ctx.Value(contextType("userId")).(string)
	uid, _ := strconv.Atoi(userIds)

	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	order, err := handler.service.GetOrder(ctx, oid)
	if err != nil {
		if errors.Is(err, o.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	if order.Uid != uid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	orderJson, _ := json.MarshalIndent(order, "", " ")
	c.JSON(http.StatusOK, map[string]interface{}{
		"order": string(orderJson),
	})
}

// @Summary get orders
// @Description get all of order from user
// @Tags order
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /user/order [GET]
func (handler *Handler) GetOrderForUser(c *gin.Context) {
	ctx := context.Background()
	userIds := ctx.Value(contextType("userId")).(string)
	uid, _ := strconv.Atoi(userIds)

	orders, err := handler.service.GetOrderForUser(ctx, uid)
	if err != nil {
		if errors.Is(err, o.ErrInternalServer) {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	ordersJson, _ := json.MarshalIndent(orders, "", " ")
	c.JSON(http.StatusOK, map[string]interface{}{
		"order": string(ordersJson),
	})
}
