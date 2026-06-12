package main

import (
	"github.com/yasinatesim/vela-commerce/api/internal/auth"
	"github.com/yasinatesim/vela-commerce/api/internal/category/models"
	"github.com/yasinatesim/vela-commerce/api/internal/config"
	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/database"
	"github.com/yasinatesim/vela-commerce/api/internal/logger"
	"github.com/yasinatesim/vela-commerce/api/internal/product/models"
	"github.com/yasinatesim/vela-commerce/api/internal/rag"
	"github.com/yasinatesim/vela-commerce/api/internal/user/models"
)

type catSeed struct {
	slug, nameTr, nameEn, icon, desc string
}

type prodSeed struct {
	name, cat, badge, seller, desc string
	price, oldP                    int64
	stock                          int
}

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(db); err != nil {
		panic(err)
	}

	hash, _ := auth.HashPassword("admin12345")
	admin := usermodels.User{Email: "admin@vela.test", PasswordHash: hash, Role: constants.ROLE_ADMIN, FullName: "Vela Admin"}
	db.Where("email = ?", admin.Email).FirstOrCreate(&admin, usermodels.User{Email: admin.Email})

	for _, c := range categories() {
		row := categorymodels.Category{Slug: c.slug, NameTr: c.nameTr, NameEn: c.nameEn, Icon: c.icon, DescTr: c.desc}
		db.Where("slug = ?", c.slug).FirstOrCreate(&row, categorymodels.Category{Slug: c.slug})
	}

	for _, p := range products() {
		slug := rag.Slugify(p.name)
		row := productmodels.Product{
			Slug: slug, TitleTr: p.name, TitleEn: p.name, DescriptionTr: p.desc,
			PriceCents: p.price * 100, OldPriceCents: p.oldP * 100, Stock: p.stock,
			CategorySlug: p.cat, Badge: p.badge, Seller: p.seller, Published: true,
		}
		db.Where("slug = ?", slug).FirstOrCreate(&row, productmodels.Product{Slug: slug})
	}

	log.Info("seed complete", "admin", admin.Email, "categories", len(categories()), "products", len(products()))
}

func categories() []catSeed {
	return []catSeed{
		{"elektronik", "Elektronik", "Electronics", "💻", "Telefon, kulaklık, aksesuar"},
		{"moda", "Moda", "Fashion", "👕", "Giyim, ayakkabı, aksesuar"},
		{"ev", "Ev & Yaşam", "Home & Living", "🏠", "Dekorasyon, mutfak, hobi"},
		{"spor", "Spor", "Sports", "⚽", "Spor aletleri, outdoor, bisiklet"},
		{"dogal-tas-bileklik", "Doğal Taş Bileklik", "Gemstone Bracelet", "💎", "Ametist, akik, kaplan gözü"},
		{"dogal-tas-kolye", "Doğal Taş Kolye", "Gemstone Necklace", "📿", "Yeşim, turkuaz, pembe kuvars"},
		{"dogal-tas-tesbih", "Doğal Taş Tesbih", "Gemstone Prayer Beads", "🪨", "Akik, obsidiyen, yeşim"},
		{"dogal-tas-kupe", "Doğal Taş Küpe", "Gemstone Earrings", "✨", "Sitrin, pembe kuvars, sodalit"},
	}
}

