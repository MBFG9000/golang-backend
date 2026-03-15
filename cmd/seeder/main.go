package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/MBFG9000/golang-backend/internal/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type options struct {
	usersCount   int
	maxFriends   int
	seed         int64
	truncateData bool
}

func main() {
	opts := parseFlags()
	loadEnvOptional()

	connURL := config.GetConnURL(config.GetConfig())
	db, err := sqlx.Connect("postgres", connURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	rng := rand.New(rand.NewSource(opts.seed))
	userIDs, insertedUsers, err := seedUsers(db, rng, opts)
	if err != nil {
		log.Fatalf("seed users: %v", err)
	}

	insertedFriends, err := seedFriends(db, rng, userIDs, opts.maxFriends)
	if err != nil {
		log.Fatalf("seed friends: %v", err)
	}

	fmt.Printf("Seed completed: users=%d, friend_links=%d, seed=%d\n", insertedUsers, insertedFriends, opts.seed)
}

func parseFlags() options {
	opts := options{}
	defaultSeed := time.Now().UnixNano()

	flag.IntVar(&opts.usersCount, "users", 100, "number of users to generate")
	flag.IntVar(&opts.maxFriends, "friends", 5, "max friends per user")
	flag.Int64Var(&opts.seed, "seed", defaultSeed, "random seed for deterministic generation")
	flag.BoolVar(&opts.truncateData, "truncate", false, "truncate existing users and user_friends before seeding")
	flag.Parse()

	if opts.usersCount < 1 {
		log.Fatal("--users must be >= 1")
	}
	if opts.maxFriends < 0 {
		log.Fatal("--friends must be >= 0")
	}

	return opts
}

func loadEnvOptional() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("load .env: %v", err)
		}
	}
}

func seedUsers(db *sqlx.DB, rng *rand.Rand, opts options) ([]uuid.UUID, int, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, 0, fmt.Errorf("begin users tx: %w", err)
	}
	defer tx.Rollback()

	if opts.truncateData {
		if _, err := tx.Exec("TRUNCATE TABLE user_friends, users RESTART IDENTITY"); err != nil {
			return nil, 0, fmt.Errorf("truncate tables: %w", err)
		}
	}

	stmt, err := tx.Preparex(`
		INSERT INTO users (first_name, last_name, email, phone, city, country, zip, gender, birth_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`)
	if err != nil {
		return nil, 0, fmt.Errorf("prepare user insert: %w", err)
	}
	defer stmt.Close()

	ids := make([]uuid.UUID, 0, opts.usersCount)
	for i := 0; i < opts.usersCount; i++ {
		firstName := firstNames[rng.Intn(len(firstNames))]
		lastName := lastNames[rng.Intn(len(lastNames))]
		email := fmt.Sprintf("dev.user.%06d.%d@example.com", i+1, opts.seed%100000)
		phone := fmt.Sprintf("+1-555-%04d-%04d", (i/10000)%10000, i%10000)
		city := cities[rng.Intn(len(cities))]
		country := countries[rng.Intn(len(countries))]
		zip := fmt.Sprintf("%05d", 10000+rng.Intn(89999))

		var gender any
		if rng.Intn(100) < 80 {
			gender = genders[rng.Intn(len(genders))]
		}

		var birthDate any
		if rng.Intn(100) < 85 {
			birthDate = randomBirthDate(rng)
		}

		var id uuid.UUID
		if err := stmt.QueryRow(firstName, lastName, email, phone, city, country, zip, gender, birthDate).Scan(&id); err != nil {
			return nil, 0, fmt.Errorf("insert user #%d: %w", i+1, err)
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("commit users tx: %w", err)
	}

	return ids, len(ids), nil
}

func seedFriends(db *sqlx.DB, rng *rand.Rand, userIDs []uuid.UUID, maxFriends int) (int, error) {
	if len(userIDs) <= 1 || maxFriends == 0 {
		return 0, nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("begin friends tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`
		INSERT INTO user_friends (user_id, friend_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, friend_id) DO NOTHING
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare friend insert: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for _, userID := range userIDs {
		targetCount := rng.Intn(maxFriends + 1)
		chosen := make(map[uuid.UUID]struct{}, targetCount)
		attempts := 0
		maxAttempts := len(userIDs) * 2

		for len(chosen) < targetCount && attempts < maxAttempts {
			attempts++
			friendID := userIDs[rng.Intn(len(userIDs))]
			if friendID == userID {
				continue
			}
			chosen[friendID] = struct{}{}
		}

		for friendID := range chosen {
			res, err := stmt.Exec(userID, friendID)
			if err != nil {
				return 0, fmt.Errorf("insert friend link %s -> %s: %w", userID, friendID, err)
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("friend rows affected: %w", err)
			}
			inserted += int(rows)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit friends tx: %w", err)
	}

	return inserted, nil
}

func randomBirthDate(rng *rand.Rand) time.Time {
	start := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)
	deltaDays := int(end.Sub(start).Hours() / 24)
	return start.AddDate(0, 0, rng.Intn(deltaDays+1))
}

var firstNames = []string{
	"Alex", "Maria", "Dmitry", "Anna", "Sergey", "Olga", "Ivan", "Nina", "Pavel", "Elena",
}

var lastNames = []string{
	"Ivanov", "Petrova", "Sidorov", "Kozlova", "Morozov", "Smirnova", "Popov", "Orlova", "Volkov", "Romanova",
}

var cities = []string{
	"Moscow", "Saint Petersburg", "Kazan", "Novosibirsk", "Yekaterinburg", "Samara", "Perm",
}

var countries = []string{
	"Russia", "Kazakhstan", "Belarus", "Georgia", "Armenia",
}

var genders = []string{
	"male", "female",
}
