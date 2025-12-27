package main

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	etcdClient *clientv3.Client
	isLeader   bool
)

func initEtcd() {
	var err error
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd:", err)
	}

	log.Println("Connected to etcd")
}

func startLeaderElection() {
	go func() {
		for {
			// Create lease
			leaseResp, err := etcdClient.Grant(context.Background(), 10)
			if err != nil {
				log.Println("Lease grant failed:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// Try to acquire leadership
			txn := etcdClient.Txn(context.Background()).
				If(clientv3.Compare(
					clientv3.CreateRevision("/scheduler/leader"),
					"=",
					0,
				)).
				Then(clientv3.OpPut(
					"/scheduler/leader",
					"leader",
					clientv3.WithLease(leaseResp.ID),
				))

			txnResp, err := txn.Commit()
			if err != nil {
				log.Println("Leader election txn failed:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if txnResp.Succeeded {
				isLeader = true
				log.Println("✅ This scheduler is LEADER")

				// Keep lease alive
				ch, kaErr := etcdClient.KeepAlive(context.Background(), leaseResp.ID)
				if kaErr != nil {
					isLeader = false
					log.Println("KeepAlive failed:", kaErr)
					continue
				}

				for range ch {
					// heartbeat loop
				}

				// Lease expired
				isLeader = false
				log.Println("❌ Leadership lost")
			}

			time.Sleep(2 * time.Second)
		}
	}()
}
