package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tokha04/blogging-platform-api/database"
	"github.com/tokha04/blogging-platform-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCollection *mongo.Collection = database.Client.Database("blogging-platform-api").Collection("blogs")
var validate = validator.New()

func CreateBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var blog models.Blog
		if err := ctx.BindJSON(&blog); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not bind json"})
			return
		}

		validationErr := validate.Struct(blog)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not validate a struct"})
			return
		}

		blog.ID = primitive.NewObjectID()
		blog.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		blog.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, insertionErr := blogCollection.InsertOne(c, blog)
		if insertionErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not insert"})
			return
		}

		ctx.JSON(http.StatusCreated, blog)
	}
}

func UpdateBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := ctx.Param("id")
		blogID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var blog models.Blog
		if err = ctx.BindJSON(&blog); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not bind json"})
			return
		}

		var updateBlog primitive.D

		if blog.Title != "" {
			updateBlog = append(updateBlog, bson.E{Key: "title", Value: blog.Title})
		}
		if blog.Content != "" {
			updateBlog = append(updateBlog, bson.E{Key: "content", Value: blog.Content})
		}
		if blog.Category != "" {
			updateBlog = append(updateBlog, bson.E{Key: "category", Value: blog.Category})
		}
		if blog.Tags != nil {
			updateBlog = append(updateBlog, bson.E{Key: "tags", Value: blog.Tags})
		}
		blog.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateBlog = append(updateBlog, bson.E{Key: "updatedAt", Value: blog.UpdatedAt})

		filter := bson.M{"_id": blogID}
		update := bson.M{"$set": updateBlog}
		res, err := blogCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not update"})
			return
		}

		if res.MatchedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "could not find a blog"})
			return
		}

		var updatedBlog models.Blog
		err = blogCollection.FindOne(c, bson.M{"_id": blogID}).Decode(&updatedBlog)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "could not fetch a blog"})
			return
		}

		ctx.JSON(http.StatusOK, updatedBlog)
	}
}

func DeleteBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := ctx.Param("id")
		blogID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		filter := bson.M{"_id": blogID}
		res, err := blogCollection.DeleteOne(c, filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete"})
			return
		}

		if res.DeletedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "could not find a blog"})
			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{"message": "successfully deleted"})
	}
}

func GetBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := ctx.Param("id")
		blogID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var blog models.Blog
		err = blogCollection.FindOne(c, bson.M{"_id": blogID}).Decode(&blog)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "could not find a blog"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch a blog"})
			}
			return
		}

		ctx.JSON(http.StatusOK, blog)
	}
}

func GetBlogs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{}

		term := ctx.Query("term")
		if term != "" {
			regex := bson.M{"$regex": primitive.Regex{Pattern: term, Options: "i"}}
			filter = bson.M{
				"$or": []bson.M{
					{"title": regex},
					{"content": regex},
					{"tags": regex},
				},
			}
		}

		cursor, err := blogCollection.Find(c, filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find blogs"})
			return
		}
		defer cursor.Close(c)

		var blogs []models.Blog
		for cursor.Next(c) {
			var blog models.Blog
			if err := cursor.Decode(&blog); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch a blog"})
				return
			}
			blogs = append(blogs, blog)
		}

		if err := cursor.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, blogs)
	}
}
