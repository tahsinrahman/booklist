FROM busybox:glibc
EXPOSE 12345
ADD server /bin/booklist
CMD ["booklist"]
