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
	exportProfiles      *metrics.DynamicGauges
	exportValueProfiles *metrics.DynamicGauges
	insightProfiles     *metrics.DynamicGauges
	rpc                 *rpc.Client
	url                 string
	insightPrefix       string
	timeout             time.Duration
	sync.Mutex
	exportAll bool
}

func NewClient(url, insightPrefix string, exportProfiles []string, exportAll bool, timeout, idleRemove time.Duration) (*Client, error) {
	client := new(Client)

	client.url = url
	client.insightPrefix = insightPrefix + ":"
	client.timeout = timeout

	client.exportProfiles = metrics.NewDynamicGauges(exportNamespace, "dialogs", "Exported dialog profiles", idleRemove)
	client.exportValueProfiles = metrics.NewDynamicGauges(exportNamespace, "dialogs", "Exported dialog profiles with values", idleRemove)
	client.insightProfiles = metrics.NewDynamicGauges(insightNamespace, "dialogs", "Insight dialogs with dynamic labels", idleRemove)

	client.profilesToExport = make(map[string]bool)
	for _, v := range exportProfiles {
		client.profilesToExport[v] = true
	}

	client.exportAll = exportAll

	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	if rpcClient, err := rpc.DialContext(ctx, url); err != nil {
		return nil, err
	} else {
		client.rpc = rpcClient
	}

	return client, nil
}
