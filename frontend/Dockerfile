FROM mintlang/mint:latest

WORKDIR /work

# copy source
COPY . .

EXPOSE 3000

ENTRYPOINT ["/bin/mint"]
CMD ["start", "-h", "0.0.0.0"]
