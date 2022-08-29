package memberlist

import (
	"bytes"
	"dist-app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/memberlist"
)

type GossipNode struct {
	memberlist *memberlist.Memberlist
	httpClient *http.Client
}

func NewGossipNode(bindPort int, leader, appPort string) GossipNode {
	list := CreateInitMemberList(bindPort)
	node := GossipNode{
		memberlist: list,
		httpClient: http.DefaultClient,
	}

	fmt.Println(list.LocalNode(), list.LocalNode().Addr, list.LocalNode().Port)

	nodes := make([]string, 0)
	if leader != "" {
		nodes = append(nodes, leader)
	}

	n, err := list.Join(nodes)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
	fmt.Printf("number of nodes contacted: %d\n", n)

	node.memberlist.LocalNode().Meta = []byte(appPort)
	node.clusterinfo()

	return node
}

func CreateInitMemberList(bindPort int) *memberlist.Memberlist {
	config := memberlist.DefaultLocalConfig()
	config.BindPort = bindPort
	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	return list
}

func (node GossipNode) clusterinfo() {
	fmt.Printf("number of nodes in the list: %d\n", len(node.memberlist.Members()))
	for i, node := range node.memberlist.Members() {
		fmt.Printf("List Member %d: %s, %s, %d\n", i, node.Name, node.Addr, node.Port)
		fmt.Printf("App service port: %s\n", string(node.Meta))
	}
	fmt.Printf("Cluster HealthScore: %d\n", node.memberlist.GetHealthScore())
}

func (node GossipNode) GracefullyLeave(timeout time.Duration) {
	log.Println("gossip service: graceful exit initiated...")
	if err := node.memberlist.Leave(timeout); err != nil {
		log.Println("gossip service failed to exit gracefully", err)
	}
	log.Println("gossip service: gracefully exited...")
}

func (node *GossipNode) HandleMessage(msg *model.PublishEvent) {
	dataInBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshaling failed: ", err)
	}

	for _, member := range node.memberlist.Members() {
		fmt.Printf("List Member %s, %s, %d\n", member.Name, member.Addr, member.Port)

		if member == node.memberlist.LocalNode() {
			continue
		}

		url := fmt.Sprintf("http://%s:%s/gossip/api/publish", member.Addr, string(member.Meta))
		resp, err := node.httpClient.Post(url, "application/json", bytes.NewBuffer(dataInBytes)) //nolint
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		resp.Body.Close()
		break
	}
}
