// Mgmt
// Copyright (C) 2013-2015+ James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	//etcd_context "github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"

	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type collectorTypeConfig struct {
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern"` // XXX: Not Implemented
}

type vertexConfig struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

type edgeConfig struct {
	Name string       `yaml:"name"`
	From vertexConfig `yaml:"from"`
	To   vertexConfig `yaml:"to"`
}

type graphConfig struct {
	Graph string `yaml:"graph"`
	Types struct {
		Noop    []NoopType    `yaml:"noop"`
		File    []FileType    `yaml:"file"`
		Service []ServiceType `yaml:"service"`
	} `yaml:"types"`
	Collector []collectorTypeConfig `yaml:"collect"`
	Edges     []edgeConfig          `yaml:"edges"`
	Comment   string                `yaml:"comment"`
}

func (c *graphConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	if c.Graph == "" {
		return errors.New("Graph config: invalid `graph`")
	}
	return nil
}

func UpdateGraphFromConfig(filename, hostname string, g *Graph, kapi etcd.KeysAPI) bool {

	var NoopMap map[string]*Vertex = make(map[string]*Vertex)
	var FileMap map[string]*Vertex = make(map[string]*Vertex)
	var ServiceMap map[string]*Vertex = make(map[string]*Vertex)

	var lookup map[string]map[string]*Vertex = make(map[string]map[string]*Vertex)
	lookup["noop"] = NoopMap
	lookup["file"] = FileMap
	lookup["service"] = ServiceMap

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return false
	}

	var config graphConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
		return false
	}
	//fmt.Printf("%+v\n", config) // debug

	g.SetName(config.Graph) // set graph name

	var keep []*Vertex // list of vertex which are the same in new graph

	for _, t := range config.Types.Noop {
		obj := NewNoopType(t.Name)
		v := g.GetVertexMatch(obj)
		if v == nil { // no match found
			v = NewVertex(obj)
			g.AddVertex(v) // call standalone in case not part of an edge
		}
		NoopMap[obj.Name] = v  // used for constructing edges
		keep = append(keep, v) // append
	}

	for _, t := range config.Types.File {
		// XXX: should we export based on a @@ prefix, or a metaparam
		// like exported => true || exported => (host pattern)||(other pattern?)
		if strings.HasPrefix(t.Name, "@@") { // exported resource
			// add to etcd storage...
			t.Name = t.Name[2:] //slice off @@
			if !EtcdPut(kapi, hostname, t.Name, "file", t) {
				log.Printf("Problem exporting file resource %v.", t.Name)
				continue
			}
		} else {
			obj := NewFileType(t.Name, t.Path, t.Content, t.State)
			v := g.GetVertexMatch(obj)
			if v == nil { // no match found
				v = NewVertex(obj)
				g.AddVertex(v) // call standalone in case not part of an edge
			}
			FileMap[obj.Name] = v  // used for constructing edges
			keep = append(keep, v) // append
		}
	}

	for _, t := range config.Types.Service {
		obj := NewServiceType(t.Name, t.State, t.Startup)
		v := g.GetVertexMatch(obj)
		if v == nil { // no match found
			v = NewVertex(obj)
			g.AddVertex(v) // call standalone in case not part of an edge
		}
		ServiceMap[obj.Name] = v // used for constructing edges
		keep = append(keep, v)   // append
	}

	// lookup from etcd graph
	// do all the graph look ups in one single step, so that if the etcd
	// database changes, we don't have a partial state of affairs...
	nodes, ok := EtcdGet(kapi)
	if ok {
		for _, t := range config.Collector {
			// XXX: use t.Type and optionally t.Pattern to collect from etcd storage
			log.Printf("Collect: %v(%v)", t.Type, t.Pattern)

			for _, x := range EtcdGetProcess(nodes, "file") {
				var obj *FileType
				if B64ToObj(x, &obj) != true {
					log.Printf("File: %v error!", x)
					continue
				}
				log.Printf("File: %v found!", obj.GetName())

				// XXX: similar to file add code:
				v := g.GetVertexMatch(obj)
				if v == nil { // no match found
					obj.Init() // initialize go channels or things won't work!!!
					v = NewVertex(obj)
					g.AddVertex(v) // call standalone in case not part of an edge
				}
				FileMap[obj.GetName()] = v // used for constructing edges
				keep = append(keep, v)     // append

			}

		}
	}

	// get rid of any vertices we shouldn't "keep" (that aren't in new graph)
	for _, v := range g.GetVertices() {
		if !HasVertex(v, keep) {
			// wait for exit before starting new graph!
			v.Type.SendEvent(eventExit, true)
			g.DeleteVertex(v)
		}
	}

	for _, e := range config.Edges {
		g.AddEdge(lookup[e.From.Type][e.From.Name], lookup[e.To.Type][e.To.Name], NewEdge(e.Name))
	}

	return true
}
