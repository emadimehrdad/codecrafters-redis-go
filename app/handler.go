package main

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"ECHO": echo,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func echo(args []Value) Value {
	return Value{typ: "bulk", bulk: args[0].bulk}
}
