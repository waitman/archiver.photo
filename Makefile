GO ?=		go
GO_FLAGS ?=	-ldflags "-s -w"
PREFIX ?= 	/usr
CGI-BIN ?=	$(PREFIX)/lib/cgi-bin
ETC ?=		/etc
INSTALL ?=	install

all:
	$(GO) build $(GO_FLAGS) -o build/poster.cgi src/poster.go src/arccommon.go
	$(GO) build $(GO_FLAGS) -o build/arcauth.cgi src/arcauth.go src/arccommon.go
	$(GO) build $(GO_FLAGS) -o build/arcuser src/arcuser.go src/arccommon.go
	$(GO) build $(GO_FLAGS) -o build/arcserial src/arcserial.go
	$(GO) build $(GO_FLAGS) -o build/arcdest src/arcdest.go src/arccommon.go
	$(GO) build $(GO_FLAGS) -o build/arcdestlist src/arcdestlist.go src/arccommon.go
	$(GO) build $(GO_FLAGS) -o build/arctoken src/arctoken.go src/arccommon.go

install:
	$(INSTALL) -s -m 0755 -o root -g root build/poster.cgi $(CGI-BIN)
	$(INSTALL) -s -m 0755 -o root -g root build/arcauth.cgi $(CGI-BIN)
	$(INSTALL) -s -m 0755 -o root -g root build/arcuser $(PREFIX)/bin
	$(INSTALL) -s -m 0755 -o root -g root build/arcserial $(PREFIX)/bin
	$(INSTALL) -s -m 0755 -o root -g root build/arcdest $(PREFIX)/bin
	$(INSTALL) -s -m 0755 -o root -g root build/arcdestlist $(PREFIX)/bin
	$(INSTALL) -s -m 0755 -o root -g root build/arctoken $(PREFIX)/bin
	$(INSTALL) -b -m 0644 -o root -g root etc/archiver.photo.json $(ETC)

clean:
	rm -f build/poster.cgi
	rm -f build/arcuser
	rm -f build/arcauth.cgi
	rm -f build/arcserial
	rm -f build/arcdest
	rm -f build/arcdestlist
	rm -f build/arctoken

deinstall:
	rm -f $(PREFIX)/bin/arcuser
	rm -f $(CGI-BIN)/poster.cgi
	rm -f $(CGI-BIN)/arcauth.cgi
	rm -f $(PREFIX)/bin/arcserial
	rm -f $(PREFIX)/bin/arcdest
	rm -f $(PREFIX)/bin/arcdestlist
	rm -f $(PREFIX)/bin/arctoken
