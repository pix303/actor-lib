# Cinecity

A simple and lightweight implementation of the Actor pattern in Go, designed for better understand the pattern, tested to verify its effectiveness in a personal project [https://github.com/pix303/localemgmt-go](localemgmt-go).
For examples of actor model implementation in Go you should see [https://github.com/vladopajic/go-actor](go-actor) or [https://github.com/anthdm/hollywood](hollywood)

## Overview

The Actor model promises to simplify and streamline the build process. Interactions between application components, providing a model that reflects reality. This is achieved through the exchange of messages within actors. Each actor is responsible for updating its own state.

An Actor can:
- Send messages to other actors
- Change their own state in response to messages
- Subscribe for future message of an actor
- Create new actors

An Actor is composed by:
- an address defined by an area and an id