func products() []prodSeed {
	return []prodSeed{
		{"Pro Kablosuz Kulaklık", "elektronik", "%18 İndirim", "Vela Tech", "Aktif gürültü engellemeli, 30 saat pil ömrü, bluetooth 5.3.", 1299, 1599, 45},
		{"Akıllı Saat Pro", "elektronik", "Çok Satan", "Vela Tech", "1.43\" AMOLED ekran, GPS, kalp ritmi, 14 gün pil.", 2499, 2999, 28},
		{"Bluetooth Hoparlör", "elektronik", "", "Vela Tech", "IPX7 su geçirmez, 20W ses, 12 saat oynatma.", 599, 0, 62},
		{"Klasik Deri Cüzdan", "moda", "%30 İndirim", "Vela Moda", "El yapımı hakiki deri, 8 kartlık, RFID korumalı.", 349, 499, 18},
		{"Pamuklu Oversize Tişört", "moda", "Yeni", "Vela Moda", "%100 organik pamuk, relaxed fit, 5 renk.", 199, 0, 120},
		{"Kargo Pantolon", "moda", "", "Vela Moda", "Rahat kesim, cepli tasarım, bel ayarlanabilir.", 549, 699, 35},
		{"Seramik Kahve Fincan Seti", "ev", "", "Vela Home", "6'lı el yapımı seramik fincan, ahşap tepsi hediyeli.", 299, 0, 24},
		{"LED Masa Lambası", "ev", "Fırsat", "Vela Home", "Dokunmatik dimmer, 3 renk sıcaklığı, USB şarjlı.", 249, 329, 40},
		{"Ahşap Baharat Takımı", "ev", "", "Vela Home", "6'lı cam kavanoz, ahşap taşıyıcı, etiket setli.", 189, 0, 55},
		{"Yoga Matı Premium", "spor", "%20 İndirim", "Vela Sport", "6mm TPE, çift taraflı kaymaz, taşıma askılı.", 399, 499, 31},
		{"Direnç Bandı Seti", "spor", "Çok Satan", "Vela Sport", "5 farklı direnç, doğal lateks, çanta hediyeli.", 149, 0, 78},
		{"Katlanabilir Su Şişesi", "spor", "Yeni", "Vela Sport", "750ml Tritan, sızdırmaz kapak, BPA free.", 129, 0, 90},
		{"Ametist Taş Bileklik", "dogal-tas-bileklik", "%25 İndirim", "Vela Doğal", "Doğal ametist taşı, elastik ipli, 8mm boncuk.", 299, 399, 34},
		{"Kaplan Gözü Bileklik", "dogal-tas-bileklik", "", "Vela Doğal", "Doğal kaplan gözü taşı, gümüş ara boncuklu.", 259, 0, 28},
		{"Lav Taşı Bileklik", "dogal-tas-bileklik", "Çok Satan", "Vela Doğal", "Volkanik lav taşı, 10mm, siyah, unisex.", 199, 0, 55},
		{"Sodalit Bileklik", "dogal-tas-bileklik", "Fırsat", "Vela Doğal", "Doğal sodalit taşı, mavi-beyaz damarlı, 8mm.", 279, 349, 20},
		{"Yeşim Taşı Kolye", "dogal-tas-kolye", "Yeni", "Vela Doğal", "Doğal yeşim taşı, gümüş 925 zincir, ayarlanabilir.", 449, 0, 16},
		{"Turkuaz Kolye", "dogal-tas-kolye", "%18 İndirim", "Vela Doğal", "Doğal turkuaz, gümüş 925 kolye ucu, 45cm.", 529, 649, 12},
		{"Pembe Kuvars Kolye", "dogal-tas-kolye", "", "Vela Doğal", "Pembe kuvars taşı, altın kaplama zincir.", 399, 0, 22},
		{"Akik Tesbih", "dogal-tas-tesbih", "", "Vela Doğal", "Doğal akik taşı, 33 boncuk, pirinç imame.", 349, 0, 18},
		{"Obsidiyen Tesbih", "dogal-tas-tesbih", "Fırsat", "Vela Doğal", "Siyah obsidiyen, 33 boncuk, gümüş ara parçalı.", 449, 549, 14},
		{"Yeşim Tesbih", "dogal-tas-tesbih", "Çok Satan", "Vela Doğal", "Doğal yeşim taşı, 33 boncuk, el işçiliği.", 499, 0, 10},
		{"Sitrin Taş Küpe", "dogal-tas-kupe", "", "Vela Doğal", "Doğal sitrin taşı, gümüş 925, iğne uçlu.", 329, 0, 25},
		{"Pembe Kuvars Küpe", "dogal-tas-kupe", "%20 İndirim", "Vela Doğal", "Pembe kuvars, altın kaplama, halka küpe.", 299, 379, 30},
		{"Sodalit Taş Küpe", "dogal-tas-kupe", "", "Vela Doğal", "Sodalit taşı, gümüş 925, damla tasarım.", 359, 0, 17},
	}
}
