package seeders

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gatehide/gatehide-api/config"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// GamenetSeeder handles seeding gamenet data
type GamenetSeeder struct {
	db *sql.DB
}

// NewGamenetSeeder creates a new gamenet seeder instance
func NewGamenetSeeder(cfg *config.Config) (*GamenetSeeder, error) {
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &GamenetSeeder{db: db}, nil
}

// GamenetData represents gamenet data
type GamenetData struct {
	Name              string
	OwnerName         string
	OwnerMobile       string
	Address           string
	Email             string
	Password          string
	LicenseAttachment *string
}

// SeedGamenets is the public seeder function that can be called by the registry
func SeedGamenets(cfg *config.Config) error {
	seeder, err := NewGamenetSeeder(cfg)
	if err != nil {
		return fmt.Errorf("failed to create gamenet seeder: %w", err)
	}
	defer seeder.Close()

	// Generate 25 gamenets for pagination testing
	gamenets := generateGamenets(25)

	return seeder.seedGamenets(gamenets)
}

// generateGamenets creates a slice of gamenet data for testing
func generateGamenets(count int) []GamenetData {
	rand.Seed(time.Now().UnixNano())

	// Persian names for variety
	names := []string{
		"گیم نت ستاره", "گیم نت طلایی", "گیم نت الماس", "گیم نت طوفان", "گیم نت اژدها",
		"گیم نت شوالیه", "گیم نت فانتزی", "گیم نت افسانه", "گیم نت قهرمان", "گیم نت پیروزی",
		"گیم نت نبرد", "گیم نت جنگجو", "گیم نت شجاع", "گیم نت قوی", "گیم نت سریع",
		"گیم نت هوشمند", "گیم نت خلاق", "گیم نت نوآور", "گیم نت پیشرفته", "گیم نت مدرن",
		"گیم نت کلاسیک", "گیم نت سنتی", "گیم نت محلی", "گیم نت خانوادگی", "گیم نت دوستانه",
	}

	ownerNames := []string{
		"علی احمدی", "محمد رضایی", "حسن کریمی", "رضا محمدی", "امیر حسینی",
		"سعید نوری", "مهدی صادقی", "حسین علیزاده", "احمد رحمانی", "علی اکبری",
		"محمد جوادی", "حسن مرادی", "رضا کرمانی", "امیر تهرانی", "سعید اصفهانی",
		"مهدی شیرازی", "حسین تبریزی", "احمد مشهدی", "علی قمی", "محمد یزدی",
		"حسن کرمانشاهی", "رضا همدانی", "امیر کرجی", "سعید قزوینی", "مهدی سمنانی",
	}

	cities := []string{
		"تهران", "مشهد", "اصفهان", "شیراز", "تبریز", "کرج", "اهواز", "قم", "کرمانشاه", "ارومیه",
		"زاهدان", "رشت", "کرمان", "یزد", "اردبیل", "بندرعباس", "گرگان", "ساری", "بوشهر", "خرم آباد",
		"سنندج", "زنجان", "قزوین", "کاشان", "نجف آباد",
	}

	streets := []string{
		"خیابان انقلاب", "خیابان ولیعصر", "خیابان کریمخان", "خیابان آزادی", "خیابان طالقانی",
		"خیابان شریعتی", "خیابان جردن", "خیابان ونک", "خیابان پاسداران", "خیابان نیایش",
		"خیابان فرشته", "خیابان پارک وی", "خیابان الهیه", "خیابان قیطریه", "خیابان سعادت آباد",
		"خیابان شهرک غرب", "خیابان پونک", "خیابان ستارخان", "خیابان صادقیه", "خیابان آریاشهر",
		"خیابان میرداماد", "خیابان گاندی", "خیابان فاطمی", "خیابان مطهری", "خیابان کریمخان",
	}

	var gamenets []GamenetData

	for i := 0; i < count; i++ {
		// Generate random data
		name := names[i%len(names)]
		if i >= len(names) {
			name = fmt.Sprintf("%s %d", name, i+1)
		}

		ownerName := ownerNames[i%len(ownerNames)]
		city := cities[i%len(cities)]
		street := streets[i%len(streets)]

		// Generate mobile number (09xxxxxxxxx)
		mobile := fmt.Sprintf("09%d%d%d%d%d%d%d%d%d",
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))

		// Generate email
		email := fmt.Sprintf("gamenet%d@example.com", i+1)

		// Generate address
		address := fmt.Sprintf("%s، %s، پلاک %d", city, street, rand.Intn(999)+1)

		// Generate password (8 digits)
		password := fmt.Sprintf("%d%d%d%d%d%d%d%d",
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))

		// Randomly assign license attachment (50% chance)
		var licenseAttachment *string
		if rand.Intn(2) == 0 {
			license := fmt.Sprintf("/uploads/licenses/gamenet_%d_license.pdf", i+1)
			licenseAttachment = &license
		}

		gamenet := GamenetData{
			Name:              name,
			OwnerName:         ownerName,
			OwnerMobile:       mobile,
			Address:           address,
			Email:             email,
			Password:          password,
			LicenseAttachment: licenseAttachment,
		}

		gamenets = append(gamenets, gamenet)
	}

	return gamenets
}

// seedGamenets seeds gamenets into the database
func (s *GamenetSeeder) seedGamenets(gamenets []GamenetData) error {
	// Check if gamenets already exist
	var count int
	checkQuery := "SELECT COUNT(*) FROM gamenets"
	err := s.db.QueryRow(checkQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing gamenets: %w", err)
	}

	if count > 0 {
		log.Printf("Gamenets already exist (%d records), skipping seeder...", count)
		return nil
	}

	// Insert gamenets
	insertQuery := `
		INSERT INTO gamenets (name, owner_name, owner_mobile, address, email, password, license_attachment, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	successCount := 0
	for i, gamenet := range gamenets {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(gamenet.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for gamenet %d: %v", i+1, err)
			continue
		}

		_, err = s.db.Exec(insertQuery,
			gamenet.Name,
			gamenet.OwnerName,
			gamenet.OwnerMobile,
			gamenet.Address,
			gamenet.Email,
			string(hashedPassword),
			gamenet.LicenseAttachment,
		)

		if err != nil {
			log.Printf("Failed to insert gamenet %d: %v", i+1, err)
			continue
		}

		successCount++
	}

	log.Printf("✅ Successfully seeded %d gamenets", successCount)
	return nil
}

// Close closes the database connection
func (s *GamenetSeeder) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
