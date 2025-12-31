# Tor-Scraper
## Nedir?
Tor ağı üzerinden anonim olarak belirlenen URL’leri tarar, erişilebilen web sayfalarının HTML içeriklerini kaydeder ve tüm işlemleri log dosyası halinde raporlayan bir scraperdır.

Bu projenin amacı anonim bağlantı sağlayarak Dark Web’teki hacker forumlarının ve hacker gruplarının sitelerinin hızlı bir şekilde taranarak hala aktif mi pasif mi olduğunu anlamak , sitenin HTML’ni kaydetmek ve tarama sonuçlarını log dosyasına kaydetmektir.

# Kullanılan Teknolojiler

Bu projede Go version 1.24.3 kullanılmıştır. 

* Standart kütüphanelere ek olarak trafiği yerel Tor Servisi (SOCK5 Proxy) ile yönlendirmek  için sadece “golang.org/x/net/proxy” kütüphanesi kullanılmıştır.
  
Ağ alt yapısı için Tor Service kullanılmaktadır.

# Kullanım 
1- Öncelikle taranacak URL'ler "targets.yaml" dosyasına satır satır eklenir. 

2- Sistemde Tor servisi çalışır duruma getirilir. 

3- Proje konumuna terminal üzerinden girilerek "go run main.go" komutu çalıştırılır. 

4- Program, Tor ağı üzerinden anonim olarak belirtilen URL’lere istek gönderir.

5- Erişilebilen web sayfalarının HTML içerikleri Results/ klasörüne kaydedilir ve tüm işlemler scan_report.log dosyasında raporlanır.
