package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/pix303/actor-lib/pkg/actor"
)

type Product struct {
	Code     string
	Quantity int
}

type ProductsState struct {
	products []Product
}

func (this *ProductsState) Process(inbox chan actor.Message) {
	for {
		msg := <-inbox
		switch paylaod := msg.Body.(type) {
		case AddProductMsg:
			fmt.Println("AddProductMsg")
			p := paylaod.Product
			this.products = append(this.products, p)
		case AddQuantityProductMsg:
			fmt.Println("AddQuantityProductMsg")
			p := paylaod.Product
			for i := range this.products {
				cp := &this.products[i]
				if cp.Code == p.Code {
					cp.Quantity += p.Quantity
				}
			}
		default:
			slog.Warn("this msg is unknown", slog.String("msg", msg.ToString()))
		}
		slog.Info("num of products", slog.Int("total", len(this.products)))
		slog.Info("num of pieces of first", slog.Int("num", this.products[0].Quantity))
	}
}

type AddProductMsg struct {
	Product Product
}

type AddQuantityProductMsg struct {
	Product Product
}

func NewProductState() *ProductsState {
	initState := ProductsState{
		products: make([]Product, 0),
	}
	return &initState
}

func main() {
	productActor := actor.NewActor(
		actor.NewPID("local", "product"),
		NewProductState(),
	)
	productActor.Activate()
	ds := actor.NewActorDispatcher()
	ds.RegisterActor(&productActor)
	msg := actor.Message{
		From: *actor.NewPID("local", "product"),
		To:   *actor.NewPID("local", "product"),
		Body: AddProductMsg{Product{Code: "ciao", Quantity: 4}},
	}
	msg2 := actor.Message{
		From: *actor.NewPID("local", "product"),
		To:   *actor.NewPID("local", "product"),
		Body: AddQuantityProductMsg{Product{Code: "ciao", Quantity: 14}},
	}
	ds.DispatchMessage(msg)
	ds.DispatchMessage(msg2)

	<-time.After(2 * time.Second)
	fmt.Println("end")
}
