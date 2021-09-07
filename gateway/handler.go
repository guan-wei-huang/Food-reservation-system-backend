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
