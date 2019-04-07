package database

import (
	"time"

	"github.com/backpulse/core/models"
	"github.com/teris-io/shortid"
	"gopkg.in/mgo.v2/bson"
)

// GetProject : return project using shortid
func GetProjectByShortID(shortID string) (models.Project, error) {
	var project models.Project
	err := DB.C(projectsCollection).Find(bson.M{
		"short_id": shortID,
	}).One(&project)

	project.Title = getDefaultProjectTitle(project)

	return project, err
}

// GetProject : return project using shortid
func GetProject(id bson.ObjectId) (models.Project, error) {
	var project models.Project
	err := DB.C(projectsCollection).FindId(id).One(&project)

	project.Title = getDefaultProjectTitle(project)

	return project, err
}

// RemoveProject : remove project from db
func RemoveProject(id bson.ObjectId) error {
	err := DB.C(projectsCollection).RemoveId(id)
	return err
}

// GetProjects : return projects of site
func GetProjects(id bson.ObjectId) ([]models.Project, error) {
	var projects []models.Project
	err := DB.C(projectsCollection).Find(bson.M{
		"site_id": id,
	}).All(&projects)

	for i := range projects {
		title := getDefaultProjectTitle(projects[i])
		projects[i].Title = title
	}

	return projects, err
}

// UpsertProject Update or insert project
func UpsertProject(project models.Project) error {
	if project.ID == "" {
		project.CreatedAt = time.Now()
		project.ID = bson.NewObjectId()
		project.ShortID, _ = shortid.Generate()
	}
	project.UpdatedAt = time.Now()
	_, err := DB.C(projectsCollection).UpsertId(project.ID, bson.M{
		"$set": project,
	})
	return err
}

// GetDefaultProjectTitle : return default project title
func getDefaultProjectTitle(project models.Project) string {
	for i := range project.Titles {
		if project.Titles[i].LanguageCode == "en" {
			return project.Titles[i].Content
		}
	}
	return project.Titles[0].Content
}
