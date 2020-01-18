If you new user from shodan, you get error:
```bash
2020/01/18 12:23:55 json: cannot unmarshal number into Go value of type shodan.APIInfo
panic: json: cannot unmarshal number into Go value of type shodan.APIInfo


goroutine 1 [running]:
log.Panicln(0xc0000ddca0, 0x1, 0x1)
        /usr/local/go/src/log/log.go:352 +0xac
main.main()
        /home/dreddsa/go/src/github.com/dreddsa5dies/goHackTools/57_shodanAPI/cmd/main.go:21 +0xed
exit status 2
```
because you received FREE API PLAN and you API ansver:
```bash
{"scan_credits": 0, "usage_limits": {"scan_credits": 0, "query_credits": 0, "monitored_ips": 0}, "plan": "oss", "https": false, "unlocked": false, "query_credits": 0, "monitored_ips": 0, "unlocked_left": 0, "telnet": false}
```
```bash
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	Https        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}
```