// Package gbp defines the Google Business Profile provider interface and the
// deterministic mock used in this build. A live implementation would call the
// Google Places API; select it with GBP_PROVIDER=live + GOOGLE_PLACES_API_KEY
// (see .env.example). Only the mock ships here — zero keys required.
package gbp

import (
	"fmt"
	"hash/fnv"
	"net/url"
	"strings"

	"zatpatsite/internal/model"
)

// Profile is what a GBP lookup returns — enough to autofill the whole form.
type Profile struct {
	Name     string            `json:"name"`
	Category string            `json:"category"`
	City     string            `json:"city"`
	Address  string            `json:"address"`
	Phone    string            `json:"phone"`
	WhatsApp string            `json:"whatsapp"`
	MapsURL  string            `json:"mapsUrl"`
	Hours    [7]model.DayHours `json:"hours"`
	Services []model.Service   `json:"services"`
	Rating   float64           `json:"rating"`
	Reviews  []model.Review    `json:"reviews"`
}

// Provider fetches a business profile by name + city.
type Provider interface {
	Fetch(name, city, category string) (Profile, error)
	Mode() string
}

// Mock is the deterministic offline provider: the same (name, city) always
// yields the same profile, seeded from an FNV-1a hash.
type Mock struct{}

func (Mock) Mode() string { return "mock" }

// rng is a tiny splitmix64 so successive picks stay deterministic.
type rng struct{ s uint64 }

func (r *rng) next() uint64 {
	r.s += 0x9e3779b97f4a7c15
	z := r.s
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func (r *rng) pick(n int) int { return int(r.next() % uint64(n)) }

func seed(name, city string) *rng {
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(strings.TrimSpace(name))))
	h.Write([]byte("|"))
	h.Write([]byte(strings.ToLower(strings.TrimSpace(city))))
	return &rng{s: h.Sum64()}
}

var streets = []string{
	"MG Road", "Station Road", "Nehru Marg", "Gandhi Chowk", "Link Road",
	"Mall Road", "Subhash Marg", "Ring Road", "Tilak Path", "Patel Chowk",
}

var localities = []string{
	"Shivaji Nagar", "Rajendra Nagar", "Civil Lines", "Sadar Bazaar",
	"Gandhi Nagar", "Lajpat Colony", "Anand Vihar", "Malviya Nagar",
	"Kalyan Nagar", "Ashok Nagar",
}

var reviewers = []string{
	"Priya Sharma", "Rahul Verma", "Sneha Kulkarni", "Arjun Nair",
	"Kavita Reddy", "Imran Shaikh", "Deepak Joshi", "Ananya Iyer",
	"Vikram Singh", "Meera Patel", "Rohit Agarwal", "Farhan Qureshi",
}

