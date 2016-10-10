package impl

import (
	"bytes"
	"fmt"
	"sync"
	"text/tabwriter"
	"time"
)

type Node struct {
	State         bool
	Name          string
	PuppetVersion string
	LastReport    int64
	Errors        bytes.Buffer
}

type EnvironmentCollection struct {
	sm   sync.Mutex
	Env  map[string]NodeCollection
	Conf *Settings
}

type NodeCollection struct {
	sm    sync.Mutex
	Nodes map[string]Node
}

func (n EnvironmentCollection) NewEnvironmentCollection() EnvironmentCollection {
	n.sm.Lock()
	defer n.sm.Unlock()
	return EnvironmentCollection{Env: make(map[string]NodeCollection)}
}

func (n *EnvironmentCollection) GetEnvironmentNodes(name string) *NodeCollection {
	n.sm.Lock()
	nodes, ok := n.Env[name]
	if !ok {
		nodes = NodeCollection{Nodes: make(map[string]Node)}
		n.Env[name] = nodes
	}
	n.sm.Unlock()
	return &nodes
}

func (n *EnvironmentCollection) RemoveNode(host string) string {
	n.sm.Lock()
	defer n.sm.Unlock()
	for _, nodes := range n.Env {
		_, exists := nodes.Nodes[host]
		if exists {
			delete(nodes.Nodes, host)
			return "OK"
		}
	}
	return "Node: " + host + " not exists"
}

func (n *EnvironmentCollection) ProcessReport(report []byte) bool {
	rep, err := ReportItem{}.FromJson(report)
	if err != nil {
		return false
	}
	//rep = rep.FromJson(report)
	nodes := *(n.GetEnvironmentNodes(rep.Environment))
	nodes.sm.Lock()
	defer nodes.sm.Unlock()
	node, exists := nodes.Nodes[rep.Host]
	if !exists {
		node = Node{Name: rep.Host}
		nodes.Nodes[rep.Host] = node

	}
	node.Name = rep.Host
	//nodes.Nodes[rep.Host] = node
	node.PuppetVersion = rep.PuppetVersion
	node.Errors.Reset()
	// analyzing report
	node.State = rep.Status != "failed" &&
		rep.Metrics.Events.Failed+rep.Metrics.Resources.Failed+rep.Metrics.Resources.FailedToRestart == 0
	for _, logI := range rep.Logs {
		if logI.Level == "err" {
			node.State = false
			node.Errors.WriteString(logI.Message + "\n")
		}
	}
	node.LastReport = time.Now().Unix()
	nodes.Nodes[rep.Host] = node
	return true
}

func (n *EnvironmentCollection) GetInfo() string {
	n.sm.Lock()
	var sb bytes.Buffer
	defer n.sm.Unlock()
	max_diff := int64(n.Conf.ControlTime * 60)
	now := time.Now().Unix()
	var state string
	w := tabwriter.NewWriter(&sb, 1, 0, 1, ' ', tabwriter.TabIndent|tabwriter.Debug)
	fmt.Fprintf(w, " %v\t %v\t %v\t %v\t\n", "Host", "\033[00mStatus\033[0m", "Last report", "Agent version")
	for _, nodes := range n.Env {
		//fmt.Fprintf(w, "environment %v\t\n", env)
		for host, node := range nodes.Nodes {
			if node.State && now-node.LastReport <= max_diff {
				state = "\033[32mOK\033[0m"
			} else {
				state = "\033[31mError\033[0m"
			}

			fmt.Fprintf(w, " %v\t %v\t %v\t %v\t\n", host, state, time.Unix(node.LastReport, 0).Format("02.01.2006 15:04:05"), node.PuppetVersion)
			//fmt.Fprintf(w,"status: %v\n", strings.Join(states, ", ")))
			//fmt.Fprintf(w,"agent version: %v\n", node.PuppetVersion))
			//fmt.Fprintf(w,"last report: %v\n", time.Unix(node.LastReport, 0).Format("02.01.2006 15:04:05")))
		}
	}
	fmt.Fprintln(w)
	w.Flush()
	return string(sb.Bytes())
}

func (n *EnvironmentCollection) ProcessCollectionState(write_errors bool) string {
	n.sm.Lock()
	defer n.sm.Unlock()
	var sb bytes.Buffer
	//control_utime := time.Now().Unix() +
	max_diff := int64(n.Conf.ControlTime * 60)
	now := time.Now().Unix()
	overall_state := true
	for env, nodes := range n.Env {
		for host, node := range nodes.Nodes {
			if node.State && now-node.LastReport <= max_diff {
				continue
			}
			overall_state = false
			if write_errors {
				if !node.State {
					sb.WriteString(fmt.Sprintf("\033[1;33m%v\033[0m: failed to implement manifest\n", host))
					sb.WriteString(fmt.Sprintf("Environment: %v\n", env))
					sb.WriteString(fmt.Sprintf("Agent version: %v\n", node.PuppetVersion))
					sb.WriteString(fmt.Sprintf("Last report: %v\n", time.Unix(node.LastReport, 0)))
					sb.Write(node.Errors.Bytes())
					sb.WriteString("\n")
				}
				if now-node.LastReport > max_diff {
					sb.WriteString(fmt.Sprintf("\033[1;33m%v\033[0m: out of sync, no report since \"%v\"\n", host, time.Unix(node.LastReport, 0)))
				}
			} else {
				sb.WriteString(fmt.Sprintf("%v\n", host))
			}
		}
	}
	if overall_state {
		return "OK"
	} else {
		return string(sb.Bytes())
	}
}
