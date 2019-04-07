package database

import (
	"time"

	"github.com/backpulse/core/constants"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"gopkg.in/mgo.v2/bson"
)

// GetSiteTotalSize : return site size in MB
func GetSiteTotalSize(siteID bson.ObjectId) float64 {
	photos, _ := GetSitePhotos(siteID)

	var totalSize float64
	for _, photo := range photos {
		totalSize = totalSize + photo.Size
	}

	files, _ := GetSiteFiles(siteID)
	for _, file := range files {
		totalSize = totalSize + file.Size
	}

	return totalSize
}

// TransferSite : Transfer the site to another user
func TransferSite(site models.Site, lastOwner models.User, newOwner models.User) error {
	collections := []string{contactCollection, aboutCollection, galleriesCollection, projectsCollection, photosCollection}
	for _, c := range collections {
		DB.C(c).Update(bson.M{
			"site_id": site.ID,
		}, bson.M{
			"$set": bson.M{
				"owner_id": newOwner.ID,
			},
		})
	}

	err := DB.C(sitesCollection).UpdateId(site.ID, bson.M{
		"$set": bson.M{
			"owner_id": newOwner.ID,
		},
	})

	RemoveCollaborator(site.ID, newOwner.Email)

	AddCollaborator(site.ID, models.Collaborator{
		Email: lastOwner.Email,
		Role:  "collaborator",
	})

	return err
}

// GetSitesOfUser : return array of Site of specific user
func GetSitesOfUser(user models.User) ([]models.Site, error) {
	var sites []models.Site
	err := DB.C(sitesCollection).Find(bson.M{
		"$or": []bson.M{
			{
				"owner_id": user.ID,
			},
			{
				"collaborators": models.Collaborator{
					Email: user.Email,
					Role:  "collaborator",
				},
			},
		},
	}).All(&sites)

	/* Dynamicly set favorite */
	for i := range sites {
		for _, id := range user.FavoriteSites {
			if sites[i].ID == id {
				sites[i].Favorite = true
			}
		}
		for _, collaborator := range sites[i].Collaborators {
			if collaborator.Email == user.Email {
				sites[i].Role = "collaborator"
			}
		}
		if sites[i].Role == "" {
			sites[i].Role = "owner"
		}
	}

	return sites, err
}

// GetOwnedSites : return only owned sites (and not the ones you might collaborate with)
func GetOwnedSites(id bson.ObjectId) ([]models.Site, error) {
	var sites []models.Site
	err := DB.C(sitesCollection).Find(bson.M{
		"owner_id": id,
	}).All(&sites)
	return sites, err
}

// AddCollaborator : add collaborator to site
func AddCollaborator(id bson.ObjectId, collaborator models.Collaborator) error {
	err := DB.C(sitesCollection).UpdateId(id, bson.M{
		"$push": bson.M{
			"collaborators": collaborator,
		},
	})
	return err
}

// RemoveCollaborator : remove collaborator from site
func RemoveCollaborator(id bson.ObjectId, email string) error {
	err := DB.C(sitesCollection).UpdateId(id, bson.M{
		"$pull": bson.M{
			"collaborators": bson.M{
				"role":  "collaborator",
				"email": email,
			}},
	})
	return err
}

// GetSiteByName : return specific site
func GetSiteByName(name string) (models.Site, error) {
	var site models.Site
	err := DB.C(sitesCollection).Find(bson.M{
		"name": name,
	}).One(&site)
	return site, err
}

// GetSiteByID : return specific site by ObjectID
func GetSiteByID(id bson.ObjectId) (models.Site, error) {
	var site models.Site
	err := DB.C(sitesCollection).FindId(id).One(&site)
	return site, err
}

// AddModule : add module to site
func AddModule(site models.Site, module string) error {
	err := DB.C(sitesCollection).UpdateId(site.ID, bson.M{
		"$push": bson.M{
			"modules": constants.Module(module),
		},
	})
	return err
}

// RemoveModule : remove module from site
func RemoveModule(site models.Site, module string) error {
	err := DB.C(sitesCollection).UpdateId(site.ID, bson.M{
		"$pull": bson.M{
			"modules": constants.Module(module),
		},
	})

	if err != nil {
		return err
	}

	//TODO: Remove data from other modules (videos, articles, music...)
	// Or ask user to delete it first idk

	if module == "galleries" {
		galleries, _ := GetGalleries(site.ID)

		for _, g := range galleries {
			photos, _ := GetGalleryPhotos(g.ID)
			utils.RemoveGoogleCloudPhotos(photos)
		}
		_, err := DB.C(photosCollection).RemoveAll(bson.M{
			"site_id": site.ID,
		})
		_, err = DB.C(galleriesCollection).RemoveAll(bson.M{
			"site_id": site.ID,
		})
		return err
	}

	if module == "projects" {
		_, err := DB.C(projectsCollection).RemoveAll(bson.M{
			"site_id": site.ID,
		})
		return err
	}
	//TODO remove articles
	return err
}

// UpdateSite : update site data
func UpdateSite(id bson.ObjectId, data models.Site) error {

	err := DB.C(sitesCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"updated_at":   time.Now(),
			"display_name": data.DisplayName,
			"name":         data.Name,
		},
	})
	return err
}

// CreateSite : insert site in db
func CreateSite(site models.Site) (models.Site, error) {
	site.CreatedAt = time.Now()
	site.UpdatedAt = time.Now()
	site.ID = bson.NewObjectId()

	err := DB.C(sitesCollection).Insert(site)

	UpdateContact(site.ID, models.ContactContent{OwnerID: site.OwnerID, SiteID: site.ID})
	UpdateAbout(site.ID, models.AboutContent{OwnerID: site.OwnerID, SiteID: site.ID})

	return site, err
}

// DeleteSite : complete erase
func DeleteSite(site models.Site) error {
	collections := []string{contactCollection, aboutCollection, galleriesCollection, projectsCollection, photosCollection}
	for _, c := range collections {
		DB.C(c).RemoveAll(bson.M{
			"site_id": site.ID,
		})
	}
	return DB.C(sitesCollection).RemoveId(site.ID)
}

// SiteExists : check if site exists
func SiteExists(name string) bool {
	count, _ := DB.C(sitesCollection).Find(bson.M{
		"name": name,
	}).Count()
	return count > 0
}

// SetSiteFavorite : set site favorite
func SetSiteFavorite(user models.User, siteID bson.ObjectId) error {
	isFavorite := false
	for _, sID := range user.FavoriteSites {
		if siteID == sID {
			isFavorite = true
		}
	}
	if isFavorite {
		return DB.C(usersCollection).UpdateId(user.ID, bson.M{
			"$pull": bson.M{
				"favorite_sites": siteID,
			},
		})
	}
	return DB.C(usersCollection).UpdateId(user.ID, bson.M{
		"$push": bson.M{
			"favorite_sites": siteID,
		},
	})
}
