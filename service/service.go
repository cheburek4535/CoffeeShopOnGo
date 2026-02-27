package service
import ("log"
)

func InitDB() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	if err := db.DB.AutoMigrate(&User{}); err != nil {}
}


var StaffService = NewDefaultCrud[Staff, uint]


func GetStaff(id uint) (Staff) {
	crud = StaffService
	staff, err := StaffService.Get
}