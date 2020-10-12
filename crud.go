package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listDemo(ns string) {
	ctx := context.Background()
	opts := v1.ListOptions{}
	delpoys, err := classicalClientSet.AppsV1().Deployments(ns).List(ctx, opts)
	if err != nil {
		panic(err)
	}
	for _, d := range delpoys.Items {
		fmt.Println(d.Name)
		for _, c := range d.Spec.Template.Spec.Containers {
			fmt.Println("   ", c.Name, c.Image, c.ImagePullPolicy)
		}
		fmt.Println("=========")
	}
}

func updateDeploy() {

}
