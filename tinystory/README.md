# tinystory

This is a very simple implementation of an interactive story teller. I
currently have two people in mind that will only be able to have minor
control over things with their hands -- so basically probably will be
able to use a phone in one hand, and limited as well (probably only
thumb control).

Hence, to provide maybe something somewhat fun, I thought of coming up
with a very cheap prototype where interactive stories could be put
together very quickly, and inserted into a service (like a web
server), where said phones can point to and run the interactive
stories.

The stories should be of this form:

```
you see a platypus. you choose to:
- greet it
- ignore it
- scream at it
```

So the interface should be very simple (buttons, that lead you to
different parts of the story).

# build / run

Run `make` in the root directory. Then you can run the binary in this
directory. You just need to pass the proper flags to the binary to
point to a story repository, and an assets directory (with the html
templates).

# todo

Since we're dealing with graphs, might be nice to have some checks to
make sure that certain nodes are not orphaned.
