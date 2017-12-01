FROM busybox:glibc
EXPOSE 12345
ADD booklist /bin/booklist
CMD ["booklist"]
