package database

import (
	"log"
	"time"

	"github.com/backpulse/core/models"
	"github.com/teris-io/shortid"
	"gopkg.in/mgo.v2/bson"
)

// GetArticles : return array of articles of site
func GetArticles(id bson.ObjectId) ([]models.Article, error) {
	var articles []models.Article
	err := DB.C(articlesCollection).Find(bson.M{
		"site_id": id,
	}).All(&articles)
	return articles, err
}

// RemoveArticle : remove article from db
func RemoveArticle(id bson.ObjectId) error {
	err := DB.C(articlesCollection).RemoveId(id)
	return err
}

// GetArticle : return single article based by short_id
func GetArticleByShortID(shortID string) (models.Article, error) {
	var article models.Article
	err := DB.C(articlesCollection).Find(bson.M{
		"short_id": shortID,
	}).One(&article)
	return article, err
}

// GetArticle : return single article based by short_id
func GetArticle(id bson.ObjectId) (models.Article, error) {
	var article models.Article
	err := DB.C(articlesCollection).FindId(id).One(&article)
	return article, err
}

// UpsertArticle : create/update article
func UpsertArticle(article models.Article) (models.Article, error) {
	article.UpdatedAt = time.Now()
	if article.ID == "" {
		article.CreatedAt = time.Now()
		article.ID = bson.NewObjectId()
		article.ShortID, _ = shortid.Generate()
		err := DB.C(articlesCollection).Insert(article)
		log.Println(err)
		return article, err
	}
	err := DB.C(articlesCollection).UpdateId(article.ID, bson.M{
		"$set": bson.M{
			"title":      article.Title,
			"content":    article.Content,
			"updated_at": article.UpdatedAt,
		},
	})
	log.Println(err)

	return article, err
}
