override CFLAGS += -std=c99 -fPIC -Wall -Wextra -Werror -g

all: libwtf.a libwtf.so

walrus.o: walrus.c walrus.h

libwtf.a: walrus.o
	$(AR) rcs $@ $^

libwtf.so: walrus.o
	$(CC) -shared $^ -o $@ -lhiredis

clean:
	rm -f *.o
	rm -f *.a
	rm -f *.so

install:
	install -D walrus.h /usr/local/include/wtf/walrus.h
	install -D libwtf.a /usr/local/lib
	install -D libwtf.so /usr/local/lib

uninstall:
	rm -rf /usr/local/include/wtf
	rm -f /usr/local/lib/libwtf.a
	rm -f /usr/local/lib/libwtf.so
