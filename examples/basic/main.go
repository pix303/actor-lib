package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/pix303/cinecity/pkg/actor"
)

// Product represents a product in the inventory.
type Product struct {
	Code     string
	Quantity int
}

// ProductsState represents the state of the products inventory.
type ProductsState struct {
	products []Product
}

func NewProductState() *ProductsState {
	initState := ProductsState{
		products: make([]Product, 0),
	}
	return &initState
}

func (state *ProductsState) getProduct(code string) *Product {
	for i := range state.products {
		cp := &state.products[i]
		if cp.Code == code {
			return cp
		}
	}
	return nil
}

// Messages payload to process

type AddNewProductPayload struct {
	Product Product
}

type AddQuantityToProductPayload struct {
	Product Product
}

type RemoveQuantityToProductPayload struct {
	Product Product
}

// Process processes incoming messages and updates the state accordingly.
func (state *ProductsState) Process(inbox <-chan actor.Message) {
	for {
		msg := <-inbox
		switch paylaod := msg.Body.(type) {

		case AddNewProductPayload:
			slog.Info("AddProductMsg")
			p := paylaod.Product
			state.products = append(state.products, p)

		case AddQuantityToProductPayload:
			slog.Info("AddQuantityProductMsg")
			p := paylaod.Product
			cp := state.getProduct(p.Code)
			if cp != nil {
				cp.Quantity += p.Quantity
				slog.Info("AddQuantityProductMsg update with success", slog.Int("qty", p.Quantity))
			} else {
				slog.Info("AddQuantityProductMsg update fail")
			}

		case RemoveQuantityToProductPayload:
			slog.Info("RemoveQuantityProductMsg")
			p := paylaod.Product
			cp := state.getProduct(p.Code)
			if cp != nil {
				cp.Quantity -= p.Quantity
				slog.Info("RemoveQuantityProductMsg update with success", slog.Int("qty", p.Quantity))
			} else {
				slog.Info("RemoveQuantityProductMsg update fail")
			}

		default:
			slog.Warn("this msg is unknown", slog.String("msg", msg.String()))
		}

		slog.Info("---------------------------------------------------")
		slog.Info("-- Quantity of items in first product", slog.Int("num", state.products[0].Quantity))
		slog.Info("---------------------------------------------------")
	}
}

// Shutdown cleans up the state when the actor is shutting down.
func (state *ProductsState) Shutdown() {
	state.products = make([]Product, 0)
	slog.Info("all product cleaned")
}

func main() {
	slog.Info("---- start of basic example -----")
	productActor, err := actor.NewActor(
		actor.NewAddress("local", "product"),
		NewProductState(),
	)
	if err != nil {
		os.Exit(1)
	}

	err = actor.RegisterActor(&productActor)
	if err != nil {
		os.Exit(1)
	}

	msg := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		AddNewProductPayload{Product{Code: "ABC", Quantity: 5}},
		nil,
	)
	msg2 := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		AddQuantityToProductPayload{Product{Code: "ABC", Quantity: 10}},
		nil,
	)
	msg3 := actor.NewMessage(
		actor.NewAddress("local", "product"),
		actor.NewAddress("local", "product"),
		RemoveQuantityToProductPayload{Product{Code: "ABC", Quantity: 2}},
		nil,
	)

	actor.SendMessage(msg)
	actor.SendMessage(msg2)
	actor.SendMessage(msg3)

	<-time.After(1 * time.Second)
	slog.Info("---- end of basic example -------")
}
