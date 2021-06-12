package main

import (
	"context"
	"flag"

	log "github.com/gogap/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//openssl req -new -key apiserver-host.key -subj "/CN=kube-apiserver," -out apiserver-host.csr
// openssl x509 -req -in apiserver-host.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial -out apiserver-host.key.crt -days 365 -extfile apiserver-host.ext
type Client struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

const (
	defaultpodname    = "testpod"
	defaultnamespace  = "prj-install"
	defaultkubeconfig = "kubeconfig"
)

func main() {
	//kubelet.kubeconfig  是文件对应地址
	kubeconfig := flag.String("kubeconfig", defaultkubeconfig, "(optional) absolute path to the kubeconfig file")
	namespace := flag.String("namespace", defaultnamespace, "(optional) which namespace pod create")
	flag.Parse()

	// 解析到config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// 创建连接
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	c := Client{
		Clientset: clientset,
		Namespace: *namespace,
	}

	// List pod
	if pods, err := c.ListPod(c.Namespace); err != nil {
		log.Error("list pod error: ", err)
		return
	} else {
		log.Info(c.Namespace, " have ", len(pods.Items), " pods")
	}
	// create pod
	if err := c.CreatePod(v1.Pod{}); err != nil {
		log.Error("create pod error: ", err)
		return
	} else {
		log.Info("succeed to create pod")
	}

	// delete pod
	if err := c.DeletePod(c.Namespace, defaultpodname); err != nil {
		log.Error("delete pod error: ", err)

	} else {
		log.Info("succeed to delete pod")
	}
}

func (c Client) ListPod(namespace string) (*v1.PodList, error) {
	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error("list pod error: ", err)
		return nil, err
	}
	return pods, nil
}

func (c Client) CreatePod(podtpl v1.Pod) error {
	pod := &v1.Pod{
		TypeMeta:   podtpl.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: defaultpodname, Namespace: defaultnamespace, Labels: map[string]string{"name": defaultpodname}},
		Spec: v1.PodSpec{
			Volumes:        podtpl.Spec.Volumes,
			InitContainers: podtpl.Spec.InitContainers,
			Containers: []v1.Container{
				{
					Name:  "lee-test",
					Image: "nginx:latest",
					//ImagePullPolicy: podtpl.Spec.Containers[0].ImagePullPolicy,
					//Command:         []string{""},
					//Args:            []string{""},
					//Env: podtpl.Spec.Containers[0].Env,
					//Ports:           podtpl.Spec.Containers[0].Ports,
					//Resources: podtpl.Spec.Containers[0].Resources,
					//VolumeMounts:    podtpl.Spec.Containers[0].VolumeMounts,
					//SecurityContext: podtpl.Spec.Containers[0].SecurityContext,
					//Stdin:           podtpl.Spec.Containers[0].Stdin,
				},
			},
			TerminationGracePeriodSeconds: podtpl.Spec.TerminationGracePeriodSeconds,
			Affinity:                      &v1.Affinity{},
		},
		Status: v1.PodStatus{},
	}
	op := metav1.CreateOptions{}

	_, err := c.Clientset.CoreV1().Pods(c.Namespace).Create(context.Background(), pod, op)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) DeletePod(namespace, podname string) error {
	if err := c.Clientset.CoreV1().Pods(c.Namespace).Delete(context.Background(), podname, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}
