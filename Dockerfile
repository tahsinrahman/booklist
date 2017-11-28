FROM ubuntu
ADD server server
EXPOSE 12345
CMD ["./server"]
