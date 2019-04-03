package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type node struct {
	id     string
	parent string
	root   string
	height int
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /nodes/nodePath to get node information. Example: /nodes/a/c or /nodes/ to get root")
		fmt.Fprintf(w, "POST /nodes/nodePath to create nodes. Example: /nodes/f/o/o/b/a/r")
		fmt.Fprintf(w, "PUT /nodes/nodePath?newParent=newParentNodePath to change parent node. Example: /nodes/a/c?newParent=b")
	})

	http.HandleFunc("/nodes/", func(w http.ResponseWriter, r *http.Request) {
		nodePath := r.URL.Path[6:]

		switch r.Method {
		case "GET":
			getNodeInfo(w, nodePath)
		case "POST":
			createNode(w, nodePath)
		case "PUT":
			newParentPath := r.URL.Query().Get("newParent")
			moveParent(w, nodePath, newParentPath)
		default:
			http.NotFound(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(getPort(), nil))
}

func createNode(w http.ResponseWriter, nodePath string) {
	err := os.MkdirAll(getStoragePath(nodePath), 0777)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to create node: %s\n\n", nodePath)
	} else {
		fmt.Fprintf(w, "Created node if it didn't already exist: %s\n\n", nodePath)
	}

}

func getNodeInfo(w http.ResponseWriter, nodePath string) {

	nodeNames, err := getChildren(nodePath)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Node not found: %s\n\n", nodePath)
	} else if len(nodeNames) == 0 {
		fmt.Fprintf(w, "Node %s has no children\n\n", nodePath)
	} else {

		fmt.Fprintf(w, "Child nodes of %s:\n\n", nodePath)
		pathParts := getPathParts(nodePath)

		for _, nodeName := range nodeNames {
			childNode := node{
				id:     path.Join(nodePath, nodeName),
				parent: nodePath,
				root:   "/",
				height: len(pathParts) + 1}

			fmt.Fprintf(w, "%+v\n", childNode)
		}

	}

}

func moveParent(w http.ResponseWriter, nodePath, newParentPath string) {
	if newParentPath != "" {
		_, nodeName := path.Split(nodePath)
		newPath := path.Join(newParentPath, nodeName)
		err := os.Rename(getStoragePath(nodePath), getStoragePath(newPath))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unable to change parent for node: %s\n\n", nodePath)
		} else {
			fmt.Fprintf(w, "Moved %s to new parent: %s\n\n", nodePath, newParentPath)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No new parent given\n\n")
	}
}

func getChildren(nodePath string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(getStoragePath(nodePath))

	if err != nil {
		return nil, err
	}

	var nodeNames []string

	for _, file := range fileInfos {
		if file.IsDir() {
			nodeNames = append(nodeNames, file.Name())
		}
	}

	return nodeNames, nil
}

func getPathParts(nodePath string) []string {
	var pathParts []string
	for _, part := range strings.Split(nodePath, "/") {
		if part != "" {
			pathParts = append(pathParts, part)
		}
	}

	return pathParts
}

func getStoragePath(nodePath string) string {
	rootPath := getEnv("STORAGE_ROOT", "/storage")
	return path.Join(rootPath, nodePath)
}

func getPort() string {
	return fmt.Sprintf(":%s", getEnv("PORT", "80"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
