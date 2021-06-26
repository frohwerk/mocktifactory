FROM registry.access.redhat.com/ubi8/ubi
RUN mkdir -p /app && mkdir -p /repo
WORKDIR /app
COPY ./mocktifactory ./
RUN chmod g+rwx -R /repo
CMD ["/app/mocktifactory"]