var reviewTexts = map[string][]string{
	"salon": {
		"Best haircut I've had in years. They actually listen to what you want.",
		"Clean, friendly and always on time. My whole family comes here now.",
		"The bridal package was worth every rupee. Highly recommended.",
		"Booked on WhatsApp, walked in, done in 40 minutes. Zero fuss.",
	},
	"restaurant": {
		"The dal makhani here is better than my mother-in-law's. Don't tell her.",
		"Generous portions, quick service, and the masala papad is a must.",
		"We host every family dinner here. Consistently excellent for 3 years.",
		"Ordered on WhatsApp for pickup — food was hot and ready on the dot.",
	},
	"kirana": {
		"Uncle ji keeps everything in stock and delivers within the hour.",
		"Fair prices, fresh atta, and they remember my monthly list. Gem of a shop.",
		"Sent a WhatsApp message at 8pm, groceries at my door by 8:30.",
		"The only shop in the colony that never runs out of anything.",
	},
	"coaching": {
		"My daughter jumped from 62% to 88% in one year. The faculty is superb.",
		"Small batches, weekly tests, honest feedback to parents. Rare these days.",
		"Cleared my banking exam in the first attempt thanks to these classes.",
		"Doubt-clearing sessions till late evening. Genuinely dedicated teachers.",
	},
	"clinic": {
		"Doctor sahab explains everything patiently. Never pushes unnecessary tests.",
		"Clean clinic, short wait, and they follow up on WhatsApp after visits.",
		"Been our family doctor for a decade. Completely trustworthy.",
		"Got an appointment the same morning I called. Very well managed.",
	},
	"boutique": {
		"Got my sister's wedding lehenga stitched here. Flawless fitting.",
		"They copied a design from a photo perfectly. True craftsmanship.",
		"Blouse fitting done in two days flat, right before the function.",
		"Fabric selection is gorgeous and the prices are honest.",
	},
	"gym": {
		"Trainers actually train you here instead of scrolling their phones.",
		"Clean equipment, no waiting, and the morning batch is super motivating.",
		"Lost 9 kg in 5 months with their diet plan. Life changing.",
		"Best value gym in the area. The annual plan is a steal.",
	},
	"bakery": {
		"Their pineapple cake made my son's birthday. Ordered again the next week.",
		"Fresh pav every morning at 7. The whole lane smells amazing.",
		"Eggless chocolate truffle that actually tastes premium. Rare find.",
		"Ordered a photo cake on WhatsApp at night, delivered by noon. Perfect.",
	},
}

var serviceMenus = map[string][]model.Service{
	"salon": {
		{Name: "Haircut (Men)", Price: 250}, {Name: "Haircut (Women)", Price: 550},
		{Name: "Hair Spa", Price: 900}, {Name: "Facial (Gold)", Price: 1200},
		{Name: "Bridal Makeup", Price: 15000}, {Name: "Beard Styling", Price: 180},
		{Name: "Hair Colour (Ammonia-free)", Price: 1800},
	},
	"restaurant": {
		{Name: "Paneer Butter Masala", Price: 280}, {Name: "Dal Makhani", Price: 240},
		{Name: "Butter Naan (2 pc)", Price: 90}, {Name: "Veg Thali", Price: 220},
		{Name: "Chicken Biryani", Price: 320}, {Name: "Masala Dosa", Price: 140},
		{Name: "Gulab Jamun (2 pc)", Price: 80},
	},
	"kirana": {
		{Name: "Aashirvaad Atta 10kg", Price: 425}, {Name: "Tata Salt 1kg", Price: 30},
		{Name: "Fortune Sunflower Oil 1L", Price: 145}, {Name: "Basmati Rice 5kg", Price: 620},
		{Name: "Amul Butter 500g", Price: 275}, {Name: "Free Home Delivery (2 km)", Price: 0},
		{Name: "Monthly Ration Bundle", Price: 3500},
	},
	"coaching": {
		{Name: "Class 10 Maths + Science", Price: 18000}, {Name: "Class 12 PCM", Price: 24000},
		{Name: "JEE Foundation (Class 9)", Price: 30000}, {Name: "Spoken English (3 months)", Price: 6500},
		{Name: "Banking / SSC Batch", Price: 12000}, {Name: "Crash Course (Boards)", Price: 9000},
	},
	"clinic": {
		{Name: "General Consultation", Price: 300}, {Name: "Follow-up Visit", Price: 150},
		{Name: "Blood Pressure Check", Price: 50}, {Name: "Diabetes Package", Price: 800},
		{Name: "Child Vaccination", Price: 600}, {Name: "Home Visit (within 3 km)", Price: 700},
	},
	"boutique": {
		{Name: "Blouse Stitching", Price: 450}, {Name: "Salwar Suit Stitching", Price: 850},
		{Name: "Lehenga (Custom)", Price: 6500}, {Name: "Saree Fall & Pico", Price: 120},
		{Name: "Kurta (Men)", Price: 700}, {Name: "Designer Gown", Price: 4500},
	},
	"gym": {
		{Name: "Monthly Membership", Price: 1200}, {Name: "Quarterly Plan", Price: 3000},
		{Name: "Annual Plan", Price: 9999}, {Name: "Personal Training (month)", Price: 5000},
		{Name: "Diet Consultation", Price: 800}, {Name: "Zumba Batch (month)", Price: 1500},
	},
	"bakery": {
		{Name: "Pineapple Cake (500g)", Price: 350}, {Name: "Chocolate Truffle (500g)", Price: 450},
		{Name: "Photo Cake (1kg)", Price: 900}, {Name: "Fresh Pav (6 pc)", Price: 30},
		{Name: "Veg Puff", Price: 25}, {Name: "Cookies Box (250g)", Price: 120},
	},
}

