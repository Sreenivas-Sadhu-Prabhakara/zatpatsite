package generator

import "strings"

// strset holds every user-visible template string for one language.
type strset struct {
	LangCode       string
	HeadingAbout   string
	HeadingHours   string
	HeadingReviews string
	HeadingVisit   string
	OpenNow        string
	ClosedNow      string
	ClosedDay      string
	Today          string
	CallCTA        string
	WhatsAppCTA    string
	DirectionsCTA  string
	ReviewsOn      string
	MadeWith       string
	Days           [7]string
	WAMessage      string // {name} placeholder
	MetaDesc       string // {name} {category} {city} placeholders
	AddressLabel   string
	PhoneLabel     string
	ChatOnWhatsApp string
}

var langEN = strset{
	LangCode:       "en",
	HeadingAbout:   "About Us",
	HeadingHours:   "Opening Hours",
	HeadingReviews: "What Customers Say",
	HeadingVisit:   "Find Us",
	OpenNow:        "Open now",
	ClosedNow:      "Closed now",
	ClosedDay:      "Closed",
	Today:          "Today",
	CallCTA:        "Call Now",
	WhatsAppCTA:    "WhatsApp Us",
	DirectionsCTA:  "Get Directions",
	ReviewsOn:      "Google reviews",
	MadeWith:       "Website by ZatpatSite",
	Days: [7]string{
		"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
	},
	WAMessage:      "Namaste! I found {name} on your website and would like to know more.",
	MetaDesc:       "{name} — {category} in {city}. {tagline} Call us or message on WhatsApp.",
	AddressLabel:   "Address",
	PhoneLabel:     "Phone",
	ChatOnWhatsApp: "Chat on WhatsApp",
}

var langHI = strset{
	LangCode:       "hi",
	HeadingAbout:   "हमारे बारे में",
	HeadingHours:   "खुलने का समय",
	HeadingReviews: "ग्राहक क्या कहते हैं",
	HeadingVisit:   "हम यहाँ हैं",
	OpenNow:        "अभी खुला है",
	ClosedNow:      "अभी बंद है",
	ClosedDay:      "बंद",
	Today:          "आज",
	CallCTA:        "कॉल करें",
	WhatsAppCTA:    "व्हाट्सऐप करें",
	DirectionsCTA:  "रास्ता देखें",
	ReviewsOn:      "Google समीक्षाएँ",
	MadeWith:       "वेबसाइट ZatpatSite से बनी",
	Days: [7]string{
		"रविवार", "सोमवार", "मंगलवार", "बुधवार", "गुरुवार", "शुक्रवार", "शनिवार",
	},
	WAMessage:      "नमस्ते! मुझे आपकी वेबसाइट से {name} के बारे में पता चला। कृपया अधिक जानकारी दें।",
	MetaDesc:       "{name} — {city} में {category}। {tagline} कॉल करें या व्हाट्सऐप पर संदेश भेजें।",
	AddressLabel:   "पता",
	PhoneLabel:     "फ़ोन",
	ChatOnWhatsApp: "व्हाट्सऐप पर बात करें",
}

func langFor(code string) strset {
	if code == "hi" {
		return langHI
	}
	return langEN
}

// servicesHeading varies by category: eateries get a menu, kiranas a rate list.
func servicesHeading(category, lang string) string {
	switch category {
	case "restaurant", "bakery":
		if lang == "hi" {
			return "मेनू"
		}
		return "Menu"
	case "kirana":
		if lang == "hi" {
			return "रेट लिस्ट"
		}
		return "Price List"
	default:
		if lang == "hi" {
			return "सेवाएँ"
		}
		return "Services"
	}
}

// CategoryLabel is the human-readable category name.
func CategoryLabel(category, lang string) string {
	en := map[string]string{
		"salon": "Salon", "restaurant": "Restaurant", "kirana": "Kirana Store",
		"coaching": "Coaching Classes", "clinic": "Clinic", "boutique": "Boutique",
		"gym": "Gym", "bakery": "Bakery",
	}
	hi := map[string]string{
		"salon": "सैलून", "restaurant": "रेस्टोरेंट", "kirana": "किराना स्टोर",
		"coaching": "कोचिंग क्लासेस", "clinic": "क्लिनिक", "boutique": "बुटीक",
		"gym": "जिम", "bakery": "बेकरी",
	}
	m := en
	if lang == "hi" {
		m = hi
	}
	if v, ok := m[category]; ok {
		return v
	}
	return m["salon"]
}

var taglinesEN = map[string]string{
	"salon":      "Look sharp. Feel brand new.",
	"restaurant": "Ghar jaisa khana, every single day.",
	"kirana":     "Your daily needs, two minutes away.",
	"coaching":   "Strong basics. Better marks.",
	"clinic":     "Caring for your family, since day one.",
	"boutique":   "Stitched for you, and only you.",
	"gym":        "Show up. We'll handle the rest.",
	"bakery":     "Fresh from the oven, every morning.",
}

