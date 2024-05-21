package controllers

import (
	"context"
	"fmt"
	"gin/database"
	helper "gin/helper"
	"gin/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users_data")
var blogCollection *mongo.Collection = database.OpenCollection(database.Client, "recent_Actions")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"codeforce_handle": user.Codeforce_handle})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the Codeforce handle"})

		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this codeforce handle or Email already exists"})
			return
		}
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.User_name, *user.Codeforce_handle, *&user.User_id)
		// user.Token = &token
		// user.Refresh_token = &refreshToken

		user.Subscribedblogs = &[]string{}

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		// c.SetCookie("token", token, 86400, "/", "localhost", false, true)
		// c.SetCookie("refreshtoken", refreshToken, 86400, "/", "localhost", false, true)
		fmt.Printf("resultInsertionNumber =%v", resultInsertionNumber)
		c.JSON(http.StatusOK, gin.H{"token": token, "refreshtoken": refreshToken})

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.User_name, *foundUser.Codeforce_handle, foundUser.User_id)
		// helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// c.SetCookie("token", token, 86400, "/", "localhost", false, true)
		fmt.Printf("%v", foundUser)
		c.JSON(http.StatusOK, gin.H{"token": token, "refreshtoken": refreshToken})
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}
// 		page, err1 := strconv.Atoi(c.Query("page"))
// 		if err1 != nil || page < 1 {
// 			page = 1
// 		}

// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err = strconv.Atoi(c.Query("startIndex"))

// 		matchStage := bson.D{{"$match", bson.D{{}}}}
// 		groupStage := bson.D{{"$group", bson.D{
// 			{"_id", bson.D{{"_id", "null"}}},
// 			{"total_count", bson.D{{"$sum", 1}}},
// 			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
// 		projectStage := bson.D{
// 			{"$project", bson.D{
// 				{"_id", 0},
// 				{"total_count", 1},
// 				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
// 		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage, groupStage, projectStage})
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
// 		}
// 		var allusers []bson.M
// 		if err = result.All(ctx, &allusers); err != nil {
// 			log.Fatal(err)
// 		}
// 		c.JSON(http.StatusOK, allusers[0])
// 	}
// }

// func GetUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userId := c.Param("user_id")

// 		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		var user models.User
// 		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, user)
// 	}
// }

// type Server struct {
// 	r     *gin.Engine
// 	store *store.MongoStore
// }

// func (srv *Server) RecentActionsHandler(ctx *gin.Context) {
// 	recentActions, err := srv.store.QueryRecentActions()
// 	if err != nil {
// 		log.Printf("Error occurred while fetching recentActions: %v", err)
// 		ctx.String(http.StatusBadRequest, "Error while getting recent actions")
// 	}

//		ctx.JSON(http.StatusOK, recentActions)
//	}
func RecentActionsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		pipeline := bson.A{
			bson.M{"$group": bson.M{"_id": "$blogEntry.id", "doc": bson.M{"$first": "$$ROOT"}}},
			bson.M{"$replaceRoot": bson.M{"newRoot": "$doc"}},
		}
		result, err := blogCollection.Aggregate(c, pipeline)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}
		var uniqueBlogs []bson.M
		if err = result.All(c, &uniqueBlogs); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, uniqueBlogs)
	}
}

func Subscribe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blogId := ctx.Param("blog_id")
		uid, ok := ctx.Get("uid")
		userID := uid.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user id is not found in token"})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.M{"user_id": userID}
		update := bson.M{"$addToSet": bson.M{"subscribedblogs": blogId}}
		_, err := userCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"InsertedBlogID": blogId})
	}
}
func Unsubscribe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blogId := ctx.Param("blog_id")
		uid, ok := ctx.Get("uid")
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user id is not found in token"})
			return
		}
		filter := bson.M{"user_id": uid}
		update := bson.M{"$pull": bson.M{"subscribedblogs": blogId}}
		_, err := userCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"deletedblogid": blogId})

	}
}

