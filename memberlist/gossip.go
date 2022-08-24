package memberlist

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
	"net/http"
	"time"
)

type GossipNode struct {
	memberlist *memberlist.Memberlist
	httpClient *http.Client
}

func NewGossipNode(bindPort int, leader string, appPort string) GossipNode {
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

func (node GossipNode) GracefullyLeave() {
	log.Println("gossip service: graceful exit initiated...")
	if err := node.memberlist.Leave(time.Second * 5); err != nil {
		log.Println("gossip service failed to exit gracefully", err)
	}
	log.Println("gossip service: gracefully exited...")
}

func (node *GossipNode) HandleMessage(msg []byte) {
	for _, member := range node.memberlist.Members() {
		fmt.Printf("List Member %s, %s, %d\n", member.Name, member.Addr, member.Port)

		if member == node.memberlist.LocalNode() {
			continue
		}

		url := fmt.Sprintf("http://%s:%s/gossip/api/publish", member.Addr, string(member.Meta))
		resp, err := node.httpClient.Post(url, "application/json", bytes.NewBuffer(msg))

		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			continue
		}
		break
	}
}
