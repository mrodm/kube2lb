/*
Copyright 2016 Tuenti Technologies S.L. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var apiserver, kubecfg, domain, config, template, notify string
	flag.StringVar(&apiserver, "apiserver", "", "Kubernetes API server URL")
	flag.StringVar(&kubecfg, "kubecfg", "", "Path to kubernetes client configuration (Optional)")
	flag.StringVar(&domain, "domain", "local", "DNS domain for the cluster")
	flag.StringVar(&config, "config", "", "Configuration path to generate")
	flag.StringVar(&template, "template", "", "Configuration source template")
	flag.StringVar(&notify, "notify", "", "Kubernetes API server URL")
	flag.Parse()

	if _, err := os.Stat(template); err != nil {
		log.Fatalf("Template not defined or doesn't exist")
	}

	if notify == "" {
		log.Fatalf("Notifier cannot be empty")
	}

	if f, err := os.OpenFile(config, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		log.Fatalf("Cannot open configuration file to write: %v", err)
	} else {
		f.Close()
	}

	notifier, err := NewNotifier(notify)
	if err != nil {
		log.Fatalf("Couldn't initialize notifier: %s", err)
	}

	client, err := NewKubernetesClient(kubecfg, apiserver, domain)
	if err != nil {
		log.Fatalf("Couldn't connect with Kubernetes API server: %s", err)
	}

	client.AddTemplate(NewTemplate(template, config))
	client.AddNotifier(notifier)

	if err := client.Watch(); err != nil {
		log.Fatalf("Couldn't watch Kubernetes API server: %s", err)
	}
}
