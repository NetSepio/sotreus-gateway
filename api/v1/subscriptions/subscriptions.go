package subscriptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/NetSepio/sotreus-gateway/api/middleware/auth/paseto"
	"github.com/NetSepio/sotreus-gateway/config/dbconfig"
	"github.com/NetSepio/sotreus-gateway/config/envconfig"
	"github.com/NetSepio/sotreus-gateway/models"
	"github.com/NetSepio/sotreus-gateway/util/pkg/logwrapper"
	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
	"gorm.io/gorm"
)

func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/subscriptions")
	{
		g.POST("webhook", HandleWebhook)
		g.Use(paseto.PASETO(false))
		g.POST("/create-payment", CreatePaymentIntent)
	}
}

func CreatePaymentIntent(c *gin.Context) {
	userId := c.GetString(paseto.CTX_USER_ID)
	db := dbconfig.GetDb()
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1099),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		logwrapper.Errorf("failed to create new payment intent: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}

	// insert in above table
	err = db.Create(&models.UserStripePi{
		Id:           uuid.NewString(),
		UserId:       userId,
		StripePiId:   pi.ID,
		StripePiType: models.SotreusSubscription,
	}).Error
	if err != nil {
		logwrapper.Errorf("failed to insert into users_stripe_pi: %v", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "internal server error").SendD(c)
		return
	}

	httpo.NewSuccessResponseP(200, "Created new charge", gin.H{"clientSecret": pi.ClientSecret}).SendD(c)
}

func HandleWebhook(c *gin.Context) {
	db := dbconfig.GetDb()
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logwrapper.Errorf("failed to create new payment intent: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}

	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), envconfig.EnvVars.STRIPE_WEBHOOK_SECRET)

	if err != nil {
		logwrapper.Errorf("failed to create new payment intent: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}

	switch event.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		var userStripePi models.UserStripePi
		if err := db.Where("stripe_pi_id = ?", paymentIntent.ID).First(&userStripePi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//warn and return success
				logwrapper.Warnf("No user found with stripe_pi_id: %v", err)
				c.JSON(http.StatusOK, gin.H{"status": "received"})
				return
			}
			logwrapper.Errorf("Error getting user with stripe_pi_id: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		var user models.User
		if err := db.Where("user_id = ?", userStripePi.UserId).First(&user).Error; err != nil {
			logwrapper.Errorf("Error getting user with user_id: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		subscription := models.Subscription{
			UserId:    user.UserId,
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}

		if err = db.Model(models.Subscription{}).Create(&subscription).Error; err != nil {
			logwrapper.Errorf("Error creating subscription: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		fmt.Println("minting nft -- 111NFT")

	case stripe.EventTypePaymentIntentCanceled:
		err := HandleCanceledOrFailedPaymentIntent(event.Data.Raw)
		if err != nil {
			logwrapper.Errorf("Error handling canceled payment intent: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "recieved"})
}

func HandleCanceledOrFailedPaymentIntent(eventDataRaw json.RawMessage) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(eventDataRaw, &paymentIntent)
	if err != nil {
		return fmt.Errorf("error parsing webhook JSON: %w", err)
	}

	return nil
}
