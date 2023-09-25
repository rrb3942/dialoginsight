package osipcollect

import (
	"context"
	"dialoginsight/metrics"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	profilesToExport    map[string]bool
	replicationHints    map[string][]string
	exportProfiles      *metrics.DynamicGauges
	exportValueProfiles *metrics.DynamicGauges
	insightProfiles     *metrics.DynamicGauges
	rpc                 *rpc.Client
	url                 string
	insightPrefix       string
	timeout             time.Duration
	mu                  sync.Mutex
	exportAll           bool
}

func NewClient(url, insightPrefix string, exportProfiles []string, exportAll bool, timeout, idleRemove time.Duration) (client *Client, err error) {
	client = new(Client)
	client.url = url
	client.insightPrefix = insightPrefix + ":"
	client.timeout = timeout

	client.exportProfiles = metrics.NewDynamicGauges(exportNamespace, "dialogs", "Exported dialog profiles", idleRemove)
	client.exportValueProfiles = metrics.NewDynamicGauges(exportNamespace, "dialogs", "Exported dialog profiles with values", idleRemove)
	client.insightProfiles = metrics.NewDynamicGauges(insightNamespace, "dialogs", "Insight dialogs with dynamic labels", idleRemove)

	client.profilesToExport, client.replicationHints = parseSharedTags(exportProfiles)

	client.exportAll = exportAll

	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	client.rpc, err = rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	return client, nil
}
