package api

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

type nodeView struct {
	Name    string            `json:"name"`
	Ready   bool              `json:"ready"`
	Roles   []string          `json:"roles"`
	Version string            `json:"version"`
	OS      string            `json:"os"`
	Age     string            `json:"age"`
	Labels  map[string]string `json:"labels"`
}

type deploymentView struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	Replicas  int32  `json:"replicas"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
}

type scaleRequest struct {
	Replicas *int32 `json:"replicas" binding:"required"`
}

func ageSince(t metav1.Time, now time.Time) string {
	if t.Time.Unix() <= 0 {
		return ""
	}
	return now.Sub(t.Time).Truncate(time.Second).String()
}

func registerK8sRoutes(r *gin.RouterGroup, client *kubernetes.Clientset) {
	r.GET("/namespaces", func(c *gin.Context) {
		list, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		out := make([]namespaceView, 0, len(list.Items))
		for _, ns := range list.Items {
			out = append(out, namespaceView{
				Name:   ns.Name,
				Status: string(ns.Status.Phase),
				Age:    ageSince(ns.CreationTimestamp, now),
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": out, "total": len(out)})
	})

	r.GET("/nodes", func(c *gin.Context) {
		list, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		out := make([]nodeView, 0, len(list.Items))
		for _, n := range list.Items {
			ready := false
			for _, cond := range n.Status.Conditions {
				if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
					ready = true
					break
				}
			}
			roles := make([]string, 0)
			for k := range n.Labels {
				const prefix = "node-role.kubernetes.io/"
				if strings.HasPrefix(k, prefix) {
					role := strings.TrimPrefix(k, prefix)
					if role == "" {
						role = "worker"
					}
					roles = append(roles, role)
				}
			}
			if len(roles) == 0 {
				roles = append(roles, "worker")
			}
			out = append(out, nodeView{
				Name:    n.Name,
				Ready:   ready,
				Roles:   roles,
				Version: n.Status.NodeInfo.KubeletVersion,
				OS:      n.Status.NodeInfo.OSImage,
				Age:     ageSince(n.CreationTimestamp, now),
				Labels:  n.Labels,
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
		now := time.Now()
		out := make([]podView, 0, len(list.Items))
		for _, p := range list.Items {
			var restarts int32
			for _, cs := range p.Status.ContainerStatuses {
				restarts += cs.RestartCount
			}
			out = append(out, podView{
				Name:      p.Name,
				Namespace: p.Namespace,
				Phase:     string(p.Status.Phase),
				Node:      p.Spec.NodeName,
				Restarts:  restarts,
				Age:       ageSince(p.CreationTimestamp, now),
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": out, "total": len(out), "namespace": ns})
	})

	r.GET("/pods/:namespace/:name/logs", func(c *gin.Context) {
		ns := c.Param("namespace")
		name := c.Param("name")
		tail := int64(200)
		if v := c.Query("tail"); v != "" {
			if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
				tail = n
			}
		}
		opts := &corev1.PodLogOptions{TailLines: &tail}
		if container := c.Query("container"); container != "" {
			opts.Container = container
		}
		stream, err := client.CoreV1().Pods(ns).GetLogs(name, opts).Stream(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stream.Close()
		data, err := io.ReadAll(stream)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"namespace": ns,
			"name":      name,
			"tail":      tail,
			"logs":      string(data),
		})
	})

	r.POST("/pods/:namespace/:name/restart", func(c *gin.Context) {
		ns := c.Param("namespace")
		name := c.Param("name")
		err := client.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "pod deleted; controller will recreate if managed",
			"namespace": ns,
			"name":      name,
		})
	})

	r.GET("/deployments", func(c *gin.Context) {
		ns := c.Query("namespace")
		list, err := client.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		out := make([]deploymentView, 0, len(list.Items))
		for _, d := range list.Items {
			var replicas int32
			if d.Spec.Replicas != nil {
				replicas = *d.Spec.Replicas
			}
			out = append(out, deploymentView{
				Name:      d.Name,
				Namespace: d.Namespace,
				Ready:     strconv.Itoa(int(d.Status.ReadyReplicas)) + "/" + strconv.Itoa(int(replicas)),
				Replicas:  replicas,
				Available: d.Status.AvailableReplicas,
				Age:       ageSince(d.CreationTimestamp, now),
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": out, "total": len(out), "namespace": ns})
	})

	r.POST("/deployments/:namespace/:name/scale", func(c *gin.Context) {
		ns := c.Param("namespace")
		name := c.Param("name")
		var req scaleRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Replicas == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "replicas is required (integer >= 0)"})
			return
		}
		if *req.Replicas < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "replicas must be >= 0"})
			return
		}
		patch := []byte(`{"spec":{"replicas":` + strconv.Itoa(int(*req.Replicas)) + `}}`)
		_, err := client.AppsV1().Deployments(ns).Patch(
			context.Background(),
			name,
			types.StrategicMergePatchType,
			patch,
			metav1.PatchOptions{},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"namespace": ns,
			"name":      name,
			"replicas":  *req.Replicas,
			"message":   "scaled",
		})
	})
}
