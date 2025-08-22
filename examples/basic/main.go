package main

import (
	"fmt"
	"log/slog"
	"os"
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

func (this *ProductsState) getProduct(code string) *Product {
	for i := range this.products {
		cp := &this.products[i]
		if cp.Code == code {
			return cp
		}
	}
	return nil
}

func (this *ProductsState) Process(inbox <-chan actor.Message) {
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
			cp := this.getProduct(p.Code)
			if cp != nil {
				cp.Quantity += p.Quantity
				fmt.Println("AddQuantityProductMsg update with success")
			} else {
				fmt.Println("AddQuantityProductMsg update fail")
			}

		case RemoveQuantityProductMsg:
			fmt.Println("RemoveQuantityProductMsg")
			p := paylaod.Product
			cp := this.getProduct(p.Code)
			if cp != nil {
				cp.Quantity -= p.Quantity
				fmt.Println("RemoveQuantityProductMsg update with success")
			} else {
				fmt.Println("RemoveQuantityProductMsg update fail")
			}

		default:
			slog.Warn("this msg is unknown", slog.String("msg", msg.String()))
		}
		slog.Info("num of products", slog.Int("total", len(this.products)))
		slog.Info("num of pieces of first", slog.Int("num", this.products[0].Quantity))
	}
}

func (this *ProductsState) Shutdown() {
	this.products = make([]Product, 0)
	slog.Info("all product cleaned")
}

type AddProductMsg struct {
	Product Product
}

type AddQuantityProductMsg struct {
	Product Product
}

type RemoveQuantityProductMsg struct {
	Product Product
}

type CheckStoreRefillProductMsg struct{}
type ReturnStoreRefillReportMsg struct {
	Products []Product
}

func NewProductState() *ProductsState {
	initState := ProductsState{
		products: make([]Product, 0),
	}
	return &initState
}

func main() {
	productActor, err := actor.NewActor(
		actor.NewAddress("local", "product"),
		NewProductState(),
	)
	if err != nil {
		os.Exit(1)
	}

	actor.RegisterActor(&productActor)
	msg := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		AddProductMsg{Product{Code: "ciao", Quantity: 4}},
		nil,
	)
	msg2 := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		AddQuantityProductMsg{Product{Code: "ciao", Quantity: 14}},
		nil,
	)
	msg3 := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		RemoveQuantityProductMsg{Product{Code: "ciao", Quantity: 14}},
		nil,
	)
	msg4 := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		CheckStoreRefillProductMsg{},
		nil,
	)
	actor.SendMessage(msg)
	actor.SendMessage(msg2)
	actor.SendMessage(msg3)
	actor.SendMessage(msg4)

	<-time.After(2 * time.Second)
	fmt.Println()
}
