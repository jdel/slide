FROM scratch
ENV SLIDE_SERVER_WS_PORT=4333
COPY slide /slide
ENTRYPOINT ["/slide"]
CMD ["server"]