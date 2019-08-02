package models

func MigrationDB() {
	db := Database()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&AuthTable{})
}