var hoursByCategory = map[string][7]model.DayHours{
	"salon":      weekOf("10:00", "20:00", 1), // Monday off
	"restaurant": weekOf("11:00", "23:00", -1),
	"kirana":     weekOf("07:30", "21:30", -1),
	"coaching":   weekOf("07:00", "19:30", 0), // Sunday off
	"clinic":     weekOf("09:30", "19:30", 0),
	"boutique":   weekOf("10:30", "20:30", 1),
	"gym":        weekOf("06:00", "22:00", -1),
	"bakery":     weekOf("08:00", "21:00", -1),
}

// weekOf builds a uniform week with an optional closed day (-1 = none).
func weekOf(open, close string, closedDay int) [7]model.DayHours {
	var h [7]model.DayHours
	for i := range h {
		h[i] = model.DayHours{Open: open, Close: close}
	}
	if closedDay >= 0 && closedDay < 7 {
		h[closedDay] = model.DayHours{Closed: true}
	}
	return h
}

// Fetch deterministically synthesizes a believable Indian business profile.
func (Mock) Fetch(name, city, category string) (Profile, error) {
	r := seed(name, city)
	if !model.ValidCategory(category) {
		category = model.Categories[r.pick(len(model.Categories))]
	}
	shopNo := 1 + r.pick(60)
	street := streets[r.pick(len(streets))]
	locality := localities[r.pick(len(localities))]
	pin := 110001 + r.pick(689999)
	phone := fmt.Sprintf("98%08d", r.next()%1_0000_0000)
	wa := phone
	if r.pick(3) == 0 { // sometimes a separate WhatsApp number
		wa = fmt.Sprintf("98%08d", r.next()%1_0000_0000)
	}

	menu := serviceMenus[category]
	count := 5 + r.pick(2)
	if count > len(menu) {
		count = len(menu)
	}
	start := r.pick(len(menu))
	services := make([]model.Service, 0, count)
	for i := 0; i < count; i++ {
		services = append(services, menu[(start+i)%len(menu)])
	}

	rating := 3.9 + float64(r.pick(10))/10.0
	texts := reviewTexts[category]
	nameStart := r.pick(len(reviewers))
	textStart := r.pick(len(texts))
	reviews := make([]model.Review, 0, 3)
	for i := 0; i < 3; i++ {
		rv := rating + float64(r.pick(2)) // 4s and 5s around the aggregate
		if rv > 5 {
			rv = 5
		}
		reviews = append(reviews, model.Review{
			Author: reviewers[(nameStart+i*3)%len(reviewers)],
			Rating: float64(int(rv)), // whole stars per review
			Text:   texts[(textStart+i)%len(texts)],
		})
	}

	trimmedName := strings.TrimSpace(name)
	trimmedCity := strings.TrimSpace(city)
	return Profile{
		Name:     trimmedName,
		Category: category,
		City:     trimmedCity,
		Address:  fmt.Sprintf("Shop %d, %s, %s, %s - %06d", shopNo, street, locality, trimmedCity, pin),
		Phone:    phone,
		WhatsApp: wa,
		MapsURL: "https://www.google.com/maps/search/?api=1&query=" +
			url.QueryEscape(trimmedName+" "+trimmedCity),
		Hours:    hoursByCategory[category],
		Services: services,
		Rating:   rating,
		Reviews:  reviews,
	}, nil
}