// func BlogById() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		blog_id := ctx.Param("blog_id")
// 		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		count, err := blogCollection.CountDocuments(c, bson.M{"blogEntry.id": blog_id})
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if count < 1 {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog Id does not exist"})
// 			return
// 		}
// 		result, err := blogCollection.Find(c, bson.M{"blogEntry": bson.M{"id": 125364}})
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		var allblogs []bson.M
// 		if err = result.All(c, &allblogs); err != nil {
// 			log.Fatal(err)
// 		}
// 		ctx.JSON(http.StatusOK, allblogs)

// 	}
// }

type BlogEntry struct {
	ID                   int      `bson:"id"`
	OriginalLocale       string   `bson:"originallocale"`
	CreationTimeSeconds  int32    `bson:"creationtimeseconds"`
	AuthorHandle         string   `bson:"authorhandle"`
	Title                string   `bson:"title"`
	Content              string   `bson:"content"`
	Locale               string   `bson:"locale"`
	ModificationTimeSecs int32    `bson:"modificationtimeseconds"`
	AllowViewHistory     bool     `bson:"allowviewhistory"`
	Tags                 []string `bson:"tags"`
	Rating               int32    `bson:"rating"`
}

type Comment struct {
	ID                int    `bson:"id"`
	CreationTimeSecs  int32  `bson:"creationTimeSeconds"`
	CommentatorHandle string `bson:"commentatorHandle"`
	Locale            string `bson:"locale"`
	Text              string `bson:"text"`
	ParentCommentID   int    `bson:"parentCommentId"`
	Rating            int32  `bson:"rating"`
}

type BlogWithComments struct {
	ID          primitive.ObjectID `bson:"_id"`
	TimeSeconds int64              `bson:"timeSeconds"`
	BlogEntry   BlogEntry          `bson:"blogEntry"`
	Comment     *Comment           `bson:"comment"`
}

type Result struct {
	SubscribedBlogs []string `bson:"subscribedblogs"`
}

// Subscribedblogs handler
func Subscribedblogs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := ctx.Get("uid")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
			return
		}
		uidStr, ok := uid.(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Step 1: Retrieve subscribed blogs for the user
		pipeline := bson.A{
			bson.M{"$match": bson.M{"user_id": uidStr}},
			bson.M{"$project": bson.M{"subscribedblogs": 1, "_id": 0}},
		}
		result, err := userCollection.Aggregate(c, pipeline)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var r []Result
		if err := result.All(c, &r); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var subscribedBlogs []int
		for _, res := range r {
			for _, id := range res.SubscribedBlogs {
				intID, err := strconv.Atoi(id)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				subscribedBlogs = append(subscribedBlogs, intID)
			}
		}
		log.Printf("Subscribed blogs for user %s: %v", uidStr, subscribedBlogs)

		// Step 2: Fetch blog documents based on subscribed blog IDs
		pipeline = bson.A{
			bson.M{"$match": bson.M{"blogEntry.id": bson.M{"$in": subscribedBlogs}}},
			bson.M{"$group": bson.M{"_id": "$blogEntry.id", "doc": bson.M{"$first": "$$ROOT"}}},
			bson.M{"$replaceRoot": bson.M{"newRoot": "$doc"}},
		}
		cursor, err := blogCollection.Aggregate(c, pipeline)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c)

		var blogs []BlogWithComments
		if err := cursor.All(c, &blogs); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Retrieved blogs for user %s: %v", uidStr, blogs)
		ctx.JSON(http.StatusOK, blogs)
	}
}

func CommentHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := blogCollection.Find(c, bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(c, &allusers); err != nil {
			log.Fatal(err)
		}
		ctx.JSON(http.StatusOK, allusers)

	}
}
func GetCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx.String(http.StatusNotFound, "No token cookie found")
				return
			}
			// Handle other errors appropriately
			ctx.String(http.StatusInternalServerError, "Error retrieving cookie")
			return
		}

		// Access the cookie value
		ctx.String(http.StatusOK, "Token cookie value: %s", token)
	}
}
func CheckSub() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("blog_id")
		uid, ok := ctx.Get("uid")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.M{
			"user_id":         uid, // Replace with the user ID
			"subscribedblogs": id,  // Replace with the ID you want to check
		}
		count, err := userCollection.CountDocuments(c, filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "false"})
			return
		}
		if count < 1 {
			ctx.JSON(http.StatusOK, gin.H{"mag": "false"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "true"})

	}
}
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := ctx.Get("uid")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"islogin": "false"})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		count, err := userCollection.CountDocuments(c, bson.M{"user_id": uid})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"islogin": "false"})
			return
		}
		if count < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"islogin": "false"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"islogin": "true"})
	}
}
