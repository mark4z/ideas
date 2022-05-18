package main

import (
	"log"
	"testing"
)

func TestGrafana(t *testing.T) {
	log.Printf("Starting UC")
	uc := center{
		pod:       "(user-center-lifecycle|user-center-privilege)",
		nameSpace: "ecommerce-user-center",
		db:        "pd-ecom-uc-common-auroramysql",
		redis: `"pd-ecom-uc-lifecycle-redis-0001-001",
             "pd-ecom-uc-lifecycle-redis-0001-002",
             "pd-ecom-uc-privilege-redis-0001-001",
             "pd-ecom-uc-privilege-redis-0001-002"`,
		es: "test",
	}
	uc.call()
	log.Printf("\n")
	log.Printf("Starting MGC")
	mgc := center{
		pod:       "message-center",
		nameSpace: "ecommerce-message-center",
		db:        "pd-ecom-mgc-common-auroramysql",
		redis: `"pd-ecom-mgc-app-redis-0001-001",
			 "pd-ecom-mgc-app-redis-0001-002"`,
	}
	mgc.call()
}
