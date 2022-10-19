[![Go Report Card](https://goreportcard.com/badge/github.com/dreddsa5dies/goHackTools)](https://goreportcard.com/report/github.com/dreddsa5dies/goHackTools) ![License](https://img.shields.io/badge/License-MIT-blue.svg)  

![IMAGE](img/goHackTools.png)  

## Hacker tools on Go (Golang)  
<details>
  <summary>I used examples from the books & materials</summary>

* "Violent Python" TJ O'Connor;
* "Black Hat Python" Python Programming for Hackers and Pentesters by Justin Seitz;
* "Security with Go" by John Daniel Leon;
* "Python Web Penetration Testing Cookbook" by C.Buchanan, T.Ip, B.May, D.Mound, A.Mabbit;
* [asecuritysite](https://asecuritysite.com/); 
* [Криптография с Python](https://vk.com/doc187366527_464874978?hash=45d8e4c6fd48820484&dl=8e644ab04c8ad6520d);  
* "Black Hat Go" Go Programming For Hackers and Pentesters by Tom Steele, Chris Patten, and Dan Kottmann.  

</details>

<details>
  <summary>List of projects</summary>

- [Перебор паролей в passwd](projects/01_crackUnixPass/)
- [Перебор паролей к архиву](projects/02_crackZipPass/)
- [Перебор паролей SSH](projects/12_sshCrack/)
- [Перебор HTML-форм](projects/36_bruteHtmlForm/)
- [Сканер портов](projects/03_tcpScanner/)
- [Сканер портов через Nmap](projects/04_goNmapScan/)
- [Конкурентное сканирование портов](projects/56_PerformingConcurrentScanning/)
- [Определение IP и адреса](projects/08_geoIp/)
- [Определение IP и адреса II](projects/10_buildGoogleMap/)
- [Поиск устройств в сети](projects/43_netScan/)
- [Поиск сетевых устройств](projects/61_findNetDevs/)
- [Проверка хоста по IP](projects/31_lookupIP/)
- [Получение IP-адреса хоста](projects/32_lookupHost/)
- [Получение MX записей](projects/33_getMXRec/)
- [Получение имен серверов DNS](projects/34_getServName/)
- [Checker ресурса](projects/28_webChecker/)
- [Тест SSH](projects/05_sshGexpectShavac/)
- [Перебор сетевых пакетов](projects/09_packetParser/)
- [Получение DNS записей](projects/58_dnsGetA/)
- [Перебор поддоменов](projects/59_subdomains/)
- [Исследование Sqlite браузера](projects/06_forensicMozillaSQLITE/)
- [Получение данных PDF](projects/07_metaDataPdf/)
- [Определение типа файла (изображения)](projects/17_forensicImage/)
- [Определение типа файла (расширенное)](projects/18_forensicFile/)
- [Получение геотегов изображения](projects/19_getGeoTagPhoto/)
- [Получение cookies](projects/24_getCookies/)
- [Получение информации по MAC-адресу](projects/25_calculateYouByMac/)
- [Получение информации о файлах](projects/26_getFileInfo/)
- [Шифр RSA](projects/21_RSAapp/)
- [Шифр Цезаря](projects/22_CipherOfCaesar/)
- [Шифр ROT-13](projects/53_cipherROT13/)
- [Шифр Виженера](projects/54_Vigenere/)
- [Книжный шифр](projects/55_bookCipher/)
- [Шифровальщик](projects/23_Cryptographer/)
- [Кодирование base-64](projects/44_base64/)
- [Стеганография](projects/29_steganoImgArch/)
- [Поиск данных в изображении и их извлечение](projects/30_detectSteganoImgAndExtractIt/)
- [Получение заголовков HTTP](projects/37_httpHead/)
- [Поиск комментариев на web-страницах](projects/38_findHtmlComm/)
- [Поиск скрытых файлов на web-сервере](projects/39_findFilesOnWebServ/)
- [Подмена User Agent](projects/40_userAgent/)
- [Получение заголовков](projects/41_getHeader/)
- [Grabbing сетевого устройства](projects/42_grabbing/)
- [Поиск изменяемых файлов](projects/45_findWrFiles/)
- [Поиск уязвимых файлов](projects/60_filepathInfoSearcher/)
- [Изменение атрибутов файлов](projects/46_fileTimestamp/)
- [Определение прав доступа к файлам](projects/47_filePerm/)
- [Определение принадлежности файлов](projects/48_fileOwnership/)
- [Получение ссылок для Maltego](projects/50_extrLinkToMaltego/)
- [Получение данных email для Maltego](projects/51_extrEmailToMaltego/)
- [Идентификация альтернативных сайтов путем подмены данных agent](projects/52_findAltSites/)
- [Использование API Shodan](projects/57_shodanAPI/)
- [SSH бот-сеть](projects/13_sshSwarm/)
- [TCP proxy](projects/14_TCPProxy/)
- [SSH туннелирование](projects/15_sshTunnel/)
- [Удаленный шелл](projects/16_shell/)
- [SYN flood](projects/62_synFlood/)
- [Netcat](projects/11_netcat/)
- [Взлом OSPF](projects/20_ospfGetAuth/)
- [Загрузка своих данных в icmp пакет](projects/63_icmpPayload/)
- [Чтение загрузочного сектора](projects/27_readingBootSector/)
- [Генератор случайных цифр (псевдо)](projects/35_CSPRNG/)
- [Создание карт объектов](projects/49_buildPNGmap/)

</details>

## [SCC](https://github.com/boyter/scc)
```
───────────────────────────────────────────────────────────────────────────────
Language                 Files     Lines   Blanks  Comments     Code Complexity
───────────────────────────────────────────────────────────────────────────────
Go                          65      4758      878       269     3611        888
Plain Text                   9     21110        0         0    21110          0
Markdown                     2        20        3         0       17          0
───────────────────────────────────────────────────────────────────────────────
Total                       76     25888      881       269    24738        888
───────────────────────────────────────────────────────────────────────────────
Estimated Cost to Develop (organic) $784,563
Estimated Schedule Effort (organic) 12.54 months
Estimated People Required (organic) 5.56
───────────────────────────────────────────────────────────────────────────────
Processed 262599 bytes, 0.263 megabytes (SI)
───────────────────────────────────────────────────────────────────────────────
```

## Thank you
Thanks to all authors of amazing books on information security

## Predict
In spite of the license, I PREDICT that all the examples here are for reference only, and not for criminal (or malicious) purposes. 

## Packages
Use [Go Modules](https://blog.golang.org/using-go-modules) && install
```bash
sudo apt-get install libpcap-dev 
```

## The code contains comments in Russian

## License
This project is licensed under MIT license. Please read the [LICENSE](https://github.com/dreddsa5dies/goHackTools/tree/master/LICENSE.md) file.  

## Contribute
Welcomes any kind of contribution. Please read the [CONTRIBUTING](https://github.com/dreddsa5dies/goHackTools/tree/master/CONTRIBUTING.md), [ISSUE TEMPLATE](https://github.com/dreddsa5dies/goHackTools/tree/master/ISSUE_TEMPLATE.md) and [CODE_OF_CONDUCT](https://github.com/dreddsa5dies/goHackTools/tree/master/CODE_OF_CONDUCT.md) file.
