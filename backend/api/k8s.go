package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type namespaceView struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Age    string `json:"age"`
}

type podView struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Phase     string `json:"phase"`
	Node      string `json:"node"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
}

func registerK8sRoutes(r *gin.RouterGroup, client *kubernetes.Clientset) {
	r.GET("/namespaces", func(c *gin.Context) {
		list, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out := make([]namespaceView, 0, len(list.Items))
		now := time.Now()
		for _, ns := range list.Items {
			age := ""
			if ns.CreationTimestamp.Time.Unix() > 0 {
				age = now.Sub(ns.CreationTimestamp.Time).Truncate(time.Second).String()
			}
			out = append(out, namespaceView{
				Name:   ns.Name,
				Status: string(ns.Status.Phase),
				Age:    age,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": out, "total": len(out)})
	})

	r.GET("/pods", func(c *gin.Context) {
		ns := c.Query("namespace")
		list, err := client.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out := make([]podView, 0, len(list.Items))
		now := time.Now()
		for _, p := range list.Items {
			var restarts int32
			for _, cs := range p.Status.ContainerStatuses {
				restarts += cs.RestartCount
			}
			age := ""
			if p.CreationTimestamp.Time.Unix() > 0 {
				age = now.Sub(p.CreationTimestamp.Time).Truncate(time.Second).String()
			}
			out = append(out, podView{
				Name:      p.Name,
				Namespace: p.Namespace,
				Phase:     string(p.Status.Phase),
				Node:      p.Spec.NodeName,
				Restarts:  restarts,
				Age:       age,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"items":     out,
			"total":     len(out),
			"namespace": ns,
		})
	})
}