var taglinesHI = map[string]string{
	"salon":      "निखार ऐसा, कि हर कोई पूछे — कहाँ से?",
	"restaurant": "घर जैसा खाना, हर रोज़।",
	"kirana":     "रोज़मर्रा का सामान, बस दो मिनट दूर।",
	"coaching":   "मज़बूत बुनियाद, बेहतर नतीजे।",
	"clinic":     "आपके परिवार की सेहत, हमारी ज़िम्मेदारी।",
	"boutique":   "सिलाई ऐसी, जो सिर्फ़ आपके लिए बनी हो।",
	"gym":        "बस आइए, बाकी हम संभाल लेंगे।",
	"bakery":     "हर सुबह, ताज़ा ओवन से।",
}

// Tagline returns the category tagline in the requested language.
func Tagline(category, lang string) string {
	m := taglinesEN
	if lang == "hi" {
		m = taglinesHI
	}
	if v, ok := m[category]; ok {
		return v
	}
	return m["salon"]
}

var aboutEN = map[string]string{
	"salon":      "{name} has been making {city} look its best with skilled stylists, honest prices and a chair that always has time for you. Walk in for a quick trim or settle in for a full makeover — you'll leave lighter, brighter and camera-ready.",
	"restaurant": "At {name}, the tadka is fresh, the rotis are hot and the welcome is warm. From weekday thalis to weekend family feasts, we serve {city} the kind of food that tastes like home — only with someone else doing the dishes.",
	"kirana":     "{name} is the shop your street depends on. From morning milk to midnight Maggi, we keep {city}'s kitchens stocked — and home delivery is just one WhatsApp message away.",
	"coaching":   "{name} believes marks follow understanding. Small batches, patient teachers and weekly practice tests help students across {city} build strong basics — and the confidence to walk into any exam hall smiling.",
	"clinic":     "{name} offers unhurried, honest care to families across {city}. Clear explanations, sensible prescriptions and appointments that respect your time — healthcare the way it should be.",
	"boutique":   "Every outfit at {name} begins with a conversation and a measuring tape. From festival blouses to bridal lehengas, we stitch {city}'s occasions into clothes that fit like they were made for you — because they were.",
	"gym":        "{name} is where {city} shows up for itself. Clean equipment, trainers who actually train, and plans that fit real life — no intimidation, just steady progress you can see.",
	"bakery":     "The ovens at {name} wake up before {city} does. Fresh pav, warm cookies and celebration cakes made to order — baked in small batches, sold the same day, gone by evening.",
}

var aboutHI = map[string]string{
	"salon":      "{name} में हर कुर्सी पर आपके लिए वक़्त है। हुनरमंद स्टाइलिस्ट, वाजिब दाम और वह निखार जो {city} में सबकी नज़र खींच ले — छोटे ट्रिम से लेकर पूरे मेकओवर तक।",
	"restaurant": "{name} में तड़का ताज़ा है, रोटियाँ गरम हैं और स्वागत दिल से है। हफ़्ते की थाली हो या वीकेंड की पारिवारिक दावत — {city} को हम वही खाना खिलाते हैं जो घर की याद दिला दे।",
	"kirana":     "{name} वह दुकान है जिस पर पूरा मोहल्ला भरोसा करता है। सुबह के दूध से लेकर रात की मैगी तक — {city} की रसोई का सारा सामान, बस एक व्हाट्सऐप संदेश की दूरी पर।",
	"coaching":   "{name} का मानना है कि समझ बनेगी तो नंबर खुद आएँगे। छोटे बैच, धैर्यवान शिक्षक और हर हफ़्ते टेस्ट — {city} के विद्यार्थियों की बुनियाद मज़बूत, आत्मविश्वास और भी मज़बूत।",
	"clinic":     "{name} में इलाज बिना जल्दबाज़ी के होता है। साफ़ बात, सही दवा और वक़्त की कदर — {city} के परिवारों के लिए भरोसेमंद देखभाल, जैसी होनी चाहिए।",
	"boutique":   "{name} में हर पोशाक की शुरुआत बातचीत और नाप से होती है। त्योहार का ब्लाउज़ हो या शादी का लहंगा — {city} के हर मौके के लिए ऐसी सिलाई जो सिर्फ़ आप पर जँचे।",
	"gym":        "{name} वह जगह है जहाँ {city} खुद के लिए वक़्त निकालता है। साफ़ मशीनें, सच में ट्रेन करने वाले ट्रेनर और असल ज़िंदगी में फ़िट होने वाले प्लान।",
	"bakery":     "{name} के ओवन {city} से पहले जागते हैं। ताज़ा पाव, गरम कुकीज़ और ऑर्डर पर बने सेलिब्रेशन केक — छोटे बैच में बनते हैं और शाम तक बिक जाते हैं।",
}

// AutoAbout writes a warm intro when the owner leaves the about box empty.
func AutoAbout(name, city, category, lang string) string {
	m := aboutEN
	if lang == "hi" {
		m = aboutHI
	}
	t, ok := m[category]
	if !ok {
		t = m["salon"]
	}
	if city == "" {
		city = "your city"
		if lang == "hi" {
			city = "आपके शहर"
		}
	}
	return strings.NewReplacer("{name}", name, "{city}", city).Replace(t)
}
