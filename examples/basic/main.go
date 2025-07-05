package main

import (
	"errors"
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

func (this *ProductsState) ProcessSync(msg actor.Message) (actor.Message, error) {
	switch msg.Body.(type) {
	case CheckStoreRefillProductMsg:
		ps := make([]Product, 0)
		for _, p := range this.products {
			if p.Quantity < 2 {
				slog.Warn("product on adding ", slog.Any("product", p))
				ps = append(ps, p)
			}
		}
		slog.Warn("products check returned", slog.Any("products", ps))
		rm := actor.NewMessage(
			msg.To,
			msg.From,
			ReturnStoreRefillReportMsg{
				Products: ps,
			},
		)
		slog.Warn("processed in sync and now return", slog.String("msg", msg.String()))
		return rm, nil
	default:
		slog.Warn("this msg is unknown for sync processing", slog.String("msg", msg.String()))
	}
	return actor.EmptyMessage(), errors.New("no message to process")
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
	msg := actor.Message{
		From: actor.NewAddress("local", "product"),
		To:   actor.NewAddress("local", "product"),
		Body: AddProductMsg{Product{Code: "ciao", Quantity: 4}},
	}
	msg2 := actor.Message{
		From: actor.NewAddress("local", "product"),
		To:   actor.NewAddress("local", "product"),
		Body: AddQuantityProductMsg{Product{Code: "ciao", Quantity: 14}},
	}
	msg3 := actor.Message{
		From: actor.NewAddress("local", "product"),
		To:   actor.NewAddress("local", "product"),
		Body: RemoveQuantityProductMsg{Product{Code: "ciao", Quantity: 14}},
	}
	msg4 := actor.Message{
		From: actor.NewAddress("local", "product"),
		To:   actor.NewAddress("local", "product"),
		Body: CheckStoreRefillProductMsg{},
	}
	actor.DispatchMessage(msg)
	actor.DispatchMessage(msg2)
	actor.DispatchMessage(msg3)
	actor.DispatchMessage(msg3)

	<-time.After(2 * time.Second)

	rm, err := actor.DispatchMessageWithReturn(msg4)
	slog.Info("end", slog.Any("rm", rm.Body.(ReturnStoreRefillReportMsg).Products[0]))
}
