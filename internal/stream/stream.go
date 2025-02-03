package stream

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const StreamURL = "https://stream.wikimedia.org/v2/stream/recentchange"

type RecentChange struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	User       string `json:"user"`
	Timestamp  int64  `json:"timestamp"`
	Meta       struct {
		Domain string `json:"meta"`
	} `json:"meta"`
}

func FormatTimestamp(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func GetRecentChanges(lang string, maxChanges int) ([]RecentChange, error) {
	changes := make([]RecentChange, 0, maxChanges)

	client := &http.Client{}
	req, err := http.NewRequest("GET", StreamURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonData := line[6:]
			var change RecentChange
			if err := json.Unmarshal([]byte(jsonData), &change); err != nil {
				continue
			}

			serverParts := strings.Split(change.ServerName, ".")
			if len(serverParts) < 2 || serverParts[0] != lang {
				continue
			}

			changes = append(changes, change)

			if len(changes) >= maxChanges {
				break
			}
		}
	}

	return changes, nil
}
