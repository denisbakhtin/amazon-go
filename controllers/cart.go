package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/denisbakhtin/amazon-go/aws"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//CartGet processes GET /cart/
func CartGet(c *gin.Context) {
	session := sessions.Default(c)

	apiRes := &aws.CartResponse{}
	if session.Get("cart_id") != nil && session.Get("cart_hmac") != nil && session.Get("cart_json") != nil {
		if err := json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes); err != nil {
			ShowErrorPage(c, 500, err)
			return
		}
	}

	H := DefaultH(c)
	H["Title"] = "My cart"
	H["BackUrl"] = "/"
	H["Cart"] = apiRes
	H["Flash"] = session.Flashes()
	session.Save()
	c.HTML(200, "cart/show", H)
}

//CartJSONGet processes GET /cart/get
func CartJSONGet(c *gin.Context) {
	reply := aws.CartAddReply{}
	session := sessions.Default(c)

	apiRes := &aws.CartResponse{}
	if session.Get("cart_id") != nil && session.Get("cart_hmac") != nil && session.Get("cart_json") != nil {
		if err := json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes); err != nil {
			c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
			return
		}
		reply.Quantity = apiRes.TotalCount()
	}

	c.JSON(200, reply)
}

//CartAdd processes POST /cart/add/:id
func CartAdd(c *gin.Context) {
	id := c.Param("id")
	variation := models.Variation{}
	session := sessions.Default(c)

	models.DB.First(&variation, id)
	if variation.ID == 0 {
		ShowErrorPage(c, 400, nil)
		return
	}
	apiRes := &aws.CartResponse{}
	var err error
	if session.Get("cart_id") == nil || session.Get("cart_hmac") == nil || session.Get("cart_id") == "" || session.Get("cart_hmac") == "" || session.Get("cart_json") == nil || session.Get("cart_json") == "" {
		//cart create request
		apiRes, err = aws.CartCreate(variation.Asin)
		if err != nil {
			ShowErrorPage(c, 500, err)
			return
		}
	} else {
		//cart add or modify request
		err = json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes)
		if err != nil {
			ShowErrorPage(c, 500, err)
			return
		}
		if foundItem := apiRes.FindItem(variation.Asin); foundItem != nil {
			apiRes, err = aws.CartModify(session.Get("cart_id").(string), session.Get("cart_hmac").(string), variation.Asin, foundItem.Quantity+1, apiRes)
		} else {
			apiRes, err = aws.CartAdd(session.Get("cart_id").(string), session.Get("cart_hmac").(string), variation.Asin, 1, apiRes)
		}
		if err != nil {
			//get previous cart state
			err = json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes)
			if err != nil {
				ShowErrorPage(c, 500, err)
				return
			}
		}
	}

	s, err := json.Marshal(*apiRes)
	if err != nil {
		ShowErrorPage(c, 500, err)
		return
	}
	session.Set("cart_id", apiRes.CartID)
	session.Set("cart_hmac", apiRes.Hmac)
	session.Set("cart_json", string(s))
	session.Save()
	c.Redirect(http.StatusSeeOther, "/cart")
}

//CartAddJSONPost processes POST /cart/add
//obsolete
func CartAddJSONPost(c *gin.Context) {
	cartItem := aws.CartItemForm{}
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)

	apiRes := &aws.CartResponse{}
	var err error
	if session.Get("cart_id") == nil || session.Get("cart_hmac") == nil || session.Get("cart_id") == "" || session.Get("cart_hmac") == "" || session.Get("cart_json") == nil || session.Get("cart_json") == "" {
		//cart create request
		apiRes, err = aws.CartCreate(cartItem.Asin)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		//cart add or modify request
		err = json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if foundItem := apiRes.FindItem(cartItem.Asin); foundItem != nil {
			apiRes, err = aws.CartModify(session.Get("cart_id").(string), session.Get("cart_hmac").(string), cartItem.Asin, foundItem.Quantity+cartItem.Quantity, apiRes)
		} else {
			apiRes, err = aws.CartAdd(session.Get("cart_id").(string), session.Get("cart_hmac").(string), cartItem.Asin, cartItem.Quantity, apiRes)
		}
		if err != nil {
			//get previous cart state
			err = json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	s, err := json.Marshal(*apiRes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session.Set("cart_id", apiRes.CartID)
	session.Set("cart_hmac", apiRes.Hmac)
	session.Set("cart_json", string(s))
	session.Save()

	reply := aws.CartAddReply{Quantity: apiRes.TotalCount()}
	c.JSON(200, reply)
}

//CartUpdatePost processes POST /cart/update
func CartUpdatePost(c *gin.Context) {
	cartItem := aws.CartItemForm{}
	if err := c.ShouldBind(&cartItem); err != nil {
		c.HTML(400, "errors/400", gin.H{"Error": err.Error()})
		return
	}

	session := sessions.Default(c)
	if session.Get("cart_id") == nil || session.Get("cart_hmac") == nil || session.Get("cart_id") == "" || session.Get("cart_hmac") == "" || session.Get("cart_json") == nil || session.Get("cart_json") == "" {
		c.HTML(500, "errors/500", gin.H{"Error": "Cart has not been initialized"})
		return
	}
	apiRes := &aws.CartResponse{}
	var err error

	if err := json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes); err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	apiRes, err = aws.CartModify(session.Get("cart_id").(string), session.Get("cart_hmac").(string), cartItem.Asin, cartItem.Quantity, apiRes)
	if err != nil {
		session.AddFlash(err.Error())
	} else {
		s, err := json.Marshal(*apiRes)
		if err != nil {
			c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
			return
		}
		session.Set("cart_json", string(s))
	}
	session.Save()

	c.Redirect(http.StatusFound, "/cart")
}

//CartCheckoutGet processes GET /cart/checkout
func CartCheckoutGet(c *gin.Context) {

	session := sessions.Default(c)

	apiRes := &aws.CartResponse{}
	if session.Get("cart_id") != nil && session.Get("cart_hmac") != nil && session.Get("cart_json") != nil {

		if err := json.Unmarshal([]byte(session.Get("cart_json").(string)), apiRes); err != nil {
			c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
			return
		}
		//clear cart cache and redirect to amazon
		session.Delete("cart_id")
		session.Delete("cart_hmac")
		session.Delete("cart_json")
		session.Save()
		c.Redirect(http.StatusFound, apiRes.PurchaseURL)
	} else {
		c.Redirect(http.StatusFound, "/cart")
	}
}
