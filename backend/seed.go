//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

var (
	firstNames = []string{
		"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Hank",
		"Iris", "Jack", "Karen", "Leo", "Mona", "Nick", "Olivia", "Paul",
		"Quinn", "Rita", "Sam", "Tina", "Uma", "Vic", "Wendy", "Xander",
		"Yara", "Zach", "Amara", "Brett", "Cleo", "Derek", "Elena", "Felix",
		"Gina", "Hugo", "Isla", "Jude", "Kira", "Liam", "Maya", "Noah",
		"Opal", "Petra", "Ravi", "Suki", "Tomas", "Ursula", "Vera", "Wes",
		"Ximena", "Yusuf",
	}

	lastNames = []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller",
		"Davis", "Rodriguez", "Martinez", "Anderson", "Taylor", "Thomas",
		"Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
		"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
		"Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen",
		"Hill", "Flores", "Green", "Adams", "Nelson", "Baker", "Hall",
		"Rivera", "Campbell", "Mitchell", "Carter", "Roberts", "Gomez",
		"Phillips", "Evans", "Turner",
	}

	severities = []string{"low", "medium", "high", "critical"}

	alarmTitles = []string{
		"CPU usage exceeded threshold",
		"Memory utilization critical",
		"Disk space running low",
		"Network latency spike detected",
		"Service response time degraded",
		"Database connection pool exhausted",
		"SSL certificate expiring soon",
		"Unauthorized access attempt",
		"Container restart loop detected",
		"API error rate above normal",
		"Load balancer health check failing",
		"Queue depth exceeding limit",
		"Replica lag increasing",
		"Cache hit ratio dropping",
		"DNS resolution failures",
		"Deployment rollback triggered",
		"Pod eviction detected",
		"Storage IOPS throttled",
		"Backup job failed",
		"Log ingestion pipeline stalled",
		"Firewall rule violation",
		"Rate limit threshold reached",
		"Deadlock detected in database",
		"Garbage collection pause too long",
		"Thread pool saturation",
		"TLS handshake failures",
		"Upstream dependency timeout",
		"Disk I/O latency elevated",
		"Memory leak suspected",
		"Cluster node unreachable",
		"Config drift detected",
		"Secret rotation overdue",
		"Ingress controller error spike",
		"Prometheus scrape failures",
		"Kafka consumer lag growing",
		"Redis eviction rate high",
		"Elasticsearch cluster yellow",
		"gRPC deadline exceeded",
		"WebSocket connection storm",
		"Cron job missed schedule",
	}

	alarmDescriptions = []string{
		"Threshold breached for over 5 minutes, investigation required.",
		"Intermittent issue detected across multiple availability zones.",
		"Triggered by automated monitoring. Manual review recommended.",
		"Correlates with recent deployment. Possible regression.",
		"Affecting production traffic. On-call team notified.",
		"First occurrence in 30 days. May be transient.",
		"Recurring issue — seen 3 times this week.",
		"Impact limited to non-critical path. Monitor closely.",
		"Customer-facing impact confirmed. Escalation in progress.",
		"Auto-remediation attempted but failed. Needs manual intervention.",
	}
)

type resource struct {
	ID string `json:"id"`
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create 50 resources
	fmt.Println("Creating 50 resources...")
	resourceIDs := make([]string, 0, 50)

	for i := 0; i < 50; i++ {
		first := firstNames[i%len(firstNames)]
		last := lastNames[i%len(lastNames)]
		displayName := fmt.Sprintf("%s %s", first, last)
		email := fmt.Sprintf("%s.%s@example.com", lower(first), lower(last))
		isUser := rng.Float64() < 0.7
		thumbnail := fmt.Sprintf("https://i.pravatar.cc/150?u=%s.%s", lower(first), lower(last))

		body, _ := json.Marshal(map[string]any{
			"displayName":      displayName,
			"email":            email,
			"isUserAssociated": isUser,
			"thumbnail":        thumbnail,
		})

		resp, err := http.Post(baseURL+"/resources", "application/json", bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Failed to create resource %d: %v", i+1, err)
		}

		var res resource
		json.NewDecoder(resp.Body).Decode(&res)
		resp.Body.Close()

		resourceIDs = append(resourceIDs, res.ID)

		if (i+1)%10 == 0 {
			fmt.Printf("  ...%d resources created\n", i+1)
		}
	}

	// Create 500 alarms
	fmt.Println("Creating 500 alarms...")

	for i := 0; i < 500; i++ {
		title := alarmTitles[rng.Intn(len(alarmTitles))]
		desc := alarmDescriptions[rng.Intn(len(alarmDescriptions))]
		severity := severities[rng.Intn(len(severities))]

		// Assign 0-3 random resources
		numAssigned := rng.Intn(4)
		assigned := make([]string, 0, numAssigned)
		used := make(map[int]bool)
		for j := 0; j < numAssigned; j++ {
			idx := rng.Intn(len(resourceIDs))
			if !used[idx] {
				used[idx] = true
				assigned = append(assigned, resourceIDs[idx])
			}
		}

		body, _ := json.Marshal(map[string]any{
			"title":             fmt.Sprintf("%s (#%d)", title, i+1),
			"description":       desc,
			"severity":          severity,
			"assignedResources": assigned,
		})

		resp, err := http.Post(baseURL+"/alarms", "application/json", bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Failed to create alarm %d: %v", i+1, err)
		}
		resp.Body.Close()

		if (i+1)%100 == 0 {
			fmt.Printf("  ...%d alarms created\n", i+1)
		}
	}

	fmt.Println("Done! Seeded 50 resources and 500 alarms.")
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
