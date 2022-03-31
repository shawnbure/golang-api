package storage

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddCollection(collection *entities.Collection) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func UpdateCollectionWhere(collection *entities.Collection, toUpdate map[string]interface{}, where string, args ...interface{}) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Model(collection).Where(where, args...).Updates(toUpdate)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func UpdateCollection(collection *entities.Collection) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateCollectionProfileWhereId(collectionId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("collections").Where("id = ?", collectionId).Update("profile_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateCollectionCoverWhereId(collectionId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("collections").Where("id = ?", collectionId).Update("cover_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func GetCollectionByAddr(addr string) (*entities.Collection, error) {
	var collection entities.Collection

	if len(addr) < 10 {
		return &collection, nil
	}

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}
	collection.ContractAddress = addr
	txRead := database.Where(collection).Find(&collection)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}
func GetCollectionById(id uint64) (*entities.Collection, error) {
	var collection entities.Collection

	if id == 0 {
		return &collection, nil
	}

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionsCreatedBy(id uint64) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collections, "creator_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionByName(name string) (*entities.Collection, error) {
	var collection entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "name = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionWithNameILike(name string) (*entities.Collection, error) {
	var collection entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "name ILIKE ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionByTokenId(tokenId string) (*entities.Collection, error) {
	var collection entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "token_id = ?", tokenId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionsWithOffsetLimit(offset int, limit int, flags []string) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit)
	for _, flag := range flags {
		txRead.Where(datatypes.JSONQuery("flags").HasKey(flag))
	}

	txRead.Order("priority desc").Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetAllCollections() ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Order("created_at desc").Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

/*
func GetAllCollectionAccounts() ([]entities.CollectionAccount, error) {

	var collections []entities.Collection

	var collectionAccounts []entities.CollectionAccount

	if err = db.Joins("JOIN artist_movies on artist_movies.artist_id=artists.id").
		Joins("JOIN movies on artist_movies.movie_id=movies.id").Where("movies.title=?", "Nayagan").
		Group("artists.id").Find(&collections).Error; err != nil {
		log.Fatal(err)
	}

	for _, ar := range artists {
		fmt.Println(ar.Name)
	}
}
*/

func GetCollectionsWithNameAlikeWithLimit(name string, limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("name ILIKE ?", name).Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsByCreatorIdWithOffsetLimit(creatorId uint64, offset int, limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&collections, "creator_id = ?", creatorId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsVerified(limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	//is_verifed and create_at in desc in order (most recent first)
	txRead := database.Limit(limit).Find(&collections, "is_verified = true AND profile_image_link <> ''").Order("created_at desc")
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsNoteworthy(limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Find(&collections, "type = 2 AND profile_image_link <> ''").Order("priority desc")
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsTrending(limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	/*
		SELECT *
		FROM public.collections
		where is_verified <> true AND type <> 2 AND profile_image_link <> ''
		ORDER BY created_at DESC
	*/

	//Do not include in verified and noteworthy accounts
	//right now, we are just getting the recently added collection
	//TODO: determine the metrics for
	txRead := database.Limit(limit).Find(&collections, "is_verified <> true AND type <> 2 AND profile_image_link <> ''").Order("created_at desc")
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}
